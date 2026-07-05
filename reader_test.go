package jsonl

import (
	"slices"
	"strings"
	"testing"
)

func TestReaderIteratesRawLines(t *testing.T) {
	r := NewReader(strings.NewReader("{\"id\":1}\n{\"id\":2}\n"))

	var lines []string
	var lineNums []int64
	for r.Next() {
		lines = append(lines, string(r.Bytes()))
		lineNums = append(lineNums, r.Line())
	}

	if err := r.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}

	wantLines := []string{`{"id":1}`, `{"id":2}`}
	if strings.Join(lines, "|") != strings.Join(wantLines, "|") {
		t.Fatalf("lines = %#v, want %#v", lines, wantLines)
	}

	wantLineNums := []int64{1, 2}
	if !slices.Equal(lineNums, wantLineNums) {
		t.Fatalf("line numbers = %#v, want %#v", lineNums, wantLineNums)
	}
}

func TestReaderSupportsCustomBufferSize(t *testing.T) {
	input := strings.Repeat("x", 32) + "\n"
	r := NewReader(strings.NewReader(input), WithBufferSize(8))

	if !r.Next() {
		t.Fatalf("Next() = false, want true; err = %v", r.Err())
	}

	if got := string(r.Bytes()); got != strings.Repeat("x", 32) {
		t.Fatalf("Bytes() length = %d, want 32", len(got))
	}

	if r.Next() {
		t.Fatal("Next() = true, want false")
	}
	if err := r.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}
