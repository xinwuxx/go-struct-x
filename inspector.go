package gostructx

import (
	"fmt"
	"reflect"
)

func Inspect(val any, opt ...Option) (InspectReport, error) {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return InspectReport{}, fmt.Errorf("invalid input value")
	}

	opts := defaultOptions()
	for _, o := range opt {
		o(opts)
	}

	c := &context{
		opts:    opts,
		visited: make(map[uintptr]bool),
		stats: Stats{
			FieldTypeCount: make(map[string]int),
		},
	}

	rootNode := InspectNode{
		Name:     "root",
		Type:     v.Type().String(),
		Children: c.inspectValue(v, 0),
	}

	report := InspectReport{
		Root:  rootNode,
		Stats: c.stats,
	}

	return report, nil
}
