package gostructx

import (
	"fmt"
	"strings"
)

// 导出成 Markdown 树
func ToMarkdown(node InspectNode) string {
	var builder strings.Builder
	writeMarkdown(&builder, node, 0)
	return builder.String()
}

// 递归写入 Markdown
func writeMarkdown(b *strings.Builder, node InspectNode, indent int) {
	prefix := strings.Repeat("  ", indent)
	if node.Name != "root" {
		// 打印当前节点
		if node.Value != nil {
			fmt.Fprintf(b, "%s- %s (%s): %v\n", prefix, node.Name, node.Type, node.Value)
		} else {
			fmt.Fprintf(b, "%s- %s (%s)\n", prefix, node.Name, node.Type)
		}
	}

	// 递归子节点
	for _, child := range node.Children {
		writeMarkdown(b, child, indent+1)
	}
}
