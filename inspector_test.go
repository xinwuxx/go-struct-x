package gostructx

import (
	"testing"
)

type Example struct {
	Name    string   `json:"name"`
	Peers   []string `json:"peers"`
	Ports   []int    `json:"ports"`
	Secret  string   `json:"-"`
	Dynamic any
}

type Host struct {
	IP     string
	Number int
}

func TestInspect(t *testing.T) {
	example := Example{
		Name:   "Test Struct",
		Peers:  []string{"node1", "node2"},
		Ports:  []int{6379, 6380},
		Secret: "hidden",
		Dynamic: Host{
			IP:     "0.0.0.0",
			Number: 41,
		},
	}

	nodes, err := Inspect(example)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if len(nodes) == 0 {
		t.Fatalf("Expected nodes, but got empty list")
	}

	var nameFound, peers0Found, peers1Found, ports0Found, ports1Found, ext1Found bool
	for _, node := range nodes {
		if node.Path == "Name" && node.Value == "Test Struct" {
			nameFound = true
		}
		if node.Path == "Peers[0]" && node.Value == "node1" {
			peers0Found = true
		}
		if node.Path == "Peers[1]" && node.Value == "node2" {
			peers1Found = true
		}
		if node.Path == "Ports[0]" && node.Value == 6379 {
			ports0Found = true
		}
		if node.Path == "Ports[1]" && node.Value == 6380 {
			ports1Found = true
		}
		if node.Path == "Dynamic.IP" && node.Value == "0.0.0.0" {
			ext1Found = true
		}
	}

	if !nameFound {
		t.Errorf("Expected Name field, but it was not found")
	}
	if !peers0Found {
		t.Errorf("Expected Peers[0] field, but it was not found")
	}
	if !peers1Found {
		t.Errorf("Expected Peers[1] field, but it was not found")
	}
	if !ports0Found {
		t.Errorf("Expected Ports[0] field, but it was not found")
	}
	if !ports1Found {
		t.Errorf("Expected Ports[1] field, but it was not found")
	}
	if !ext1Found {
		t.Errorf("Expected ext1Found field, but it was not found")
	}

	nodes, err = Inspect(example, WithShowTag(true))
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	var tagFound bool
	for _, node := range nodes {
		if node.Tag == "name" {
			tagFound = true
		}
	}

	if !tagFound {
		t.Errorf("Expected to find a tag, but it was not found")
	}

	nodes, err = Inspect(example, WithMaxDepth(1))
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
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
