// Package metaexport builds a machine-readable description of the
// version-agnostic NCPDP claim model: every registered transaction type,
// its JSON property names (the exported Go field names), and the NCPDP
// mapping metadata carried in the header/segment/group/field struct tags.
//
// The output allows non-Go consumers (e.g. the Java claim processor) to
// serialize the version-agnostic JSON format to raw NCPDP without
// duplicating the struct definitions by hand.
package metaexport

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/serde"
)

// Layout mirrors the fixed-width header layout tags (layout / layoutF6).
type Layout struct {
	Start  int `json:"start"`
	Length int `json:"length"`
	Order  int `json:"order"`
}

// Field describes one exported struct field.
//
// Role values:
//   - header:          the transaction header (Type is the header value type;
//     JSON shape is {Value}; RawValue/Size are internal and excluded from JSON)
//   - segment:         a segment struct (Code/Order from the segment tag)
//   - dynamicSegments: slice holding non-modeled segments
//   - group:           the repeating claim group slice (Max from the group tag)
//   - field:           a leaf NCPDP field (Code/Order/DecimalPlaces/Overpunch/Format)
//   - struct:          an untagged nested struct whose fields flatten into the parent
//   - repeating:       a slice of structs representing repeating fields
//   - dynamicFields:   slice holding non-modeled fields
//   - raw:             raw data captures (code=rawgroup/rawsegment); never
//     serialized and excluded from JSON
type Field struct {
	Name          string  `json:"name"`
	Role          string  `json:"role"`
	Type          string  `json:"type,omitempty"`
	GoType        string  `json:"goType,omitempty"`
	Code          string  `json:"code,omitempty"`
	Order         int     `json:"order,omitempty"`
	DecimalPlaces int     `json:"decimalPlaces,omitempty"`
	Overpunch     bool    `json:"overpunch,omitempty"`
	Format        string  `json:"format,omitempty"`
	Max           int     `json:"max,omitempty"`
	Layout        *Layout `json:"layout,omitempty"`
	LayoutF6      *Layout `json:"layoutF6,omitempty"`

	// CountFor marks an automatically derived counter: the sibling slice
	// field whose length this field reports, or "index" for a 1-based
	// position counter within a repeating item. Serializers should compute
	// these rather than require them as input.
	CountFor string `json:"countFor,omitempty"`
}

// Type describes one struct participating in the claim model.
type Type struct {
	Kind   string  `json:"kind"` // transaction | header | struct
	Fields []Field `json:"fields"`
}

// Metadata is the root of the exported model description.
type Metadata struct {
	// Transactions maps "<transactionCode>|request" / "<transactionCode>|response"
	// to the entry in Types describing the transaction.
	Transactions map[string]string `json:"transactions"`
	Types        map[string]*Type  `json:"types"`
}

const (
	headerTagName   = "header"
	segmentTagName  = "segment"
	groupTagName    = "group"
	fieldTagName    = "field"
	layoutTagName   = "layout"
	layoutF6TagName = "layoutF6"

	dynamicCode    = "dynamic"
	rawGroupCode   = "rawgroup"
	rawSegmentCode = "rawsegment"
)

var timeType = reflect.TypeOf(time.Time{})

// Build walks every registered transaction type and returns the model description.
func Build() (*Metadata, error) {
	meta := &Metadata{
		Transactions: map[string]string{},
		Types:        map[string]*Type{},
	}

	for key, tranType := range serde.TransactionTypes() {
		name, err := addType(meta, tranType, "transaction")
		if err != nil {
			return nil, fmt.Errorf("transaction %v: %w", key, err)
		}
		meta.Transactions[key] = name
	}

	return meta, nil
}

// Add a struct type (and everything it references) to the metadata.
// Returns the name the type is registered under.
func addType(meta *Metadata, structType reflect.Type, kind string) (string, error) {
	structType = baseType(structType)
	name := typeName(structType)

	if existing, ok := meta.Types[name]; ok {
		// Transactions may reuse struct types; keep the strongest kind.
		if kind == "transaction" {
			existing.Kind = kind
		}
		return name, nil
	}

	if structType.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %v", structType.Kind())
	}

	typeMeta := &Type{Kind: kind}
	meta.Types[name] = typeMeta

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		fieldMeta, err := buildFieldMeta(meta, field)
		if err != nil {
			return "", fmt.Errorf("%v.%v: %w", name, field.Name, err)
		}
		if fieldMeta != nil {
			typeMeta.Fields = append(typeMeta.Fields, *fieldMeta)
		}
	}

	return name, nil
}

// Build the description of a single struct field, or nil when the field
// does not participate in serialization.
func buildFieldMeta(meta *Metadata, field reflect.StructField) (*Field, error) {
	if tag := field.Tag.Get(headerTagName); tag != "" {
		return buildHeaderFieldMeta(meta, field)
	}

	if tag := field.Tag.Get(segmentTagName); tag != "" {
		return buildSegmentFieldMeta(meta, field, tag)
	}

	if tag := field.Tag.Get(groupTagName); tag != "" {
		return buildGroupFieldMeta(meta, field, tag)
	}

	if tag := field.Tag.Get(fieldTagName); tag != "" {
		return buildLeafFieldMeta(meta, field, tag)
	}

	// Untagged fields: nested structs flatten into the parent, slices repeat.
	switch baseType(field.Type).Kind() {
	case reflect.Struct:
		name, err := addType(meta, field.Type, "struct")
		if err != nil {
			return nil, err
		}
		return &Field{Name: field.Name, Role: "struct", Type: name}, nil

	case reflect.Slice:
		name, err := addType(meta, baseType(field.Type).Elem(), "struct")
		if err != nil {
			return nil, err
		}
		return &Field{Name: field.Name, Role: "repeating", Type: name}, nil
	}

	// Untagged scalars are not serialized.
	return nil, nil
}

func buildHeaderFieldMeta(meta *Metadata, field reflect.StructField) (*Field, error) {
	valueField, ok := baseType(field.Type).FieldByName("Value")
	if !ok {
		return nil, fmt.Errorf("header type %v has no Value field", field.Type)
	}

	name, err := addHeaderType(meta, valueField.Type)
	if err != nil {
		return nil, err
	}

	return &Field{Name: field.Name, Role: "header", Type: name}, nil
}

func buildSegmentFieldMeta(meta *Metadata, field reflect.StructField, tag string) (*Field, error) {
	attr, err := serde.GetSegmentAttribute(tag)
	if err != nil {
		return nil, err
	}

	if attr.Dynamic {
		return &Field{Name: field.Name, Role: "dynamicSegments", Order: attr.Order}, nil
	}

	segType := baseType(field.Type)
	if segType.Kind() == reflect.Slice {
		segType = segType.Elem()
	}

	name, err := addType(meta, segType, "struct")
	if err != nil {
		return nil, err
	}

	return &Field{
		Name:  field.Name,
		Role:  "segment",
		Type:  name,
		Code:  attr.Code,
		Order: attr.Order,
	}, nil
}

func buildGroupFieldMeta(meta *Metadata, field reflect.StructField, tag string) (*Field, error) {
	if baseType(field.Type).Kind() != reflect.Slice {
		return nil, fmt.Errorf("group tag on non-slice field")
	}

	max := 0
	if v := rawTagValue(tag, "max"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &max); err != nil {
			return nil, fmt.Errorf("invalid group max %q: %w", v, err)
		}
	}

	name, err := addType(meta, baseType(field.Type).Elem(), "struct")
	if err != nil {
		return nil, err
	}

	return &Field{Name: field.Name, Role: "group", Type: name, Max: max}, nil
}

func buildLeafFieldMeta(meta *Metadata, field reflect.StructField, tag string) (*Field, error) {
	attr, err := serde.GetFieldAttribute(tag)
	if err != nil {
		return nil, err
	}

	switch attr.Code {
	case dynamicCode:
		return &Field{Name: field.Name, Role: "dynamicFields"}, nil
	case rawGroupCode, rawSegmentCode:
		return &Field{Name: field.Name, Role: "raw", Code: attr.Code}, nil
	}

	return &Field{
		Name:          field.Name,
		Role:          "field",
		GoType:        goTypeName(field.Type),
		Code:          attr.Code,
		Order:         attr.Order,
		DecimalPlaces: attr.DecimalPlaces,
		Overpunch:     attr.Overpunch,
		// serde.GetFieldAttribute converts the format to a Go time layout;
		// export the readable tag value (e.g. YYYYMMdd) instead.
		Format:   rawTagValue(tag, "format"),
		CountFor: attr.CountFor,
	}, nil
}

// Add a fixed-width header value type (layout/layoutF6 tags).
func addHeaderType(meta *Metadata, headerType reflect.Type) (string, error) {
	headerType = baseType(headerType)
	name := typeName(headerType)

	if _, ok := meta.Types[name]; ok {
		return name, nil
	}

	typeMeta := &Type{Kind: "header"}
	meta.Types[name] = typeMeta

	for i := 0; i < headerType.NumField(); i++ {
		field := headerType.Field(i)

		layout, err := parseLayout(field.Tag.Get(layoutTagName))
		if err != nil {
			return "", fmt.Errorf("%v.%v: %w", name, field.Name, err)
		}
		layoutF6, err := parseLayout(field.Tag.Get(layoutF6TagName))
		if err != nil {
			return "", fmt.Errorf("%v.%v: %w", name, field.Name, err)
		}

		if layout == nil && layoutF6 == nil {
			continue
		}

		typeMeta.Fields = append(typeMeta.Fields, Field{
			Name:     field.Name,
			Role:     "field",
			GoType:   goTypeName(field.Type),
			Layout:   layout,
			LayoutF6: layoutF6,
		})
	}

	return name, nil
}

func parseLayout(tag string) (*Layout, error) {
	if tag == "" {
		return nil, nil
	}

	layout := Layout{}
	if _, err := fmt.Sscanf(tag, "start=%d,length=%d,order=%d", &layout.Start, &layout.Length, &layout.Order); err != nil {
		return nil, fmt.Errorf("invalid layout tag %q: %w", tag, err)
	}

	return &layout, nil
}

// Extract the raw value for a key from a comma-delimited struct tag.
func rawTagValue(tag, key string) string {
	for _, part := range strings.Split(tag, ",") {
		if value, found := strings.CutPrefix(part, key+"="); found {
			return value
		}
	}

	return ""
}

// Unwrap pointers.
func baseType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}

// Short package-qualified type name, e.g. "request.Billing".
func typeName(t reflect.Type) string {
	return strings.ReplaceAll(t.String(), "github.com/transactrx/NCPDPSerDe/pkg/", "")
}

// JSON-relevant scalar type name.
func goTypeName(t reflect.Type) string {
	t = baseType(t)

	if t == timeType {
		return "time"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float64"
	case reflect.Bool:
		return "bool"
	default:
		return t.Kind().String()
	}
}
