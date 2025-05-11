package main

import (
	"encoding/json"
	"fmt"

	gostructx "github.com/xinwuxx/go-struct-x"
)

func main() {
	type Example struct {
		Name  string `json:"username"`
		Ports []int  `json:"ports"`
	}

	ex := Example{
		Name:  "server",
		Ports: []int{8080, 8081, 8082},
	}

	report, _ := gostructx.Inspect(ex,
		gostructx.WithMaxDepth(5),
		gostructx.WithShowTag(true),
		gostructx.WithSkipEmpty(true),
	)

	// 打印统计信息
	fmt.Println("总字段数：", report.Stats.TotalFields)
	fmt.Println("最大嵌套层级：", report.Stats.MaxDepth)
	fmt.Println("循环引用数量：", report.Stats.CircularRef)

	// 打印树
	jsonData, _ := json.MarshalIndent(report.Root, "", "  ")
	fmt.Println(string(jsonData))
}
