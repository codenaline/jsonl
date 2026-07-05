package jsonl

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

const defaultBufferSize = 64 * 1024

// Reader streams raw JSON Lines records from an io.Reader.
type Reader[T any] struct {
	r *bufio.Reader

	line []byte
	buf  []byte
	err  error

	lineNum     int64
	offset      int64
	nextOffset  int64
	decoder     decoderFunc
	bufferSize  int
	maxLineSize int
}

// NewReader creates a reader for JSON Lines input.
func NewReader[T any](r io.Reader, opts ...Option) *Reader[T] {
	cfg := readerConfig{
		bufferSize: defaultBufferSize,
		decoder:    json.Unmarshal,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Reader[T]{
		r:           bufio.NewReaderSize(r, cfg.bufferSize),
		decoder:     cfg.decoder,
		bufferSize:  cfg.bufferSize,
		maxLineSize: cfg.maxLineSize,
	}
}

// Next advances the reader to the next line.
func (r *Reader[T]) Next() bool {
	if r.err != nil {
		return false
	}

	start := r.nextOffset
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
	r.offset = start
	r.line = trimLineEnding(line)
	return true
}

func (r *Reader[T]) readLine() ([]byte, error) {
	start := r.nextOffset
	part, err := r.r.ReadSlice('\n')
	r.nextOffset += int64(len(part))
	if !errors.Is(err, bufio.ErrBufferFull) {
		if tooLong := r.lineTooLong(start, len(part)); tooLong != nil {
			return nil, tooLong
		}
		return part, err
	}

	r.buf = append(r.buf[:0], part...)
	for {
		part, err = r.r.ReadSlice('\n')
		r.buf = append(r.buf, part...)
		r.nextOffset += int64(len(part))
		if tooLong := r.lineTooLong(start, len(r.buf)); tooLong != nil {
			r.shrinkBuffer()
			return nil, tooLong
		}
		if err == nil || !errors.Is(err, bufio.ErrBufferFull) {
			line := r.buf
			r.shrinkBuffer()
			return line, err
		}
	}
}

func (r *Reader[T]) shrinkBuffer() {
	if cap(r.buf) > r.bufferSize*4 {
		r.buf = make([]byte, 0, r.bufferSize)
	}
}

func (r *Reader[T]) lineTooLong(start int64, size int) error {
	if r.maxLineSize <= 0 || size <= r.maxLineSize {
		return nil
	}
	return &LineTooLongError{
		Line:   r.lineNum + 1,
		Offset: start,
		Size:   size,
		Max:    r.maxLineSize,
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
// The returned slice is only valid until the next call to Next.
func (r *Reader[T]) Bytes() []byte {
	return r.line
}

// Line returns the current 1-based input line number.
func (r *Reader[T]) Line() int64 {
	return r.lineNum
}

// Offset returns the byte offset where the current line begins.
func (r *Reader[T]) Offset() int64 {
	return r.offset
}

// DecodeInto decodes the current JSON Lines record into dst.
func (r *Reader[T]) DecodeInto(dst *T) error {
	if err := r.decoder(r.line, dst); err != nil {
		return &DecodeError{
			Line:   r.lineNum,
			Offset: r.offset,
			Err:    err,
		}
	}
	return nil
}

// Value decodes the current JSON Lines record into T.
func (r *Reader[T]) Value() (T, error) {
	var v T
	if err := r.DecodeInto(&v); err != nil {
		return v, err
	}
	return v, nil
}

// Err returns the terminal error that stopped iteration.
func (r *Reader[T]) Err() error {
	return r.err
}
