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

	// Build groups
	rawGrps, err := buildGroups(claimType, claimObjectRef)
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
func buildGroups(structType reflect.Type, structVal reflect.Value) (string, error) {
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

	builder := strings.Builder{}

	for i := 0; i < groupSlice.Len(); i++ {
		groupElement := groupSlice.Index(i)

		rawSegs, err := buildSegments(groupElement.Type(), groupElement)
		if err != nil {
			return serde.Empty, err
		}

		builder.WriteByte(ncpdp.GROUP)
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
	rawFieldMap, err := buildFieldMap(baseType, baseVal)
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

func buildFieldMap(structType reflect.Type, structVal reflect.Value) (map[int]string, error) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	switch baseVal.Interface().(type) {
	case dynamic.DynamicStruct:
		ds, _ := baseVal.Interface().(dynamic.DynamicStruct)
		baseType = reflectionutils.GetElementType(ds.DynamicType)
		baseVal = reflectionutils.GetElementValue(reflect.ValueOf(ds.Value))
	}

	rawFieldMap := make(map[int]string)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)
		fieldVal := baseVal.FieldByName(field.Name)

		tag := field.Tag.Get(serde.FieldTag)
		fieldAttribute, err := serde.GetFieldAttribute(tag)
		if err != nil {
			return rawFieldMap, err
		}

		fieldResult, err := buildStructField(field.Type, fieldVal, &fieldAttribute)
		if err != nil {
			return rawFieldMap, err
		}

		for k, v := range fieldResult {
			rawFieldMap[k] = v
		}
	}

	return rawFieldMap, nil
}

func buildStructField(ft reflect.Type, fieldVal reflect.Value, fieldAttribute *serde.FieldAttribute) (map[int]string, error) {
	rawFieldMap := make(map[int]string)

	switch ft.Kind() {
	case reflect.Struct:
		return buildFieldMap(ft, fieldVal)

	case reflect.Pointer:
		if fieldAttribute != nil && fieldAttribute.Order > 0 {
			if !fieldVal.IsNil() {
				rawField, err := buildField(*fieldAttribute, fieldVal.Elem())
				if err != nil {
					return rawFieldMap, err
				}

				rawFieldMap[fieldAttribute.Order] = rawField
			}

			break
		}

		return buildStructField(ft.Elem(), fieldVal.Elem(), fieldAttribute)

	case reflect.Slice, reflect.Array:
		builder := strings.Builder{}
		order := math.MaxInt

		for i := 0; i < fieldVal.Len(); i++ {
			element := fieldVal.Index(i)

			fm, err := buildFieldMap(element.Type(), element)
			if err != nil {
				return rawFieldMap, err
			}

			// Get keys and sort
			fmKeys := make([]int, 0, len(fm))
			for k := range fm {
				fmKeys = append(fmKeys, k)
			}

			slices.Sort(fmKeys)

			// Concat all fields for element
			for _, fmOrder := range fmKeys {
				fmVal, ok := fm[fmOrder]
				if ok {
					builder.WriteString(fmVal)

					if fmOrder < order {
						order = fmOrder
					}
				}
			}
		}

		rawString := builder.String()
		if rawString != serde.Empty {
			rawFieldMap[order] = rawString
		}

	default:
		if fieldAttribute != nil && fieldAttribute.Order > 0 {
			rawField, err := buildField(*fieldAttribute, fieldVal)
			if err != nil {
				return rawFieldMap, err
			}

			rawFieldMap[fieldAttribute.Order] = rawField
		}
	}

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
