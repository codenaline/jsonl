package jsonl

import "fmt"

// DecodeError describes a failure to decode the current JSON Lines record.
type DecodeError struct {
	Line   int64
	Offset int64
	Err    error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("jsonl: decode error at line %d offset %d: %v", e.Line, e.Offset, e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}
