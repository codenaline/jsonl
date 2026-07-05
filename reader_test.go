package jsonl

import (
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
