package jsonl

import (
	"bufio"
	"encoding/json"
	"io"
)

// Writer writes JSON Lines records to an io.Writer.
type Writer struct {
	w       *bufio.Writer
	marshal marshalFunc
}

// NewWriter creates a buffered JSON Lines writer.
func NewWriter(w io.Writer, opts ...WriterOption) *Writer {
	cfg := writerConfig{
		bufferSize: defaultBufferSize,
		marshal:    json.Marshal,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Writer{
		w:       bufio.NewWriterSize(w, cfg.bufferSize),
		marshal: cfg.marshal,
	}
}

// Write marshals v as JSON and writes it as one JSON Lines record.
func (w *Writer) Write(v any) error {
	b, err := w.marshal(v)
	if err != nil {
		return err
	}
	return w.WriteBytes(b)
}

// WriteBytes writes pre-encoded JSON bytes as one JSON Lines record.
func (w *Writer) WriteBytes(b []byte) error {
	if _, err := w.w.Write(b); err != nil {
		return err
	}
	return w.w.WriteByte('\n')
}

// Flush writes buffered data to the underlying writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}
