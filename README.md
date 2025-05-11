# GoStructX

[![Go Report Card](https://goreportcard.com/badge/github.com/xinwuxx/go-struct-x)](https://goreportcard.com/report/github.com/xinwuxx/go-struct-x)
[![Go](https://github.com/xinwuxx/go-struct-x/actions/workflows/go.yml/badge.svg)](https://github.com/xinwuxx/go-struct-x/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/xinwuxx/go-struct-x.svg)](https://pkg.go.dev/github.com/xinwuxx/go-struct-x)
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)

GoStructX 是一个轻量的 Go 语言结构体遍历与输出工具。支持树形结构输出、多种格式导出、防止循环引用死循环，可以满足一些调试、日志、配置管理的简单需求。

## ✨ 特性

- 🌳 树形结构输出（清晰展示对象层次）
- 📦 支持多种格式（JSON / Markdown / YAML）
- ♻️ 防止循环引用（自动检测）
- 🎯 最大递归深度控制（防止爆栈）
- 🚀 跳过空字段（更简洁的输出）
- 🎯 Slice/Map 最大展开数量限制（防止爆炸）
- 📊 字段类型统计（比如 int、string、map 数量）
- 🔥 Option 链式配置（灵活使用）

## 📦 安装

```bash
go get github.com/xinwuxx/go-struct-x
```

## 🛠 快速使用示例

```go
import "github.com/xinwuxx/go-struct-x"

report, _ := inspector.Inspect(myStruct,
    inspector.WithMaxDepth(5),
    inspector.WithShowTag(true),
    inspector.WithSkipEmpty(true),
)

fmt.Println(inspector.ToMarkdown(report.Root))
fmt.Println(report.Stats.TotalFields)
fmt.Println(report.Stats.MaxDepth)
```

## 🎨 示例输出（Markdown）

```markdown
- username (string): shinwu
- age (int): 25
- ports ([]int)
  - [0] (int): 8080
  - [1] (int): 8081
  - ... (truncated): 1 items truncated
```

## 🔧 支持的 Option 列表

|配置                               |描述                      |
| --------------------------------- | ------------------------ |
|WithMaxDepth(depth int)            |设置最大递归深度            |
|WithSkipTag(tag string)	          |跳过某些字段（比如 json:"-"）|
|WithFilterPrefix(prefix string)	  |仅包含特定路径前缀           |
|WithShowTag(show bool)             |显示字段 tag               |
|WithSkipEmpty(skip bool)	          |跳过空字段                 |
|WithMaxSliceMapLen(max int)	      |`slice/map` 最多展开元素数  |

## 📋 统计信息 (Stats)

- TotalFields：统计总字段数
- MaxDepth：结构体最大嵌套层级
- CircularRef：循环引用检测次数
- FieldTypeCount：各字段类型数量统计