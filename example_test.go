package jsonl_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/codenaline/jsonl"
)

type exampleEvent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ExampleNewReader() {
	input := strings.NewReader("{\"id\":1,\"name\":\"alice\"}\n{\"id\":2,\"name\":\"bob\"}\n")
	r := jsonl.NewReader[exampleEvent](input)

	for r.Next() {
		event, err := r.Value()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%d %s\n", event.ID, event.Name)
	}
	if err := r.Err(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1 alice
	// 2 bob
}

func ExampleReader_DecodeInto() {
	input := strings.NewReader("{\"id\":1,\"name\":\"alice\"}\n{\"id\":2,\"name\":\"bob\"}\n")
	r := jsonl.NewReader[exampleEvent](input)

	var event exampleEvent
	for r.Next() {
		if err := r.DecodeInto(&event); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%d %s\n", event.ID, event.Name)
	}
	if err := r.Err(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1 alice
	// 2 bob
}

func ExampleNewWriter() {
	var buf bytes.Buffer
	w := jsonl.NewWriter(&buf)

	if err := w.Write(exampleEvent{ID: 1, Name: "alice"}); err != nil {
		fmt.Println(err)
	}
	if err := w.WriteBytes([]byte(`{"id":2,"name":"bob"}`)); err != nil {
		fmt.Println(err)
	}
	if err := w.Flush(); err != nil {
		fmt.Println(err)
	}

	fmt.Print(buf.String())

	// Output:
	// {"id":1,"name":"alice"}
	// {"id":2,"name":"bob"}
}
