# 📦 jsonl

A fast, streaming Go package for reading and writing **JSON Lines**.

It is designed for logs, datasets, ingestion pipelines, and large files where memory use, line-level diagnostics, and JSON decoding choice matter.

---

## 📥 Installation

```bash
go get github.com/codenaline/jsonl
```

---

## 🚀 Quick Start

```go
package main

import (
	"log"
	"strings"

	"github.com/codenaline/jsonl"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	input := strings.NewReader("{\"id\":1,\"name\":\"alice\"}\n{\"id\":2,\"name\":\"bob\"}\n")
	r := jsonl.NewReader[Event](input)

	for r.Next() {
		event, err := r.Value()
		if err != nil {
			log.Printf("bad record at line %d offset %d: %v", r.Line(), r.Offset(), err)
			continue
		}

		log.Printf("event: %+v", event)
	}

	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
```

---

## ✨ Features

- Streaming reader powered by `bufio.ReadSlice`, not `bufio.Scanner`
- Generic `Reader[T]` for typed JSON Lines records
- Copy-free short-line fast path for efficient raw line reads
- Line number and byte-offset tracking for diagnostics
- `DecodeInto(*T)` for caller-owned object reuse on hot paths
- Recoverable per-record decode errors
- Optional `WithMaxLineSize` guard for untrusted input
- Custom JSON unmarshal support through `WithUnmarshal`
- Buffered JSON Lines writer
- `Write` for Go values and `WriteBytes` for pre-encoded JSON
- Custom writer marshal support through `WithMarshal`
- No required third-party dependencies

---

## 🛠️ Reading

Create a reader from any `io.Reader`:

```go
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

	process(event)
}

if err := r.Err(); err != nil {
	log.Fatal(err)
}
```

Each call to `Next()` advances to the next JSON Lines record. Decode errors from `Value()` do not stop iteration, so bad records can be logged while valid records continue processing.

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

`Value()` is the ergonomic API. `DecodeInto(*T)` is the high-throughput API.

---

## ✍️ Writing

Create a buffered writer from any `io.Writer`:

```go
var buf bytes.Buffer
w := jsonl.NewWriter(&buf)

_ = w.Write(Event{ID: 1, Name: "alice"})
_ = w.WriteBytes([]byte(`{"id":2,"name":"bob"}`))
_ = w.Flush()
```

Use `Write` when you have a Go value. Use `WriteBytes` when you already have encoded JSON bytes.

---

## ⚙️ Options

### Reader Options

```go
jsonl.WithReaderBufferSize(128 * 1024)
jsonl.WithMaxLineSize(8 * 1024 * 1024)
jsonl.WithUnmarshal(customUnmarshal)
```

### Writer Options

```go
jsonl.WithWriterBufferSize(128 * 1024)
jsonl.WithMarshal(customMarshal)
```

### API Table

| API | Option | Purpose |
| --- | --- | --- |
| `NewReader[T]` | `WithReaderBufferSize` | Set the internal reader buffer size. |
| `NewReader[T]` | `WithMaxLineSize` | Reject lines larger than the configured limit. |
| `NewReader[T]` | `WithUnmarshal` | Replace `encoding/json.Unmarshal`. |
| `NewWriter` | `WithWriterBufferSize` | Set the internal writer buffer size. |
| `NewWriter` | `WithMarshal` | Replace `encoding/json.Marshal`. |

`WithUnmarshal` accepts functions with the same shape as `encoding/json.Unmarshal`:

```go
func(data []byte, v any) error
```

`WithMarshal` accepts functions with the same shape as `encoding/json.Marshal`:

```go
func(v any) ([]byte, error)
```

---

## 🧯 Error Handling

The reader separates stream errors from record decode errors. Empty lines are preserved as records and will return a decode error when decoded as JSON.

- `Next()` reports stream-level failures through `Err()`.
- `Value()` and `DecodeInto()` report decode failures for the current record.
- Decode failures do not stop the iterator.
- `Value()` returns the zero value of `T` on decode failure.
- `DecodeInto(*T)` may leave caller-owned storage partially mutated if unmarshaling fails.

---

## 🔒 Safety

For trusted local files, the default reader is unlimited and avoids imposing arbitrary limits.

For untrusted input, network streams, or user uploads, set `WithMaxLineSize`:

```go
r := jsonl.NewReader[Event](reader, jsonl.WithMaxLineSize(8*1024*1024))
```

This prevents a single oversized line from forcing unbounded memory growth. Oversized accumulated buffers are shrunk back after large records or line-size failures.

---

## 📊 Benchmarks

Current baseline on Linux amd64, Intel i5-1135G7:

```text
BenchmarkReaderRawLines-8      21732 ns/op  1225.09 MB/s   65832 B/op      5 allocs/op
BenchmarkReaderValue-8       673061 ns/op    39.56 MB/s  319785 B/op   6149 allocs/op
BenchmarkReaderDecodeInto-8  665572 ns/op    40.00 MB/s  295233 B/op   5126 allocs/op
BenchmarkWriterWrite-8       169683 ns/op   156.90 MB/s   98431 B/op   1028 allocs/op
BenchmarkWriterWriteBytes-8   14848 ns/op  1793.16 MB/s   65632 B/op      4 allocs/op
```

Run benchmarks locally:

```bash
go test -bench=. -benchmem -run=^$ ./...
```

---

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

Run vet:

```bash
go vet ./...
```

This package includes tests for reader iteration, typed decoding, `DecodeInto`, decode error recovery, maximum line size enforcement, writer buffering, writer options, custom marshal functions, and Go doc examples.

---

## 📌 Status

`jsonl` is ready for the first `v0.1.0` release tag. Reader and writer APIs are implemented, documented, benchmarked, and covered by release-hardening tests.

---

## 🤝 Contributing

Pull requests are welcome! Please see [CONTRIBUTING](CONTRIBUTING.md) for details.

## 👨🏻‍💻 Credits

- [Mahdi Rezaei](https://github.com/mahdirezaei-dev)

## 📄 License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
