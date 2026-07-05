package jsonl

import (
	"errors"
	"fmt"
)

// ErrNilDecodeTarget is returned when DecodeInto receives a nil destination.
var ErrNilDecodeTarget = errors.New("jsonl: decode target is nil")

// DecodeError describes a failure to decode the current JSON Lines record.
type DecodeError struct {
	Line   int64
	Offset int64
	Err    error
}

// LineTooLongError describes a line that exceeded the configured maximum size.
type LineTooLongError struct {
	Line   int64
	Offset int64
	Size   int
	Max    int
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("jsonl: decode error at line %d offset %d: %v", e.Line, e.Offset, e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

func (e *LineTooLongError) Error() string {
	return fmt.Sprintf("jsonl: line %d too long at offset %d: %d > %d bytes", e.Line, e.Offset, e.Size, e.Max)
}
