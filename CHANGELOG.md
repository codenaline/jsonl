# Changelog

All notable changes to this project will be documented in this file.

## v0.1.0 - 2026-07-10

### Added

- Streaming generic JSON Lines reader with line and byte-offset tracking.
- `DecodeInto` for caller-owned value reuse.
- Configurable reader buffer size, maximum line size, and unmarshal function.
- Buffered JSON Lines writer with value and pre-encoded byte APIs.
- Configurable writer buffer size and marshal function.
- Reader and writer benchmarks.
- Go doc examples for reader, `DecodeInto`, and writer usage.
- Release-hardening tests for EOF, CRLF, empty lines, no-op options, and writer error propagation.

### Notes

- This is the first planned public release.
- Public option names were reviewed before tagging `v0.1.0`.
