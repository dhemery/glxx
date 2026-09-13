package describe

import (
	"fmt"
	"maps"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type Properties map[string]PropertyValue

// Describe describes the properties if the report requests them.
func (p Properties) Describe(r Report) {
	if !r.IncludeAllProperties || len(p) == 0 {
		return
	}

	r.Line()
	for label, value := range p {
		value.Describe(label, r)
	}
}

type PropertyValue interface {
	Describe(label string, r Report)
	String() string
}

type StringValue string

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Describe(label string, r Report) {
	r.Item(label, v.String())
}

type PropertyMultiValue []PropertyValue

func (v PropertyMultiValue) String() string {
	if len(v) == 0 {
		return ""
	}
	return v[0].String()
}

func (v PropertyMultiValue) Describe(label string, r Report) {
	for _, val := range v {
		val.Describe(label, r)
	}
}

type PropertyObjectValue struct {
	Value  PropertyValue
	Date   string
	Fields map[string]string
}

func (v PropertyObjectValue) String() string {
	return ""
}

func (v PropertyObjectValue) Describe(label string, r Report) {
	v.Value.Describe(label, r)
	if v.Date != "" {
		r.Item("  Date", v.Date)
	}
	for _, flabel := range slices.Sorted(maps.Keys(v.Fields)) {
		r.Item("  "+flabel, v.Fields[flabel])
	}
}

type PropertyReferenceValue struct {
	DisplayName string
	ID          string
}

func (v PropertyReferenceValue) String() string {
	return v.DisplayName
}

func (v PropertyReferenceValue) Describe(label string, r Report) {
	r.Item(label, v.String())
	r.Item("  id", v.ID)
}

type PropertyErrorValue struct {
	Description string
	Value       string
}

func (v PropertyErrorValue) String() string {
	return v.Description
}

func (v PropertyErrorValue) Describe(label string, r Report) {
	r.Item(label, v.String())
	r.Item("  value", v.Value)

}

func NewProperties(vals map[string]any, defs map[string]*glx.PropertyDefinition, archive Archive) Properties {
	out := Properties{}
	for _, name := range slices.Sorted(maps.Keys(vals)) {
		def := defs[name]
		val := vals[name]
		out[def.Label] = NewPropertyValue(val, def, archive)
	}
	return out

}

func NewPropertyValue(value any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	if isMultiValueProperty(def) {
		return NewMultiValueProperty(asArray(value), def, archive)
	}

	return NewSingleValueProperty(value, def, archive)
}

func NewMultiValueProperty(vals []any, def *glx.PropertyDefinition, archive Archive) PropertyMultiValue {
	var out PropertyMultiValue

	for _, val := range vals {
		out = append(out, NewSingleValueProperty(val, def, archive))
	}

	return out
}

func NewSingleValueProperty(v any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	if isObjectProperty(def) {
		return NewObjectPropertyValue(asObject(v), def, archive)
	}

	return NewStringPropertyValue(v, def, archive)

}

func NewObjectPropertyValue(object map[string]any, def *glx.PropertyDefinition, archive Archive) PropertyObjectValue {
	var out PropertyObjectValue
	out.Value = NewStringPropertyValue(object["value"], def, archive)
	if date, ok := object["date"]; ok {
		out.Date = fmt.Sprint(date)
	}
	out.Fields = make(map[string]string)
	if fields, ok := object["fields"].(map[string]any); ok {
		fdefs := def.Fields
		for fname, fval := range fields {
			flabel := fdefs[fname].Label
			out.Fields[flabel] = fmt.Sprint(fval)
		}

	}
	return out
}

// NewStringPropertyValue constructs a PropertyValue assuming that value is neither an array nor a map.
func NewStringPropertyValue(value any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	switch {
	case def.ReferenceType != "":
		return NewPropertyReferenceValue(value, def, archive)

	case def.ValueType != "":
		return NewStringValue(value)

	case def.VocabularyType != "":
		return NewPropertyVocabularyValue(value, def, archive)

	default:
		return NewPropertyUnknownValue(value, def)
	}
}

func isMultiValueProperty(definition *glx.PropertyDefinition) bool {
	if definition.MultiValue != nil && *definition.MultiValue {
		return true
	}
	return isTemporalProperty(definition)
}

func isObjectProperty(definition *glx.PropertyDefinition) bool {
	return len(definition.Fields) > 0 || isTemporalProperty(definition)
}

func isTemporalProperty(definition *glx.PropertyDefinition) bool {
	return definition.Temporal != nil && *definition.Temporal
}

func asObject(value any) map[string]any {
	if m, ok := value.(map[string]any); ok {
		return m
	}
	return map[string]any{
		"value": value,
	}
}

func asArray(value any) []any {
	if v, ok := value.([]any); ok {
		return v
	}
	return []any{value}

}

func NewPropertyReferenceValue(value any, definition *glx.PropertyDefinition, archive Archive) PropertyReferenceValue {
	var out PropertyReferenceValue

	id := fmt.Sprint(value)

	switch definition.ReferenceType {
	case "persons":
		if person := archive.Person(id); person != nil {
			out.DisplayName = person.Name()
		}
	case "places":
		if place := archive.Place(id); place != nil {
			out.DisplayName = place.Name()
		}

	default:
		panic("NewPropertyReferenceValue" + definition.Label + " unimplemented reference type " + definition.ReferenceType)
	}

	return out
}

func NewStringValue(v any) StringValue {
	return StringValue(fmt.Sprint(v))
}

func NewPropertyVocabularyValue(v any, definition *glx.PropertyDefinition, archive Archive) PropertyValue {
	return NewStringValue(v)
}

func NewPropertyUnknownValue(value any, definition *glx.PropertyDefinition) PropertyErrorValue {
	return PropertyErrorValue{
		Description: "UNKNOWN " + definition.Label,
		Value:       fmt.Sprintf("%#v", value),
	}

}
