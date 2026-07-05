package jsonl

import (
	"bytes"
	"testing"
)

type writerRecord struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
