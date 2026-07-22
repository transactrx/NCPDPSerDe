package dynamic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type DynamicStruct struct {
	DynamicType reflect.Type `json:"-"`
	Value       *any
}

// Field name of the raw-data field included in dynamic segment types.
const rawFieldName = "Raw"

// Field name prefix for dynamically generated fields.
const fieldNamePrefix = "Field_"

var illegalNamingRunes = map[rune]string{
	'!':  "exclamation",
	'"':  "quote",
	'#':  "hash",
	'$':  "dollar",
	'&':  "ampersand",
	'\'': "apostrophe",
	'(':  "leftparenthesis",
	')':  "rightparenthesis",
	'*':  "asterisk",
	',':  "comma",
	':':  "colon",
	';':  "semicolon",
	'<':  "lessthan",
	'=':  "equals",
	'>':  "greaterthan",
	'?':  "questionmark",
	'[':  "leftbracket",
	'\\': "backslash",
	']':  "rightbracket",
	'^':  "caret",
	'`':  "backtick",
	'{':  "leftbrace",
	'|':  "pipe",
	'}':  "rightbrace",
}

// FieldName builds the struct field name for a dynamically generated field.
func FieldName(fieldCode string, order int) string {
	return fieldNamePrefix + ReplaceIllegalCharacters(fieldCode) + "_" + strconv.Itoa(order)
}

// ReplaceIllegalCharacters replaces characters that are invalid in Go
// identifiers with word equivalents.
func ReplaceIllegalCharacters(str string) string {
	builder := strings.Builder{}

	for _, ch := range str {
		val, found := illegalNamingRunes[ch]
		if found {
			builder.WriteString(val)
		} else {
			builder.WriteRune(ch)
		}
	}

	return builder.String()
}

// restoreIllegalCharacters reverses ReplaceIllegalCharacters by matching the
// longest word equivalent at each position.
func restoreIllegalCharacters(mangled string) string {
	builder := strings.Builder{}

	for len(mangled) > 0 {
		matchedRune, matchedLen := rune(0), 0
		for ch, word := range illegalNamingRunes {
			if len(word) > matchedLen && strings.HasPrefix(mangled, word) {
				matchedRune, matchedLen = ch, len(word)
			}
		}

		if matchedLen > 0 {
			builder.WriteRune(matchedRune)
			mangled = mangled[matchedLen:]
			continue
		}

		r := []rune(mangled)
		builder.WriteRune(r[0])
		mangled = string(r[1:])
	}

	return builder.String()
}

// parseFieldName extracts the field code and order encoded in a dynamically
// generated field name (see FieldName).
func parseFieldName(name string) (code string, order int, ok bool) {
	if !strings.HasPrefix(name, fieldNamePrefix) {
		return "", 0, false
	}

	rest := name[len(fieldNamePrefix):]
	sep := strings.LastIndex(rest, "_")
	if sep < 0 {
		return "", 0, false
	}

	order, err := strconv.Atoi(rest[sep+1:])
	if err != nil {
		return "", 0, false
	}

	return restoreIllegalCharacters(rest[:sep]), order, true
}

// UnmarshalJSON rebuilds DynamicType and a typed Value from JSON produced by
// marshaling a DynamicStruct. DynamicType is excluded from the JSON (json:"-"),
// so it is reconstructed from the field names, which encode each field's code
// and order. Without this, a JSON round-trip leaves DynamicType nil and the
// struct cannot be serialized back to NCPDP format.
func (d *DynamicStruct) UnmarshalJSON(data []byte) error {
	var wrapper struct {
		Value json.RawMessage
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	d.DynamicType = nil
	d.Value = nil

	if len(wrapper.Value) == 0 || string(wrapper.Value) == "null" {
		return nil
	}

	var rawFields map[string]*string
	if err := json.Unmarshal(wrapper.Value, &rawFields); err != nil {
		// Not the flat string object a dynamic segment/field marshals to;
		// preserve the value untyped.
		var v any
		if err := json.Unmarshal(wrapper.Value, &v); err != nil {
			return err
		}
		d.Value = &v
		return nil
	}

	type fieldEntry struct {
		structField reflect.StructField
		value       *string
		order       int
	}

	entries := []fieldEntry{}

	for name, value := range rawFields {
		if name == rawFieldName {
			entries = append(entries, fieldEntry{
				structField: reflect.StructField{
					Name: rawFieldName,
					Type: reflect.TypeOf((*string)(nil)),
					Tag:  `field:"code=rawsegment"`,
				},
				value: value,
				order: 0,
			})
			continue
		}

		code, order, ok := parseFieldName(name)
		if !ok {
			continue
		}

		entries = append(entries, fieldEntry{
			structField: reflect.StructField{
				Name: name,
				Type: reflect.TypeOf((*string)(nil)),
				Tag:  reflect.StructTag(fmt.Sprintf(`field:"code=%v,order=%v"`, code, order)),
			},
			value: value,
			order: order,
		})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].order < entries[j].order })

	structFields := make([]reflect.StructField, 0, len(entries))
	for _, entry := range entries {
		structFields = append(structFields, entry.structField)
	}

	dynamicType := reflect.StructOf(structFields)
	dynamicValue := reflect.New(dynamicType).Elem()

	for i, entry := range entries {
		if entry.value != nil {
			dynamicValue.Field(i).Set(reflect.ValueOf(entry.value))
		}
	}

	iface := dynamicValue.Interface()
	d.DynamicType = dynamicType
	d.Value = &iface

	return nil
}

type DynamicFieldData struct {
	Name    string
	Type    reflect.Type
	TagName string
	Tags    map[string]string
}

type DynamicFieldList []DynamicFieldData

// Convert to struct field.
func (d DynamicFieldData) ToStructField() reflect.StructField {
	tags := []string{}

	for k, v := range d.Tags {
		tags = append(tags, fmt.Sprintf("%v=%v", k, v))
	}

	sTag := ""
	if len(tags) > 0 {
		sTag = fmt.Sprintf(`%v:"%v"`, d.TagName, strings.Join(tags[:], ","))
	}

	sf := reflect.StructField{
		Name: d.Name,
		Type: d.Type,
		Tag:  reflect.StructTag(sTag),
	}

	return sf
}

func (dfl DynamicFieldList) ToStructFieldList() []reflect.StructField {
	structFields := []reflect.StructField{}

	for _, v := range dfl {
		structFields = append(structFields, v.ToStructField())
	}

	return structFields
}
