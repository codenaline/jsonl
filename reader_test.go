package jsonl

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestReaderIteratesRawLines(t *testing.T) {
	r := NewReader[struct{}](strings.NewReader("{\"id\":1}\n{\"id\":2}\n"))

	var lines []string
	var lineNums []int64
	for r.Next() {
		lines = append(lines, string(r.Bytes()))
		lineNums = append(lineNums, r.Line())
	}

	if err := r.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}

	wantLines := []string{`{"id":1}`, `{"id":2}`}
	if strings.Join(lines, "|") != strings.Join(wantLines, "|") {
		t.Fatalf("lines = %#v, want %#v", lines, wantLines)
	}

	wantLineNums := []int64{1, 2}
	if !slices.Equal(lineNums, wantLineNums) {
		t.Fatalf("line numbers = %#v, want %#v", lineNums, wantLineNums)
	}
}

func TestReaderSupportsCustomBufferSize(t *testing.T) {
	input := strings.Repeat("x", 32) + "\n"
	r := NewReader[struct{}](strings.NewReader(input), WithBufferSize(8))

	if !r.Next() {
		t.Fatalf("Next() = false, want true; err = %v", r.Err())
	}

	if got := string(r.Bytes()); got != strings.Repeat("x", 32) {
		t.Fatalf("Bytes() length = %d, want 32", len(got))
	}

	if r.Next() {
		t.Fatal("Next() = true, want false")
	}
	if err := r.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}

func TestReaderTracksLineStartOffset(t *testing.T) {
	r := NewReader[struct{}](strings.NewReader("aa\nbbb\ncccc\n"), WithBufferSize(2))

	var offsets []int64
	for r.Next() {
		offsets = append(offsets, r.Offset())
	}

	if err := r.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}

	want := []int64{0, 3, 7}
	if !slices.Equal(offsets, want) {
		t.Fatalf("offsets = %#v, want %#v", offsets, want)
	}
}

type readerRecord struct {
	ID int `json:"id"`
}

func TestReaderDecodesTypedValues(t *testing.T) {
	r := NewReader[readerRecord](strings.NewReader("{\"id\":1}\n"))

	if !r.Next() {
		t.Fatalf("Next() = false, want true; err = %v", r.Err())
	}

	got, err := r.Value()
	if err != nil {
		t.Fatalf("Value() error = %v, want nil", err)
	}
	if got.ID != 1 {
		t.Fatalf("Value().ID = %d, want 1", got.ID)
	}
}

func TestReaderUsesCustomDecoder(t *testing.T) {
	r := NewReader[readerRecord](
		strings.NewReader("ignored\n"),
		WithDecoder(func(_ []byte, v any) error {
			v.(*readerRecord).ID = 42
			return nil
		}),
	)

	if !r.Next() {
		t.Fatalf("Next() = false, want true; err = %v", r.Err())
	}

	got, err := r.Value()
	if err != nil {
		t.Fatalf("Value() error = %v, want nil", err)
	}
	if got.ID != 42 {
		t.Fatalf("Value().ID = %d, want 42", got.ID)
	}
}

func TestReaderDecodeErrorIncludesLineAndOffset(t *testing.T) {
	r := NewReader[readerRecord](strings.NewReader("{\"id\":1}\n{bad json}\n"))

	if !r.Next() {
		t.Fatalf("first Next() = false, want true; err = %v", r.Err())
	}
	if _, err := r.Value(); err != nil {
		t.Fatalf("first Value() error = %v, want nil", err)
	}

	if !r.Next() {
		t.Fatalf("second Next() = false, want true; err = %v", r.Err())
	}
	_, err := r.Value()
	if err == nil {
		t.Fatal("second Value() error = nil, want error")
	}

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("Value() error type = %T, want *DecodeError", err)
	}
	if decodeErr.Line != 2 {
		t.Fatalf("DecodeError.Line = %d, want 2", decodeErr.Line)
	}
	if decodeErr.Offset != int64(len("{\"id\":1}\n")) {
		t.Fatalf("DecodeError.Offset = %d, want %d", decodeErr.Offset, len("{\"id\":1}\n"))
	}
	if decodeErr.Unwrap() == nil {
		t.Fatal("DecodeError.Unwrap() = nil, want wrapped decoder error")
	}
}

func TestReaderUsesReadSliceFastPathForShortLines(t *testing.T) {
	r := NewReader[struct{}](strings.NewReader("short\n"), WithBufferSize(64))

	if !r.Next() {
		t.Fatalf("Next() = false, want true; err = %v", r.Err())
	}
	if string(r.Bytes()) != "short" {
		t.Fatalf("Bytes() = %q, want short", r.Bytes())
	}
	if cap(r.buf) != 0 {
		t.Fatalf("internal buffer capacity = %d, want 0 for short-line fast path", cap(r.buf))
	}
}

func TestReaderShrinksOversizedAccumulationBuffer(t *testing.T) {
	r := NewReader[struct{}](
		strings.NewReader(strings.Repeat("x", 128)+"\nsmall\n"),
		WithBufferSize(16),
	)

	if !r.Next() {
		t.Fatalf("first Next() = false, want true; err = %v", r.Err())
	}
	if got := len(r.Bytes()); got != 128 {
		t.Fatalf("first line length = %d, want 128", got)
	}

	if !r.Next() {
		t.Fatalf("second Next() = false, want true; err = %v", r.Err())
	}
	if string(r.Bytes()) != "small" {
		t.Fatalf("second line = %q, want small", r.Bytes())
	}
	if got, wantMax := cap(r.buf), 16; got > wantMax {
		t.Fatalf("internal buffer capacity = %d, want <= %d after shrink", got, wantMax)
	}
}

func TestReaderRejectsShortLineOverMaxSize(t *testing.T) {
	r := NewReader[struct{}](
		strings.NewReader("too-long\n"),
		WithMaxLineSize(3),
	)

	if r.Next() {
		t.Fatal("Next() = true, want false")
	}

	var lineErr *LineTooLongError
	if !errors.As(r.Err(), &lineErr) {
		t.Fatalf("Err() = %T, want *LineTooLongError", r.Err())
	}
	if lineErr.Line != 1 {
		t.Fatalf("LineTooLongError.Line = %d, want 1", lineErr.Line)
	}
	if lineErr.Offset != 0 {
		t.Fatalf("LineTooLongError.Offset = %d, want 0", lineErr.Offset)
	}
	if lineErr.Size != len("too-long\n") {
		t.Fatalf("LineTooLongError.Size = %d, want %d", lineErr.Size, len("too-long\n"))
	}
	if lineErr.Max != 3 {
		t.Fatalf("LineTooLongError.Max = %d, want 3", lineErr.Max)
	}
}

func TestReaderRejectsAccumulatedLineOverMaxSize(t *testing.T) {
	r := NewReader[struct{}](
		strings.NewReader(strings.Repeat("x", 80)+"\n"),
		WithBufferSize(8),
		WithMaxLineSize(64),
	)

	if r.Next() {
		t.Fatal("Next() = true, want false")
	}

	var lineErr *LineTooLongError
	if !errors.As(r.Err(), &lineErr) {
		t.Fatalf("Err() = %T, want *LineTooLongError", r.Err())
	}
	if lineErr.Line != 1 {
		t.Fatalf("LineTooLongError.Line = %d, want 1", lineErr.Line)
	}
	if lineErr.Offset != 0 {
		t.Fatalf("LineTooLongError.Offset = %d, want 0", lineErr.Offset)
	}
	if lineErr.Size <= lineErr.Max {
		t.Fatalf("LineTooLongError.Size = %d, want > max %d", lineErr.Size, lineErr.Max)
	}
	if got, wantMax := cap(r.buf), 8; got > wantMax {
		t.Fatalf("internal buffer capacity = %d, want <= %d after line-too-long error", got, wantMax)
	}
}
