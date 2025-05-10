package gostructx

import (
	"fmt"
	"reflect"
	"sort"
)

type InspectNode struct {
	Path  string `json:"path"`
	Type  string `json:"type"`
	Value any    `json:"value"`
	Tag   string `json:"tag"`
}

type options struct {
	MaxDepth     int
	ShowTag      bool
	SkipTag      string
	FilterPrefix string
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		MaxDepth: 5,
		SkipTag:  "-",
	}
}

// WithMaxDepth 设置最大递归深度
func WithMaxDepth(depth int) Option {
	return func(o *options) {
		o.MaxDepth = depth
	}
}

// WithSkipTag 设置跳过的 tag（比如 json:"-")
func WithSkipTag(tag string) Option {
	return func(o *options) {
		o.SkipTag = tag
	}
}

// WithFilterPrefix 设置只包含特定路径前缀的字段
func WithFilterPrefix(prefix string) Option {
	return func(o *options) {
		o.FilterPrefix = prefix
	}
}

// WithShowTag 设置是否显示 tag
func WithShowTag(show bool) Option {
	return func(o *options) {
		o.ShowTag = show
	}
}

func Inspect(val any, opt ...Option) ([]InspectNode, error) {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return nil, fmt.Errorf("invalid input value")
	}

	conf := defaultOptions()
	for _, o := range opt {
		o(conf)
	}

	c := &context{
		conf:    conf,
		Visited: make(map[uintptr]bool),
	}

	nodes := c.inspectStruct(v, "", 0)
	return nodes, nil
}

type context struct {
	conf       *options
	Visited    map[uintptr]bool
	CurrentTag string
}

func (c *context) inspectStruct(val reflect.Value, path string, depth int) []InspectNode {
	var nodes []InspectNode

	if depth > c.conf.MaxDepth {
		return nodes
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			nodes = append(nodes, InspectNode{
				Path:  path,
				Type:  val.Type().String(),
				Value: nil,
			})

			return nodes
		}

		ptr := val.Pointer()
		if c.Visited[ptr] {
			node := InspectNode{
				Path:  path,
				Type:  val.Type().String(),
				Value: "<circular reference>",
			}

			if c.conf.ShowTag {
				node.Tag = c.CurrentTag
			}

			nodes = append(nodes, node)

			return nodes
		}

		c.Visited[ptr] = true
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := range val.NumField() {
			field := val.Type().Field(i)

			if c.conf.FilterPrefix != "" && startsWith(field.Name, c.conf.FilterPrefix) {
				continue
			}

			newPath := field.Name
			if path != "" {
				newPath = path + "." + newPath
			}

			if c.conf.ShowTag {
				c.CurrentTag = field.Tag.Get("json")
			}

			fieldValue := val.Field(i)
			newNodes := c.inspectStruct(fieldValue, newPath, depth+1)
			nodes = append(nodes, newNodes...)
		}
	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			newNodes := c.inspectStruct(val.Index(i), newPath, depth+1)
			nodes = append(nodes, newNodes...)
		}
	case reflect.Map:
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
		})

		for _, key := range keys {
			newPath := fmt.Sprintf("%s[%v]", path, key)
			newNodes := c.inspectStruct(val.MapIndex(key), newPath, depth+1)
			nodes = append(nodes, newNodes...)
		}
	case reflect.Interface:
	default:
		node := InspectNode{
			Path:  path,
			Type:  val.Type().String(),
			Value: val.Interface(),
		}

		if c.conf.ShowTag {
			node.Tag = c.CurrentTag
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
