# jsonl

`jsonl` is a fast, streaming JSON Lines reader and writer for Go.

It is built for large files, logs, ingestion pipelines, and datasets where memory use, line-level diagnostics, and decoder choice matter.

## ✨ Features

- 🚀 Streaming reader with `bufio.ReadSlice`, not `bufio.Scanner`
- ✍️ Buffered JSON Lines writer
- 🧠 Copy-free short-line fast path
- 📏 Optional `WithMaxLineSize` guard for untrusted input
- 📍 Line and byte-offset tracking
- 🧬 Generic `Reader[T]`
- 🔁 `DecodeInto(*T)` for caller-owned object reuse
- 🔌 Custom decoder support through `WithDecoder`
- 🧯 Recoverable per-record decode errors
- 📦 No required third-party dependencies

## 📦 Install

```bash
go get github.com/codenaline/jsonl
```

## 🚀 Quick Start

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

## ✍️ Writing

```go
var buf bytes.Buffer
w := jsonl.NewWriter(&buf)

_ = w.Write(Event{ID: 1, Name: "alice"})
_ = w.WriteBytes([]byte(`{"id":2,"name":"bob"}`))
_ = w.Flush()
```

Use `Write` for Go values and `WriteBytes` when you already have encoded JSON.

## 🔥 Hot Path API

Use `DecodeInto` when you want to reuse caller-owned storage and avoid creating a fresh result value per record.

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

## 🧯 Error Handling

The reader separates stream errors from record decode errors.

`Next()` reports stream-level failures. If reading fails, or a configured limit such as `WithMaxLineSize` is exceeded, `Next()` returns false and `Err()` returns the terminal error. After that, future calls to `Next()` also return false.

`Value()` and `DecodeInto()` report decode failures for the current record. Decode failures do not stop the iterator; callers can log or count the bad record and continue to the next line.

`Value()` returns the zero value of `T` on decode failure, so partially decoded state is not exposed. `DecodeInto(*T)` may leave caller-owned storage partially mutated if a decoder fails, which is expected for the performance-oriented API.

## 🔒 Safety

For trusted local files, the default reader is unlimited and avoids imposing arbitrary limits.

For untrusted input, network streams, or user uploads, set `WithMaxLineSize`:

```go
r := jsonl.NewReader[Event](reader, jsonl.WithMaxLineSize(8*1024*1024))
```

This prevents a single oversized line from forcing unbounded memory growth. Oversized accumulated buffers are shrunk back after large records or line-size failures.

## 📊 Benchmarks

Current baseline on Linux amd64, Intel i5-1135G7:

```text
BenchmarkReaderRawLines-8      21684 ns/op  1227.84 MB/s   65832 B/op      5 allocs/op
BenchmarkReaderValue-8        689506 ns/op    38.61 MB/s  319785 B/op   6149 allocs/op
BenchmarkReaderDecodeInto-8   682884 ns/op    38.99 MB/s  295232 B/op   5126 allocs/op
```

Run locally:

```bash
go test -bench=BenchmarkReader -benchmem -run=^$ ./...
```

## ⚙️ Options

```go
jsonl.WithBufferSize(128 * 1024)
jsonl.WithMaxLineSize(8 * 1024 * 1024)
jsonl.WithDecoder(customUnmarshal)
```

`WithDecoder` accepts functions with the same shape as `encoding/json.Unmarshal`:

```go
func(data []byte, v any) error
```

## 🚧 Status

Early development. The reader API is hardened first, and the buffered writer foundation is now available.
