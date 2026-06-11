package claimserializer

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
	"github.com/transactrx/NCPDPSerDe/pkg/serde"
)

// Serialize claim
func Serialize[V any](claim *V) (string, error) {
	if claim == nil {
		return serde.Empty, fmt.Errorf("'claim' parameter is nil")
	}

	claimType := reflect.TypeOf(claim)
	claimObjectRef := reflect.ValueOf(claim)

	builder := strings.Builder{}

	// Build header
	rawHdr, err := buildHeader(claimType, claimObjectRef)
	if err != nil {
		return serde.Empty, err
	}
	builder.WriteString(rawHdr)

	// Build shared segments
	rawSegs, err := buildSegments(claimType, claimObjectRef)
	if err != nil {
		return serde.Empty, err
	}
	builder.WriteString(rawSegs)

	// Build groups. vEB+ versions (e.g. F6) eliminated the group separator and
	// allow only a single transaction per transmission.
	omitGroupSeparators := strings.HasPrefix(rawHdr, "F")

	rawGrps, err := buildGroups(claimType, claimObjectRef, omitGroupSeparators)
	if err != nil {
		return serde.Empty, err
	}
	builder.WriteString(rawGrps)

	return builder.String(), nil
}

// Build raw header.
func buildHeader(structType reflect.Type, structVal reflect.Value) (string, error) {
	headerField, attr, err := serde.GetHeaderField(structType)
	if err != nil {
		return serde.Empty, err
	}

	if headerField == nil {
		return serde.Empty, fmt.Errorf("header is null")
	}

	if attr.Serializer == serde.Empty {
		return serde.Empty, fmt.Errorf("serializer undefined for header")
	}

	// Get header field object reference
	baseVal := reflectionutils.GetElementValue(structVal)
	headerValue := baseVal.FieldByName(headerField.Name)

	serMethod, err := serde.GetMethod(headerValue, attr.Serializer)
	if err != nil {
		return serde.Empty, err
	}

	serMethod.Call([]reflect.Value{})

	rawField := headerValue.FieldByName(attr.RawField)
	if !rawField.IsValid() {
		return serde.Empty, fmt.Errorf("raw field not found in header definition")
	}

	return rawField.String(), nil
}

// Build raw groups.
func buildGroups(structType reflect.Type, structVal reflect.Value, omitGroupSeparators bool) (string, error) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	// Find the slice field with the group tag (if present)
	groupField, err := serde.GetGroupSlice(baseType)
	if err != nil {
		return serde.Empty, nil
	}

	if groupField == nil {
		return serde.Empty, nil
	}

	groupSlice := baseVal.FieldByName(groupField.Name)

	if omitGroupSeparators && groupSlice.Len() > 1 {
		return serde.Empty, fmt.Errorf("NCPDP versions without group separators (vEB and higher) allow a single transaction per transmission; found %d claim groups", groupSlice.Len())
	}

	builder := strings.Builder{}

	for i := 0; i < groupSlice.Len(); i++ {
		groupElement := groupSlice.Index(i)

		rawSegs, err := buildSegments(groupElement.Type(), groupElement)
		if err != nil {
			return serde.Empty, err
		}

		if !omitGroupSeparators {
			builder.WriteByte(ncpdp.GROUP)
		}
		builder.WriteString(rawSegs)
	}

	return builder.String(), nil
}

// Build raw segments.
func buildSegments(structType reflect.Type, structVal reflect.Value) (string, error) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	// Get segment info using the order as the key
	segmentMap, err := serde.GetSegmentDefinitionByOrder(baseType)
	if err != nil {
		return serde.Empty, err
	}

	// Sort segmentKeys by order
	segmentKeys := make([]int, 0, len(segmentMap))

	for k := range segmentMap {
		segmentKeys = append(segmentKeys, k)
	}
	slices.Sort(segmentKeys)

	builder := strings.Builder{}

	for _, segmentOrder := range segmentKeys {
		segmentDef := segmentMap[segmentOrder]
		segmentVal := baseVal.FieldByName(segmentDef.Field.Name)

		switch segmentVal.Kind() {
		case reflect.Slice:
			for i := 0; i < segmentVal.Len(); i++ {
				sliceItemVal := segmentVal.Index(i)

				raw, err := buildSegment(sliceItemVal.Type(), sliceItemVal)
				if err != nil {
					return serde.Empty, err
				}
				builder.WriteString(raw)
			}

		default:
			raw, err := buildSegment(segmentDef.Field.Type, segmentVal)
			if err != nil {
				return serde.Empty, err
			}
			builder.WriteString(raw)
		}
	}

	return builder.String(), nil
}

// Build individual raw segment.
func buildSegment(structType reflect.Type, structVal reflect.Value) (string, error) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	switch baseVal.Interface().(type) {
	case dynamic.DynamicStruct:
		ds, _ := baseVal.Interface().(dynamic.DynamicStruct)
		baseType = reflectionutils.GetElementType(ds.DynamicType)
		baseVal = reflectionutils.GetElementValue(reflect.ValueOf(ds.Value))
	}

	builder := strings.Builder{}
	builder.WriteByte(ncpdp.SEGMENT)

	// Build map of fields and contents
	rawFieldMap, err := buildFieldMap(baseType, baseVal, nil)
	if err != nil {
		return builder.String(), err
	}

	// Sort keys by order
	fieldKeys := make([]int, 0, len(rawFieldMap))

	for k := range rawFieldMap {
		fieldKeys = append(fieldKeys, k)
	}
	slices.Sort(fieldKeys)

	// Append values in order
	for _, fieldOrder := range fieldKeys {
		fieldVal, ok := rawFieldMap[fieldOrder]
		if ok {
			builder.WriteString(fieldVal)
		}
	}

	// Exclude empty segments
	rawSegment := builder.String()
	if len(rawSegment) <= 6 {
		return serde.Empty, nil
	}

	return rawSegment, nil
}

func buildFieldMap(structType reflect.Type, structVal reflect.Value, skipCodes map[string]bool) (map[int]string, error) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	switch baseVal.Interface().(type) {
	case dynamic.DynamicStruct:
		ds, _ := baseVal.Interface().(dynamic.DynamicStruct)
		baseType = reflectionutils.GetElementType(ds.DynamicType)
		baseVal = reflectionutils.GetElementValue(reflect.ValueOf(ds.Value))
	}

	// Codes already emitted by populated repeating slices suppress their scalar
	// counterparts so dual-mapped fields (D0 scalar + F6 repeating slice) are
	// not serialized twice.
	skipCodes = mergeSliceFieldCodes(baseType, baseVal, skipCodes)

	rawFieldMap := make(map[int]string)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)
		fieldVal := baseVal.FieldByName(field.Name)

		tag := field.Tag.Get(serde.FieldTag)
		fieldAttribute, err := serde.GetFieldAttribute(tag)
		if err != nil {
			return rawFieldMap, err
		}

		fieldResult, err := buildStructField(field.Type, fieldVal, &fieldAttribute, skipCodes)
		if err != nil {
			return rawFieldMap, err
		}

		for k, v := range fieldResult {
			rawFieldMap[k] = v
		}
	}

	return rawFieldMap, nil
}

// Union of skipCodes and the field codes contained in non-empty slice fields.
func mergeSliceFieldCodes(baseType reflect.Type, baseVal reflect.Value, skipCodes map[string]bool) map[string]bool {
	var merged map[string]bool

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		if field.Type.Kind() != reflect.Slice {
			continue
		}

		fieldVal := baseVal.FieldByName(field.Name)
		if !fieldVal.IsValid() || fieldVal.Len() == 0 {
			continue
		}

		if merged == nil {
			merged = make(map[string]bool, len(skipCodes))
			for k := range skipCodes {
				merged[k] = true
			}
		}

		collectFieldCodes(field.Type.Elem(), merged)
	}

	if merged == nil {
		return skipCodes
	}

	return merged
}

// Collect all field tag codes defined on a struct type, recursing into nested structs.
func collectFieldCodes(structType reflect.Type, codes map[string]bool) {
	structType = reflectionutils.GetElementType(structType)
	if structType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		tag := field.Tag.Get(serde.FieldTag)
		if tag != serde.Empty {
			attr, err := serde.GetFieldAttribute(tag)
			if err == nil && attr.Code != serde.Empty {
				codes[attr.Code] = true
			}
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			collectFieldCodes(field.Type, codes)
		case reflect.Pointer:
			if field.Type.Elem().Kind() == reflect.Struct {
				collectFieldCodes(field.Type.Elem(), codes)
			}
		}
	}
}

func buildStructField(ft reflect.Type, fieldVal reflect.Value, fieldAttribute *serde.FieldAttribute, skipCodes map[string]bool) (map[int]string, error) {
	switch ft.Kind() {
	case reflect.Struct:
		return buildFieldMap(ft, fieldVal, skipCodes)

	case reflect.Pointer:
		return buildPointerField(ft, fieldVal, fieldAttribute, skipCodes)

	case reflect.Slice, reflect.Array:
		return buildSliceField(fieldVal)

	default:
		return buildScalarField(fieldVal, fieldAttribute, skipCodes)
	}
}

func buildPointerField(ft reflect.Type, fieldVal reflect.Value, fieldAttribute *serde.FieldAttribute, skipCodes map[string]bool) (map[int]string, error) {
	if fieldAttribute == nil || fieldAttribute.Order <= 0 {
		return buildStructField(ft.Elem(), fieldVal.Elem(), fieldAttribute, skipCodes)
	}

	rawFieldMap := make(map[int]string)

	if fieldVal.IsNil() || skipCodes[fieldAttribute.Code] {
		return rawFieldMap, nil
	}

	rawField, err := buildField(*fieldAttribute, fieldVal.Elem())
	if err != nil {
		return rawFieldMap, err
	}

	rawFieldMap[fieldAttribute.Order] = rawField

	return rawFieldMap, nil
}

func buildSliceField(fieldVal reflect.Value) (map[int]string, error) {
	rawFieldMap := make(map[int]string)

	builder := strings.Builder{}
	order := math.MaxInt

	for i := 0; i < fieldVal.Len(); i++ {
		element := fieldVal.Index(i)

		fm, err := buildFieldMap(element.Type(), element, nil)
		if err != nil {
			return rawFieldMap, err
		}

		elementOrder := writeSortedFieldMap(&builder, fm)
		if elementOrder < order {
			order = elementOrder
		}
	}

	rawString := builder.String()
	if rawString != serde.Empty {
		rawFieldMap[order] = rawString
	}

	return rawFieldMap, nil
}

// writeSortedFieldMap appends the map values to the builder in order and
// returns the lowest order present (math.MaxInt when the map is empty).
func writeSortedFieldMap(builder *strings.Builder, fm map[int]string) int {
	fmKeys := make([]int, 0, len(fm))
	for k := range fm {
		fmKeys = append(fmKeys, k)
	}

	slices.Sort(fmKeys)

	for _, fmOrder := range fmKeys {
		builder.WriteString(fm[fmOrder])
	}

	if len(fmKeys) == 0 {
		return math.MaxInt
	}

	return fmKeys[0]
}

func buildScalarField(fieldVal reflect.Value, fieldAttribute *serde.FieldAttribute, skipCodes map[string]bool) (map[int]string, error) {
	rawFieldMap := make(map[int]string)

	if fieldAttribute == nil || fieldAttribute.Order <= 0 || skipCodes[fieldAttribute.Code] {
		return rawFieldMap, nil
	}

	rawField, err := buildField(*fieldAttribute, fieldVal)
	if err != nil {
		return rawFieldMap, err
	}

	rawFieldMap[fieldAttribute.Order] = rawField

	return rawFieldMap, nil
}

// Build raw field.
func buildField(fieldAttr serde.FieldAttribute, structVal reflect.Value) (string, error) {
	baseVal := reflectionutils.GetElementValue(structVal)

	builder := strings.Builder{}

	switch baseVal.Interface().(type) {
	case int, int8, int16, int32, int64:
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString(fieldAttr.Code)
		builder.WriteString(strconv.FormatInt(int64(baseVal.Int()), 10))

	case uint, uint8, uint16, uint32, uint64:
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString(fieldAttr.Code)
		builder.WriteString(strconv.FormatUint(uint64(baseVal.Uint()), 10))

	case float32, float64:
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString(fieldAttr.Code)
		fs := serde.NewFieldSettings(fieldAttr)
		rawVal := ""
		if fs.Overpunch {
			rawVal = fs.Sign(baseVal.Interface())
		} else {
			rawVal = fs.ToImpliedDecimalString(baseVal.Interface())
		}
		builder.WriteString(rawVal)

	case time.Time:
		if baseVal.IsZero() {
			return serde.Empty, nil
		}
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString(fieldAttr.Code)
		t := baseVal.Interface().(time.Time)
		builder.WriteString(t.Format(fieldAttr.Format))

	default:
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString(fieldAttr.Code)
		builder.WriteString(baseVal.String())
	}

	return builder.String(), nil
}
