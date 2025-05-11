package main

import (
	"fmt"

	gostructx "github.com/xinwuxx/go-struct-x"
)

func main() {
	type Example struct {
		Name  string `json:"username"`
		Age   int    `json:"age"`
		Ports []int  `json:"ports"`
	}

	ex := Example{
		Name:  "server",
		Age:   18,
		Ports: []int{8080, 8081, 8082},
	}

	report, _ := gostructx.Inspect(ex,
		gostructx.WithMaxDepth(5),
		gostructx.WithShowTag(true),
		gostructx.WithSkipEmpty(true),
	)

	fmt.Println("stats:")
	for t, c := range report.Stats.FieldTypeCount {
		fmt.Printf("- %s: %d\n", t, c)
	}
}
