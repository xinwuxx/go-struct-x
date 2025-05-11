package gostructx

import (
	"testing"
)

type Example struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Peers   []string `json:"peers"`
	Ports   []int    `json:"ports"`
	Secret  string   `json:"-"`
	Dynamic any
	Empty   string
}

func TestInspect(t *testing.T) {
	example := Example{
		Name:   "Test Struct",
		Age:    25,
		Peers:  []string{"node1", "node2"},
		Ports:  []int{6379, 6380},
		Secret: "hidden",
		Dynamic: map[string]int{
			"ext1": 1000,
			"ext2": 2,
		},
	}

	// 测试基本功能
	t.Run("basic inspection", func(t *testing.T) {
		node, err := Inspect(example)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		// 检查根节点
		if node.Root.Name != "root" {
			t.Errorf("Expected root node name to be 'root', got %s", node.Root.Name)
		}

		// 检查基本字段
		foundName := false
		foundAge := false
		for _, child := range node.Root.Children {
			if child.Name == "Name" && child.Value == "Test Struct" {
				foundName = true
			}
			if child.Name == "Age" && child.Value == 25 {
				foundAge = true
			}
		}
		if !foundName {
			t.Error("Expected to find Name field with value 'Test Struct'")
		}
		if !foundAge {
			t.Error("Expected to find Age field with value 25")
		}
	})

	// 测试 WithShowTag 选项
	t.Run("with show tag", func(t *testing.T) {
		node, err := Inspect(example, WithShowTag(true))
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		foundTag := false
		for _, child := range node.Root.Children {
			if child.Name == "Name" && child.Tag == "name" {
				foundTag = true
				break
			}
		}
		if !foundTag {
			t.Error("Expected to find tag 'name' for Name field")
		}
	})

	// 测试 WithSkipTag 选项
	t.Run("with skip tag", func(t *testing.T) {
		node, err := Inspect(example, WithSkipTag("-"))
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		for _, child := range node.Root.Children {
			if child.Name == "Secret" {
				t.Error("Expected Secret field to be skipped")
			}
		}
	})

	// 测试 WithSkipEmpty 选项
	t.Run("with skip empty", func(t *testing.T) {
		node, err := Inspect(example, WithSkipEmpty(true))
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		for _, child := range node.Root.Children {
			if child.Name == "Empty" {
				t.Error("Expected Empty field to be skipped")
			}
		}
	})

	// 测试 WithMaxDepth 选项
	t.Run("with max depth", func(t *testing.T) {
		node, err := Inspect(example, WithMaxDepth(1))
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		// 检查是否只遍历到第一层
		for _, child := range node.Root.Children {
			if child.Name == "Dynamic" && len(child.Children) > 0 {
				t.Error("Expected Dynamic field to not have children due to max depth")
			}
		}
	})

	// 测试 WithFilterPrefix 选项
	t.Run("with filter prefix", func(t *testing.T) {
		node, err := Inspect(example, WithFilterPrefix("P"))
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		// 检查是否只包含以 P 开头的字段
		for _, child := range node.Root.Children {
			if child.Name[0] != 'P' {
				t.Errorf("Expected field name to start with 'P', got %s", child.Name)
			}
		}
	})
}

func TestInspectWithNil(t *testing.T) {
	// 测试空指针
	t.Run("nil pointer", func(t *testing.T) {
		var ptr *Example
		node, err := Inspect(ptr)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
		if node.Root.Value != nil {
			t.Error("Expected nil value for nil pointer")
		}
	})

	// 测试空接口
	t.Run("nil interface", func(t *testing.T) {
		var iface any
		node, err := Inspect(iface)
		if err == nil {
			t.Fatalf("Expected %v, but got error", err)
		}
		if node.Root.Value != nil {
			t.Error("Expected nil value for nil interface")
		}
	})
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		s, prefix string
		expected  bool
	}{
		{"Hello", "He", true},
		{"Hello", "Ho", false},
		{"", "He", false},
		{"Hello", "", true},
		{"Hello", "Hello", true},
		{"Hello", "HelloWorld", false},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := startsWith(test.s, test.prefix)
			if result != test.expected {
				t.Errorf("For %q, expected %v but got %v", test.s, test.expected, result)
			}
		})
	}
}
