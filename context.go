package gostructx

import (
	"fmt"
	"reflect"
	"sort"
)

type context struct {
	opts    *options
	visited map[uintptr]bool
	stats   Stats
}

func (c *context) inspectValue(val reflect.Value, depth int) []InspectNode {
	if depth > c.stats.MaxDepth {
		c.stats.MaxDepth = depth
	}
	if depth > c.opts.MaxDepth {
		return nil
	}

	switch val.Kind() {
	case reflect.Pointer:
		return c.handlePointer(val, depth)
	case reflect.Struct:
		return c.handleStruct(val, depth)
	case reflect.Slice, reflect.Array:
		return c.handleSliceArray(val, depth)
	case reflect.Map:
		return c.handleMap(val, depth)
	case reflect.Interface:
		return c.handleInterface(val, depth)
	default:
		return []InspectNode{c.createSimpleNode("", val)}
	}
}

func (c *context) handlePointer(val reflect.Value, depth int) []InspectNode {
	if val.IsNil() {
		if c.opts.SkipEmpty {
			return nil
		}
		return []InspectNode{{
			Name:  "",
			Type:  val.Type().String(),
			Value: nil,
		}}
	}

	ptr := val.Pointer()
	if c.visited[ptr] {
		c.stats.CircularRef++
		return []InspectNode{{
			Name:  "",
			Type:  val.Type().String(),
			Value: "<circular reference>",
		}}
	}

	c.visited[ptr] = true
	return c.inspectValue(val.Elem(), depth)
}

func (c *context) handleStruct(val reflect.Value, depth int) []InspectNode {
	t := val.Type()
	var nodes []InspectNode

	for i := range val.NumField() {
		c.stats.TotalFields++
		field := t.Field(i)
		fieldValue := val.Field(i)

		if c.opts.FilterPrefix != "" && !startsWith(field.Name, c.opts.FilterPrefix) {
			continue
		}
		if c.opts.SkipTag == field.Tag.Get("json") {
			continue
		}

		node := InspectNode{
			Name: field.Name,
			Type: fieldValue.Type().String(),
		}
		if c.opts.ShowTag {
			node.Tag = field.Tag.Get("json")
		}

		if isSimpleValue(fieldValue) {
			node.Value = fieldValue.Interface()
		} else {
			children := c.inspectValue(fieldValue, depth+1)
			if len(children) > 0 {
				node.Children = children
			}
		}

		if c.opts.SkipEmpty && isEmptyNode(node) {
			continue
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func (c *context) handleSliceArray(val reflect.Value, depth int) []InspectNode {
	length := val.Len()
	maxLen := length
	if c.opts.MaxSliceMapLen > 0 && length > c.opts.MaxSliceMapLen {
		maxLen = c.opts.MaxSliceMapLen
	}

	var nodes []InspectNode
	for i := range maxLen {
		itemVal := val.Index(i)
		node := InspectNode{
			Name: fmt.Sprintf("[%d]", i),
			Type: itemVal.Type().String(),
		}

		if isSimpleValue(itemVal) {
			node.Value = itemVal.Interface()
		} else {
			children := c.inspectValue(itemVal, depth+1)
			if len(children) > 0 {
				node.Children = children
			}
		}

		nodes = append(nodes, node)
	}

	if maxLen < length {
		nodes = append(nodes, InspectNode{
			Name:  "...",
			Type:  "truncated",
			Value: fmt.Sprintf("%d items truncated", length-maxLen),
		})
	}
	return nodes
}

func (c *context) handleMap(val reflect.Value, depth int) []InspectNode {
	keys := val.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
	})

	length := len(keys)
	maxLen := length
	if c.opts.MaxSliceMapLen > 0 && length > c.opts.MaxSliceMapLen {
		maxLen = c.opts.MaxSliceMapLen
	}

	var nodes []InspectNode
	for _, key := range keys[:maxLen] {
		itemVal := val.MapIndex(key)
		node := InspectNode{
			Name: fmt.Sprintf("[%v]", key.Interface()),
			Type: itemVal.Type().String(),
		}

		if isSimpleValue(itemVal) {
			node.Value = itemVal.Interface()
		} else {
			children := c.inspectValue(itemVal, depth+1)
			if len(children) > 0 {
				node.Children = children
			}
		}

		nodes = append(nodes, node)
	}

	if maxLen < length {
		nodes = append(nodes, InspectNode{
			Name:  "...",
			Type:  "truncated",
			Value: fmt.Sprintf("%d items truncated", length-maxLen),
		})
	}
	return nodes
}

func (c *context) handleInterface(val reflect.Value, depth int) []InspectNode {
	if val.IsNil() {
		return []InspectNode{{
			Name:  "",
			Type:  "interface",
			Value: nil,
		}}
	}
	return c.inspectValue(val.Elem(), depth+1)
}

func (c *context) createSimpleNode(name string, val reflect.Value) InspectNode {
	typeStr := val.Type().String()
	c.stats.FieldTypeCount[typeStr]++
	return InspectNode{
		Name:  name,
		Type:  typeStr,
		Value: val.Interface(),
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

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

func isEmptyNode(n InspectNode) bool {
	return (n.Value == nil || n.Value == "" || n.Value == 0) && len(n.Children) == 0
}
