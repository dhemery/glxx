package describe

import (
	"fmt"
	"maps"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type Property struct {
	Value      any
	Definition *glx.PropertyDefinition
}

type Properties []Property

func NewProperties(values map[string]any, definitions map[string]*glx.PropertyDefinition) Properties {
	out := Properties{}
	for _, name := range slices.Sorted(maps.Keys(values)) {
		out = append(out, NewProperty(values[name], definitions[name]))
	}
	return out

}

// Describe describes the properties if the report requests them.
func (p Properties) Describe(r Report) {
	if !r.IncludeAllProperties || len(p) == 0 {
		return
	}

	r.Line()
	for _, property := range p {
		property.Describe(r)
	}
}

func NewProperty(value any, definition *glx.PropertyDefinition) Property {
	return Property{Value: value, Definition: definition}
}

func (p Property) Describe(r Report) {
	label := p.Definition.Label
	switch v := p.Value.(type) {
	case string:
		r.Item(label, v)
	case []any:
		describeMultivalue(v, p.Definition, r)
	case map[string]any:
		describePropertyObject(v, p.Definition, r)
	default:
		r.Item(label, fmt.Sprintf("unknown type %T", v))
	}
}

func describeMultivalue(values []any, definition *glx.PropertyDefinition, r Report) {
	label := definition.Label
	for _, value := range values {
		switch v := value.(type) {
		case string:
			r.Item(label, v)
		case map[string]any:
			describePropertyObject(v, definition, r)
		default:
			r.Item(label, fmt.Sprintf("unknown type %T", v))
		}
	}
}

func describePropertyObject(v map[string]any, definition *glx.PropertyDefinition, r Report) {
	label := definition.Label
	r.Item(label, fmt.Sprintf("%#v", v["value"]))

	if date, ok := v["date"]; ok {
		r.Item("  Date", fmt.Sprintf("%q", date))
	}
	fields, ok := v["fields"].(map[string]any)
	if !ok {
		return
	}
	for n, f := range fields {
		label := definition.Fields[n].Label
		r.Item("  "+label, fmt.Sprintf("%#v", f))

	}
}
