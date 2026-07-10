package jsonl

import (
	"bytes"
	"errors"
	"testing"
)

type writerRecord struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestWriterWritesJSONLines(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.Write(writerRecord{ID: 1, Name: "alice"}); err != nil {
		t.Fatalf("Write() error = %v, want nil", err)
	}
	if err := w.Write(writerRecord{ID: 2, Name: "bob"}); err != nil {
		t.Fatalf("Write() error = %v, want nil", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}

	want := "{\"id\":1,\"name\":\"alice\"}\n{\"id\":2,\"name\":\"bob\"}\n"
	if got := buf.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriterWritesPreEncodedJSONLine(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteBytes([]byte(`{"id":1}`)); err != nil {
		t.Fatalf("WriteBytes() error = %v, want nil", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}

	if got, want := buf.String(), "{\"id\":1}\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriterBuffersUntilFlush(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteBytes([]byte(`{"id":1}`)); err != nil {
		t.Fatalf("WriteBytes() error = %v, want nil", err)
	}
	if got := buf.String(); got != "" {
		t.Fatalf("buffer before Flush() = %q, want empty", got)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}
	if got, want := buf.String(), "{\"id\":1}\n"; got != want {
		t.Fatalf("buffer after Flush() = %q, want %q", got, want)
	}
}

func TestWriterSupportsCustomBufferSize(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, WithWriterBufferSize(1))

	if err := w.WriteBytes([]byte(`{"id":1}`)); err != nil {
		t.Fatalf("WriteBytes() error = %v, want nil", err)
	}

	if got := buf.String(); got == "" {
		t.Fatal("buffer before Flush() is empty, want custom small buffer to write through")
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}
	if got, want := buf.String(), "{\"id\":1}\n"; got != want {
		t.Fatalf("buffer after Flush() = %q, want %q", got, want)
	}
}

func TestWriterIgnoresNonPositiveBufferSize(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, WithWriterBufferSize(0))

	if err := w.WriteBytes([]byte(`{"id":1}`)); err != nil {
		t.Fatalf("WriteBytes() error = %v, want nil", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}
	if got, want := buf.String(), "{\"id\":1}\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriterIgnoresNilMarshal(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, WithMarshal(nil))

	if err := w.Write(writerRecord{ID: 1, Name: "alice"}); err != nil {
		t.Fatalf("Write() error = %v, want nil", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}
	if got, want := buf.String(), "{\"id\":1,\"name\":\"alice\"}\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriterPropagatesWriteError(t *testing.T) {
	wantErr := errors.New("write failed")
	w := NewWriter(failingWriter{err: wantErr}, WithWriterBufferSize(1))

	if err := w.WriteBytes([]byte(`{"id":1}`)); !errors.Is(err, wantErr) {
		t.Fatalf("WriteBytes() error = %v, want %v", err, wantErr)
	}
}

func TestWriterPropagatesFlushError(t *testing.T) {
	wantErr := errors.New("flush failed")
	w := NewWriter(failingWriter{err: wantErr})

	if err := w.WriteBytes([]byte(`{"id":1}`)); err != nil {
		t.Fatalf("WriteBytes() error = %v, want nil before flush", err)
	}
	if err := w.Flush(); !errors.Is(err, wantErr) {
		t.Fatalf("Flush() error = %v, want %v", err, wantErr)
	}
}

func TestWriterUsesCustomMarshal(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, WithMarshal(func(v any) ([]byte, error) {
		record := v.(writerRecord)
		return []byte(`{"custom":"` + record.Name + `"}`), nil
	}))

	if err := w.Write(writerRecord{ID: 1, Name: "alice"}); err != nil {
		t.Fatalf("Write() error = %v, want nil", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}

	if got, want := buf.String(), "{\"custom\":\"alice\"}\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriterReturnsCustomMarshalError(t *testing.T) {
	wantErr := errors.New("marshal failed")
	var buf bytes.Buffer
	w := NewWriter(&buf, WithMarshal(func(any) ([]byte, error) {
		return nil, wantErr
	}))

	if err := w.Write(writerRecord{ID: 1}); !errors.Is(err, wantErr) {
		t.Fatalf("Write() error = %v, want %v", err, wantErr)
	}
	if got := buf.String(); got != "" {
		t.Fatalf("buffer = %q, want empty", got)
	}
}
