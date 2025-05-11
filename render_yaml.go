package gostructx

import (
	"fmt"
	"strings"
)

// 导出成 YAML 树
func ToYAML(node InspectNode) string {
	var builder strings.Builder
	writeYAML(&builder, node, 0, false)
	return builder.String()
}

// 递归写 YAML
func writeYAML(b *strings.Builder, node InspectNode, indent int, isListItem bool) {
	prefix := strings.Repeat("  ", indent)

	if node.Name != "root" {
		if node.Value != nil {
			if isListItem {
				prefix = strings.Repeat("  ", indent) + "- "
			}

			fmt.Fprintf(b, "%s%s: %v\n", prefix, node.Name, node.Value)
		} else {
			fmt.Fprintf(b, "%s%s:\n", prefix, node.Name)
		}
	}

	for _, child := range node.Children {
		isSliceElem := strings.HasPrefix(child.Name, "[") // 如果名字是 [0] 这样，表示是 slice
		writeYAML(b, child, indent+1, isSliceElem)
	}
}
