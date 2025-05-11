package main

import (
	"encoding/json"
	"fmt"

	gostructx "github.com/xinwuxx/go-struct-x"
)

func main() {
	addr := "Earth"
	nestedName := "NestedStruct"

	type Nested struct {
		ID   int     `json:"id"`
		Name *string `json:"name"`
	}

	type Example struct {
		Name    string            `json:"name"`
		Age     int               `json:"age"`
		Address *string           `json:"address"`
		Tags    []string          `json:"tags"`
		Ports   []int             `json:"ports"`
		Props   map[string]string `json:"props"`
		Nested  Nested            `json:"nested"`
		Secret  string            `json:"-"` // 跳过
		Dynamic any               `json:"dynamic"`
	}

	ex := Example{
		Name:    "test",
		Age:     20,
		Address: &addr,
		Tags:    []string{"golang", "reflection"},
		Ports:   []int{8080, 8081},
		Props:   map[string]string{"env": "dev", "version": "1.0"},
		Nested: Nested{
			ID:   1,
			Name: &nestedName,
		},
		Dynamic: map[string]any{
			"enabled": true,
			"count":   42,
		},
	}

	nodes, err := gostructx.Inspect(ex,
		gostructx.WithMaxDepth(5),
		gostructx.WithShowTag(true),
		gostructx.WithSkipEmpty(true),
		gostructx.WithFilterPrefix("A"),
	)

	if err != nil {
		panic(err)
	}

	// 输出漂亮的 JSON
	jsonData, _ := json.MarshalIndent(nodes, "", "  ")
	fmt.Println(string(jsonData))
}
