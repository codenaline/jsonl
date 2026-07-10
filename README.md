# 📦 jsonl

A professional Go package for reading and writing **JSON Lines** streams.

It is built for large files, logs, ingestion pipelines, and datasets where memory use, line-level diagnostics, and decoder choice matter.

---

## 🚀 Features

- Streaming reader powered by `bufio.ReadSlice`, not `bufio.Scanner`
- Generic `Reader[T]` for typed JSON Lines records
- Copy-free short-line fast path for efficient raw line reads
- Line number and byte-offset tracking for diagnostics
- `DecodeInto(*T)` for caller-owned object reuse on hot paths
- Recoverable per-record decode errors
- Optional `WithMaxLineSize` guard for untrusted input
- Custom decoder support through `WithDecoder`
- Buffered JSON Lines writer
- `Write` for Go values and `WriteBytes` for pre-encoded JSON
- Custom writer marshal support through `WithMarshal`
- No required third-party dependencies

---

## 📥 Installation

```bash
go get github.com/codenaline/jsonl
```

---

## 🛠️ Reading

Create a reader for any `io.Reader`:

```go
package main

import (
	"log"
	"os"

	"github.com/codenaline/jsonl"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	f, err := os.Open("events.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := jsonl.NewReader[Event](
		f,
		jsonl.WithMaxLineSize(8*1024*1024),
	)

	for r.Next() {
		event, err := r.Value()
		if err != nil {
			log.Printf("bad record at line %d offset %d: %v", r.Line(), r.Offset(), err)
			continue
		}

		_ = event
	}

	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
```

Each call to `Next()` advances to the next JSON Lines record.

Decode errors from `Value()` do not stop iteration, so bad records can be logged while valid records continue processing.

---

## 🔥 Hot Path Decoding

Use `DecodeInto` when you want to reuse caller-owned storage and avoid creating a fresh result value per record:

```go
var event Event

for r.Next() {
	if err := r.DecodeInto(&event); err != nil {
		continue
	}

	process(event)
}
```

`Value()` is the ergonomic API.

`DecodeInto(*T)` is the high-throughput API.

---

## ✍️ Writing

Create a buffered writer for any `io.Writer`:

```go
var buf bytes.Buffer
w := jsonl.NewWriter(&buf)

_ = w.Write(Event{ID: 1, Name: "alice"})
_ = w.WriteBytes([]byte(`{"id":2,"name":"bob"}`))
_ = w.Flush()
```

Use `Write` when you have a Go value.

Use `WriteBytes` when you already have encoded JSON bytes.

---

## ⚙️ Options

Reader options:

```go
jsonl.WithBufferSize(128 * 1024)
jsonl.WithMaxLineSize(8 * 1024 * 1024)
jsonl.WithDecoder(customUnmarshal)
```

Writer options:

```go
jsonl.WithWriterBufferSize(128 * 1024)
jsonl.WithMarshal(customMarshal)
```

### API Table

| API | Option | Purpose |
| --- | --- | --- |
| `NewReader[T]` | `WithBufferSize` | Set the internal reader buffer size. |
| `NewReader[T]` | `WithMaxLineSize` | Reject lines larger than the configured limit. |
| `NewReader[T]` | `WithDecoder` | Replace `encoding/json.Unmarshal`. |
| `NewWriter` | `WithWriterBufferSize` | Set the internal writer buffer size. |
| `NewWriter` | `WithMarshal` | Replace `encoding/json.Marshal`. |

`WithDecoder` accepts functions with the same shape as `encoding/json.Unmarshal`:

```go
func(data []byte, v any) error
```

`WithMarshal` accepts functions with the same shape as `encoding/json.Marshal`:

```go
func(v any) ([]byte, error)
```

---

## 🔒 Safety

For trusted local files, the default reader is unlimited and avoids imposing arbitrary limits.

For untrusted input, network streams, or user uploads, set `WithMaxLineSize`:

```go
r := jsonl.NewReader[Event](reader, jsonl.WithMaxLineSize(8*1024*1024))
```

This prevents a single oversized line from forcing unbounded memory growth.

Oversized accumulated buffers are shrunk back after large records or line-size failures.

---

## 🧯 Error Handling

The reader separates stream errors from record decode errors.

- `Next()` reports stream-level failures through `Err()`.
- `Value()` and `DecodeInto()` report decode failures for the current record.
- Decode failures do not stop the iterator.
- `Value()` returns the zero value of `T` on decode failure.
- `DecodeInto(*T)` may leave caller-owned storage partially mutated if a decoder fails.

---

## 📊 Benchmarks

Current baseline on Linux amd64, Intel i5-1135G7:

```text
BenchmarkReaderRawLines-8      21684 ns/op  1227.84 MB/s   65832 B/op      5 allocs/op
BenchmarkReaderValue-8       656898 ns/op    40.53 MB/s  319785 B/op   6149 allocs/op
BenchmarkReaderDecodeInto-8  644326 ns/op    41.32 MB/s  295233 B/op   5126 allocs/op
BenchmarkWriterWrite-8       159367 ns/op   167.06 MB/s   98431 B/op   1028 allocs/op
BenchmarkWriterWriteBytes-8   14699 ns/op  1811.34 MB/s   65632 B/op      4 allocs/op
```

Run benchmarks locally:

```bash
go test -bench=. -benchmem -run=^$ ./...
```

---

## 🧪 Testing

This package ships with tests for:

- Reader iteration and raw byte access
- Typed decoding with `Value`
- Caller-owned decoding with `DecodeInto`
- Decode error wrapping and recovery
- Maximum line size enforcement
- Writer buffering and flushing
- Writer options and custom marshal functions
- Go doc examples

Run the test suite:

```bash
go test ./...
```

Run vet:

```bash
go vet ./...
```

---

## 📌 Summary

- Reads and writes JSON Lines streams
- Supports typed generic readers
- Tracks line numbers and byte offsets
- Allows recoverable per-record decode errors
- Provides hot-path decoding through `DecodeInto`
- Supports custom decoder and marshal functions
- Includes reader and writer benchmarks
- Ready for `v0.1.0` release review

---

## 🤝 Contributing

Pull requests are welcome! Please see [CONTRIBUTING](CONTRIBUTING.md) for details.

## 👨🏻‍💻 Credits

- [Mahdi Rezaei](https://github.com/mahdirezaei-dev)

## 📄 License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
