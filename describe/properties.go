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
	Definition  *glx.PropertyDefinition
	Value       any
}

func (v PropertyErrorValue) String() string {
	return "ERROR: " + v.Description
}

func (v PropertyErrorValue) Describe(label string, r Report) {
	r.Item(label, v.String())
	r.Item("  value", fmt.Sprintf("%#v", v.Value))
	r.Item("  definition", fmt.Sprintf("%v", v.Definition))

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

func NewPropertyValue(val any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	if isMultiValueProperty(def) {
		return NewMultiValueProperty(asArray(val), def, archive)
	}

	return NewSingleValuedProperty(val, def, archive)
}

func NewMultiValueProperty(vals []any, def *glx.PropertyDefinition, archive Archive) PropertyMultiValue {
	var out PropertyMultiValue

	for _, val := range vals {
		out = append(out, NewSingleValuedProperty(val, def, archive))
	}

	return out
}

func NewSingleValuedProperty(val any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	if isObjectProperty(def) {
		return NewObjectPropertyValue(asObject(val), def, archive)
	}

	return NewStringPropertyValue(fmt.Sprint(val), def, archive)

}

func NewObjectPropertyValue(object map[string]any, def *glx.PropertyDefinition, archive Archive) PropertyObjectValue {
	var out PropertyObjectValue

	if val, ok := object["value"]; ok {
		out.Value = NewStringPropertyValue(fmt.Sprint(val), def, archive)
	}

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

func NewStringPropertyValue(value string, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	switch {
	case def.ReferenceType != "":
		return NewPropertyReferenceValue(value, def, archive)

	case def.ValueType != "":
		return NewStringValue(value)

	case def.VocabularyType != "":
		return NewPropertyVocabularyValue(value, def, archive)

	default:
		return PropertyErrorValue{
			Description: "Not reference, value, or vocabulary type",
			Definition:  def,
			Value:       value,
		}
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

func asObject(val any) map[string]any {
	if m, ok := val.(map[string]any); ok {
		return m
	}
	return map[string]any{
		"value": val,
	}
}

func asArray(val any) []any {
	if v, ok := val.([]any); ok {
		return v
	}
	return []any{val}

}

func NewPropertyReferenceValue(val any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	id := fmt.Sprint(val)

	out := PropertyReferenceValue{ID: id}

	switch def.ReferenceType {
	case "persons":
		if person := archive.Person(id); person != nil {
			out.DisplayName = person.Name()
		}

	case "places":
		if place := archive.Place(id); place != nil {
			out.DisplayName = place.Name()
		}

	default:
		return PropertyErrorValue{
			Description: "Unimplemented reference type",
			Definition:  def,
			Value:       val,
		}
	}

	return out
}

func NewStringValue(v any) StringValue {
	return StringValue(fmt.Sprint(v))
}

func NewPropertyVocabularyValue(val any, def *glx.PropertyDefinition, archive Archive) PropertyValue {
	return NewStringValue(val)
}
