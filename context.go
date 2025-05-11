package gostructx

import (
	"fmt"
	"reflect"
	"sort"
)

type context struct {
	opts    *options
	Visited map[uintptr]bool
	stats   Stats
}

func (c *context) inspectValue(val reflect.Value, depth int) []InspectNode {
	var nodes []InspectNode

	if depth > c.stats.MaxDepth {
		c.stats.MaxDepth = depth
	}

	if depth > c.opts.MaxDepth {
		return nodes
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			if c.opts.SkipEmpty {
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
			c.stats.CircularRef++
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
			c.stats.TotalFields++
			field := t.Field(i)

			if c.opts.FilterPrefix != "" && !startsWith(field.Name, c.opts.FilterPrefix) {
				continue
			}

			if c.opts.SkipTag == field.Tag.Get("json") {
				continue
			}

			fieldValue := val.Field(i)
			childNode := InspectNode{
				Name: field.Name,
				Type: fieldValue.Type().String(),
			}

			if c.opts.ShowTag {
				childNode.Tag = field.Tag.Get("json")
			}

			children := c.inspectValue(fieldValue, depth+1)
			if isSimpleValue(fieldValue) {
				childNode.Value = fieldValue.Interface()
			} else if len(children) > 0 {
				childNode.Children = children
			}

			if c.opts.SkipEmpty &&
				(childNode.Value == nil || childNode.Value == "" || childNode.Value == 0) &&
				len(childNode.Children) == 0 {
				continue
			}

			nodes = append(nodes, childNode)
		}
	case reflect.Slice, reflect.Array:
		length := val.Len()
		maxLen := length
		if c.opts.MaxSliceMapLen > 0 && length > c.opts.MaxSliceMapLen {
			maxLen = c.opts.MaxSliceMapLen
		}

		for i := range maxLen {
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

		if maxLen < length {
			nodes = append(nodes, InspectNode{
				Name:  "...",
				Type:  "truncated",
				Value: fmt.Sprintf("%d items truncated", length-maxLen),
			})
		}
	case reflect.Map:
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
		})

		length := len(keys)
		maxLen := length
		if c.opts.MaxSliceMapLen > 0 && length > c.opts.MaxSliceMapLen {
			maxLen = c.opts.MaxSliceMapLen
		}

		for _, key := range keys[:maxLen] {
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

		if maxLen < length {
			nodes = append(nodes, InspectNode{
				Name:  "...",
				Type:  "truncated",
				Value: fmt.Sprintf("%d items truncated", length-maxLen),
			})
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

		typeStr := val.Type().String()
		c.stats.FieldTypeCount[typeStr]++
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
