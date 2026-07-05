# jsonl

`jsonl` is a Go package for high-performance JSON Lines processing.

The package is being built incrementally with a small, stable public API. The
first production goals are:

- streaming reads without `bufio.Scanner` line-size limits
- buffered writes for JSON Lines output
- precise line and offset errors
- predictable memory use for large files
- no required third-party dependencies

## Status

Early development. Reader and writer APIs will be added in small, tested
commits.

## Error handling

The reader separates stream errors from record decode errors.

`Next()` reports stream-level failures. If reading fails, or a configured limit
such as `WithMaxLineSize` is exceeded, `Next()` returns false and `Err()`
returns the terminal error. After that, future calls to `Next()` also return
false.

`Value()` reports decode failures for the current record. A `Value()` error does
not stop the iterator; callers can log or count the bad record and continue to
the next line.

For hot paths, use `DecodeInto(*T)` to decode into caller-owned storage and
avoid creating a fresh result value for each record.
