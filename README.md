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
