package gostructx

import (
	"fmt"
	"reflect"
	"sort"
)

type InspectNode struct {
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Value    any           `json:"value,omitempty"`
	Tag      string        `json:"tag,omitempty"`
	Children []InspectNode `json:"children,omitempty"`
}

type options struct {
	MaxDepth     int
	ShowTag      bool
	SkipTag      string
	SkipEmpty    bool
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

// WithSkipEmpty 设置是否跳过空字段
func WithSkipEmpty(skip bool) Option {
	return func(o *options) {
		o.SkipEmpty = skip
	}
}

func Inspect(val any, opt ...Option) (InspectNode, error) {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return InspectNode{}, fmt.Errorf("invalid input value")
	}

	conf := defaultOptions()
	for _, o := range opt {
		o(conf)
	}

	c := &context{
		conf:    conf,
		Visited: make(map[uintptr]bool),
	}

	rootNode := InspectNode{
		Name:     "root",
		Type:     v.Type().String(),
		Children: c.inspectValue(v, 0),
	}

	return rootNode, nil
}

type context struct {
	conf       *options
	Visited    map[uintptr]bool
	CurrentTag string
}

func (c *context) inspectValue(val reflect.Value, depth int) []InspectNode {
	var nodes []InspectNode

	if depth > c.conf.MaxDepth {
		return nodes
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			if c.conf.SkipEmpty {
				return []InspectNode{}
			}

			return []InspectNode{{
				Name:  "",
				Type:  val.Type().String(),
				Value: nil,
			}}
		}

		ptr := val.Pointer()
		if c.Visited[ptr] {
			return []InspectNode{{
				Name:  "",
				Type:  val.Type().String(),
				Value: "<circular reference>",
			}}
		}

		c.Visited[ptr] = true
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := range val.NumField() {
			field := t.Field(i)

			if c.conf.FilterPrefix != "" && !startsWith(field.Name, c.conf.FilterPrefix) {
				continue
			}

			if c.conf.SkipTag == field.Tag.Get("json") {
				continue
			}

			fieldValue := val.Field(i)
			childNode := InspectNode{
				Name: field.Name,
				Type: fieldValue.Type().String(),
			}

			if c.conf.ShowTag {
				childNode.Tag = field.Tag.Get("json")
			}

			children := c.inspectValue(fieldValue, depth+1)
			if isSimpleValue(fieldValue) {
				childNode.Value = fieldValue.Interface()
			} else if len(children) > 0 {
				childNode.Children = children
			}

			if c.conf.SkipEmpty &&
				(childNode.Value == nil || childNode.Value == "" || childNode.Value == 0) &&
				len(childNode.Children) == 0 {
				continue
			}

			nodes = append(nodes, childNode)
		}
	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			itemVal := val.Index(i)
			childNode := InspectNode{
				Name: fmt.Sprintf("[%d]", i),
				Type: itemVal.Type().String(),
			}

			children := c.inspectValue(itemVal, depth+1)
			if isSimpleValue(itemVal) {
				childNode.Value = itemVal.Interface()
			} else if len(children) > 0 {
				childNode.Children = children
			}

			nodes = append(nodes, childNode)
		}
	case reflect.Map:
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
		})

		for _, key := range keys {
			itemVal := val.MapIndex(key)
			childNode := InspectNode{
				Name: fmt.Sprintf("[%v]", key.Interface()),
				Type: itemVal.Type().String(),
			}

			children := c.inspectValue(itemVal, depth+1)
			if isSimpleValue(itemVal) {
				childNode.Value = itemVal.Interface()
			} else if len(children) > 0 {
				childNode.Children = children
			}

			nodes = append(nodes, childNode)
		}
	case reflect.Interface:
		if val.IsNil() {
			return []InspectNode{{
				Name:  "",
				Type:  "interface",
				Value: nil,
			}}
		} else {
			nodes = c.inspectValue(val.Elem(), depth+1)
		}
	default:
		nodes = append(nodes, InspectNode{
			Name:  "",
			Type:  val.Type().String(),
			Value: val.Interface(),
		})
	}

	return nodes
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// 判断是否是简单类型
func isSimpleValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
