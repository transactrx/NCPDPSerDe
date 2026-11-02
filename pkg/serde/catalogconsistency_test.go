package serde_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/serde"
)

// These tests hold the model catalog itself to the invariants the parser and
// serializer rely on. Each one is a silent-data-loss vector when broken:
//
//   - a transaction or group type without a dynamic segment slice cannot
//     preserve unknown or duplicated segments (the deserializer's append is
//     best-effort at that level);
//   - a segment or repeating-item type without a dynamic field slice loses
//     the values its typed fields cannot represent, or misattributes them to
//     the wrong occurrence;
//   - a countfor tag naming a missing or non-slice sibling never derives its
//     count, so rebuilt messages lose the counter;
//   - two fields sharing one code within the same version scope both consume
//     the value, duplicating it on rebuild;
//   - two fields sharing one order collide in the serializer's order-keyed
//     field map, silently dropping one of them.
//
// They walk every registered transaction type, so a future segment or F6
// model that breaks an invariant fails here instead of in production traffic.

// structural codes that carry no order and may repeat per struct.
var structuralCodes = map[string]bool{"dynamic": true, "rawsegment": true, "rawgroup": true}

// elementStruct unwraps pointers and slices down to a struct type, or nil.
func elementStruct(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct && t.String() != "time.Time" {
		return t
	}
	return nil
}

// hasDynamicSlice mirrors the deserializer's findDynamicSliceValue at the
// type level: the dynamic slice may live on the type itself or inside a
// nested (embedded) struct, but never inside a nested slice.
func hasDynamicSlice(structType reflect.Type, tagType string) bool {
	if structType.Kind() != reflect.Struct {
		return false
	}

	if df, err := serde.GetDynamicSlice(structType, tagType); err == nil && df != nil {
		return true
	}

	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Type.Kind() == reflect.Struct &&
			hasDynamicSlice(structType.Field(i).Type, tagType) {
			return true
		}
	}

	return false
}

// catalogRoles indexes every struct type reachable from the registered
// transaction types by its structural role, with one example path each.
type catalogRoles struct {
	parts    map[reflect.Type]string // transaction roots and group element types
	segments map[reflect.Type]string // segment-tagged struct fields
	items    map[reflect.Type]string // elements of repeating slices inside segments
}

func collectRoles(t *testing.T) *catalogRoles {
	t.Helper()

	roles := &catalogRoles{
		parts:    map[reflect.Type]string{},
		segments: map[reflect.Type]string{},
		items:    map[reflect.Type]string{},
	}

	for key, tranType := range serde.TransactionTypes() {
		roles.collectPart(tranType, key)
	}

	return roles
}

func (roles *catalogRoles) collectPart(partType reflect.Type, path string) {
	if _, seen := roles.parts[partType]; seen {
		return
	}
	roles.parts[partType] = path

	for i := 0; i < partType.NumField(); i++ {
		field := partType.Field(i)

		if tag, found := field.Tag.Lookup("segment"); found {
			if strings.Contains(tag, "dynamic") {
				continue // the DynamicSegments slice itself
			}
			if segType := elementStruct(field.Type); segType != nil {
				if _, seen := roles.segments[segType]; !seen {
					roles.segments[segType] = path + "." + field.Name
					roles.collectItems(segType, path+"."+field.Name)
				}
			}
			continue
		}

		if _, found := field.Tag.Lookup("group"); found {
			if el := elementStruct(field.Type); el != nil {
				roles.collectPart(el, path+"."+field.Name)
			}
		}
	}
}

// collectItems finds repeating-item element types inside a segment subtree:
// any slice of structs other than the dynamic.DynamicStruct slices, at any
// nesting depth (items nest, e.g. amounts inside an other-payer occurrence).
func (roles *catalogRoles) collectItems(structType reflect.Type, path string) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		switch field.Type.Kind() {
		case reflect.Slice:
			el := elementStruct(field.Type)
			if el == nil || el == reflect.TypeOf(dynamic.DynamicStruct{}) {
				continue
			}
			if _, seen := roles.items[el]; !seen {
				roles.items[el] = path + "." + field.Name
				roles.collectItems(el, path+"."+field.Name)
			}

		case reflect.Struct:
			roles.collectItems(field.Type, path+"."+field.Name)
		}
	}
}

// Every transaction root and group element must be able to hold preserved
// segments; every segment and repeating item must be able to hold preserved
// fields. Without the slice, preservation fails and the parse either errors
// on well-formed-enough input (segments) or misattributes data (items).
func TestEveryCatalogTypeCanPreserveUnconsumedData(t *testing.T) {
	roles := collectRoles(t)

	for partType, where := range roles.parts {
		if !hasDynamicSlice(partType, serde.SegmentTag) {
			t.Errorf("%s (%s): no dynamic segment slice reachable — unknown or duplicate segments in this part cannot be preserved", partType, where)
		}
	}

	for segType, where := range roles.segments {
		if !hasDynamicSlice(segType, serde.FieldTag) {
			t.Errorf("%s (%s): no dynamic field slice reachable — unconsumed fields on this segment make the parse fail loudly instead of preserving", segType, where)
		}
	}

	for itemType, where := range roles.items {
		if !hasDynamicSlice(itemType, serde.FieldTag) {
			t.Errorf("%s (%s): repeating item type has no dynamic field slice — an unrepresentable value in one occurrence would be preserved at segment level, losing its occurrence attribution", itemType, where)
		}
	}
}

// taggedField is one field-tagged leaf found in a segment subtree.
type taggedField struct {
	path  string
	owner reflect.Type // the struct that declares the field
	attr  serde.FieldAttribute

	// dualMapped marks leaves inside an F6 dual-mapped composite (a slice
	// whose own tag is code-less with sinceVersion set, like Patient.Ids):
	// their codes and orders deliberately shadow the D0 scalars they repeat,
	// and the serializer reconciles the two (mergeSliceFieldCodes), so they
	// are exempt from the duplicate-code and duplicate-order checks.
	dualMapped bool
}

// collectTaggedFields walks a segment subtree — embedded structs share the
// segment's order space, and repeating items contribute their own orders to
// it too (orders are globally unique per segment in the NCPDP models).
func collectTaggedFields(t *testing.T, structType reflect.Type, path string, dualMapped bool, visited map[reflect.Type]bool, out *[]taggedField) {
	t.Helper()

	if visited[structType] {
		return
	}
	visited[structType] = true

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldPath := path + "." + field.Name
		elementDualMapped := dualMapped

		if tag, found := field.Tag.Lookup(serde.FieldTag); found {
			attr, err := serde.GetFieldAttribute(tag)
			if err != nil {
				t.Errorf("%s: field tag does not parse: %v", fieldPath, err)
				continue
			}
			if attr.Code == serde.Empty && attr.SinceVersion != serde.Empty && field.Type.Kind() == reflect.Slice {
				elementDualMapped = true
			}
			if attr.Code != serde.Empty {
				*out = append(*out, taggedField{path: fieldPath, owner: structType, attr: attr, dualMapped: dualMapped})
			}
		}

		if el := elementStruct(field.Type); el != nil && el != reflect.TypeOf(dynamic.DynamicStruct{}) {
			collectTaggedFields(t, el, fieldPath, elementDualMapped, visited, out)
		}
	}
}

func TestCatalogFieldTagInvariants(t *testing.T) {
	roles := collectRoles(t)

	for segType, where := range roles.segments {
		fields := []taggedField{}
		collectTaggedFields(t, segType, where, false, map[reflect.Type]bool{}, &fields)

		orders := map[int][]string{}
		codes := map[string][]string{}

		for _, f := range fields {
			if structuralCodes[f.attr.Code] {
				continue
			}

			// countfor must name a sibling slice on the declaring struct (or
			// be the reserved per-item index counter): the serializer derives
			// the count from that slice's length via FieldByName.
			if f.attr.CountFor != serde.Empty && f.attr.CountFor != serde.CountForIndex {
				target, found := f.owner.FieldByName(f.attr.CountFor)
				if !found || target.Type.Kind() != reflect.Slice {
					t.Errorf("%s: countfor=%q does not name a sibling slice field — the counter can never derive", f.path, f.attr.CountFor)
				}
			}

			if f.dualMapped {
				continue
			}

			orders[f.attr.Order] = append(orders[f.attr.Order], f.path)
			codeKey := f.attr.Code + "|" + f.attr.SinceVersion
			codes[codeKey] = append(codes[codeKey], f.path)
		}

		for order, paths := range orders {
			if len(paths) > 1 {
				t.Errorf("%s: order %d is shared by %v — the serializer's order-keyed field map would silently drop all but one", segType, order, paths)
			}
		}

		for codeKey, paths := range codes {
			if len(paths) > 1 {
				t.Errorf("%s: code %s is mapped by %v within the same version scope — the deserializer would consume one received value into both", segType, codeKey, paths)
			}
		}
	}
}
