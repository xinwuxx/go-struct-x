package gostructx

type InspectReport struct {
	Root  InspectNode `json:"root"`
	Stats Stats       `json:"stats"`
}

type Stats struct {
	TotalFields    int            `json:"total_fields"`
	MaxDepth       int            `json:"max_depth"`
	CircularRef    int            `json:"circular_refs"`
	FieldTypeCount map[string]int `json:"field_type_count"`
}

type InspectNode struct {
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Value    any           `json:"value,omitempty"`
	Tag      string        `json:"tag,omitempty"`
	Children []InspectNode `json:"children,omitempty"`
}
