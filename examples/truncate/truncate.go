package main

import (
	"fmt"

	gostructx "github.com/xinwuxx/go-struct-x"
)

func main() {
	type Example struct {
		Ports []int             `json:"ports"`
		Name  map[string]string `json:"name"`
	}

	ex := Example{
		Ports: []int{8080, 8081, 8082, 8083, 8084, 8085, 8086},
		Name: map[string]string{
			"os":     "mac",
			"cpu":    "intel",
			"gpu":    "amd",
			"memory": "64gb",
		},
	}

	report, _ := gostructx.Inspect(ex,
		gostructx.WithMaxDepth(5),
		gostructx.WithShowTag(true),
		gostructx.WithSkipEmpty(true),
		gostructx.WithMaxSliceMapLen(3), // ✅ 最多展开 3 个元素
	)

	markdown := gostructx.ToMarkdown(report.Root)
	fmt.Println(markdown)
}
