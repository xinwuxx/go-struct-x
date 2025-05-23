package main

import (
	"fmt"

	gostructx "github.com/xinwuxx/go-struct-x"
)

func main() {
	type Example struct {
		Name  string `json:"name"`
		Ports []int  `json:"ports"`
	}

	ex := Example{
		Name:  "server",
		Ports: []int{8080, 8081},
	}

	tree, _ := gostructx.Inspect(ex,
		gostructx.WithMaxDepth(5),
		gostructx.WithShowTag(true),
		gostructx.WithSkipEmpty(true),
	)

	// 转成 YAML
	yamlOutput := gostructx.ToYAML(tree.Root)
	fmt.Println(yamlOutput)
}
