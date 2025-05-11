# GoStructX

GoStructX 是一个轻量的 Go 语言结构体遍历与输出工具。

## 特性

- 支持树形结构输出
- 支持 JSON / Markdown / YAML 格式输出
- 支持防止循环引用死循环
- 支持最大递归深度控制
- 支持跳过空字段
- 支持 Slice / Map 最大展开数量限制
- 支持字段类型统计（比如 int有多少个、string有多少个）
- 超易用的 API，一行调用

## 安装

```bash
go get github.com/xinwuxx/go-struct-x
```

## 使用示例
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

## 示例输出

```markdown
- username (string): shinwu
- age (int): 25
- ports ([]int)
  - [0] (int): 8080
  - [1] (int): 8081
  - ... (truncated): 1 items truncated
```