package jsonl

import (
	"bufio"
	"errors"
	"io"
)

const defaultBufferSize = 64 * 1024

// Reader streams raw JSON Lines records from an io.Reader.
type Reader struct {
	r *bufio.Reader

	line []byte
	buf  []byte
	err  error

	lineNum int64
}

// NewReader creates a reader for JSON Lines input.
func NewReader(r io.Reader, opts ...Option) *Reader {
	cfg := readerConfig{bufferSize: defaultBufferSize}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Reader{
		r: bufio.NewReaderSize(r, cfg.bufferSize),
	}
}

// Next advances the reader to the next line.
func (r *Reader) Next() bool {
	line, err := r.readLine()
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) == 0 {
			r.err = nil
			return false
		}
		if !errors.Is(err, io.EOF) {
			r.err = err
			return false
		}
	}

	r.lineNum++
	r.line = trimLineEnding(line)
	return true
}

func (r *Reader) readLine() ([]byte, error) {
	r.buf = r.buf[:0]

	for {
		part, err := r.r.ReadSlice('\n')
		r.buf = append(r.buf, part...)
		if err == nil || !errors.Is(err, bufio.ErrBufferFull) {
			return r.buf, err
		}
	}
}

func trimLineEnding(line []byte) []byte {
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line
}

// Bytes returns the raw bytes for the current JSON Lines record.
func (r *Reader) Bytes() []byte {
	return r.line
}

// Line returns the current 1-based input line number.
func (r *Reader) Line() int64 {
	return r.lineNum
}

// Err returns the terminal error that stopped iteration.
func (r *Reader) Err() error {
	return r.err
}
