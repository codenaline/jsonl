package jsonl

import (
	"bytes"
	"testing"
)

type benchmarkRecord struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var benchmarkJSONL = bytes.Repeat([]byte(`{"id":123,"name":"alice"}`+"\n"), 1024)

func BenchmarkReaderRawLines(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchmarkJSONL)))

	for i := 0; i < b.N; i++ {
		r := NewReader[struct{}](bytes.NewReader(benchmarkJSONL))
		var count int
		for r.Next() {
			count += len(r.Bytes())
		}
		if err := r.Err(); err != nil {
			b.Fatal(err)
		}
		if count == 0 {
			b.Fatal("no records read")
		}
	}
}

func BenchmarkReaderValue(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchmarkJSONL)))

	for i := 0; i < b.N; i++ {
		r := NewReader[benchmarkRecord](bytes.NewReader(benchmarkJSONL))
		var sum int
		for r.Next() {
			v, err := r.Value()
			if err != nil {
				b.Fatal(err)
			}
			sum += v.ID
		}
		if err := r.Err(); err != nil {
			b.Fatal(err)
		}
		if sum == 0 {
			b.Fatal("no records decoded")
		}
	}
}

func BenchmarkReaderDecodeInto(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchmarkJSONL)))

	for i := 0; i < b.N; i++ {
		r := NewReader[benchmarkRecord](bytes.NewReader(benchmarkJSONL))
		var record benchmarkRecord
		var sum int
		for r.Next() {
			if err := r.DecodeInto(&record); err != nil {
				b.Fatal(err)
			}
			sum += record.ID
		}
		if err := r.Err(); err != nil {
			b.Fatal(err)
		}
		if sum == 0 {
			b.Fatal("no records decoded")
		}
	}
}
