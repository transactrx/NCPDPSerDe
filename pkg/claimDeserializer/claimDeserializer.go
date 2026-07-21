package claimdeserializer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
	"github.com/transactrx/NCPDPSerDe/pkg/serde"
	sliceutils "github.com/transactrx/NCPDPSerDe/pkg/sliceUtils"
	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

const (
	rawGroupTag   = "rawgroup"
	rawSegmentTag = "rawsegment"
)

const (
	dZeroRequestLength  = 56
	dZeroResponseLength = 31
	f6RequestLength     = 58
)

// Parse raw data.
// Good if you don't know anything about the data.
func Deserialize(rawClaimString string) (interface{}, error) {
	serde.RegisterTypes()

	rawClaimString = strings.TrimSpace(rawClaimString)

	// Determine type by first separator index
	firstSeparatorIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})
	if firstSeparatorIndex == -1 {
		return nil, fmt.Errorf("improperly formatted data: no NCPDP separators present")
	}

	// Assume request
	if firstSeparatorIndex == dZeroRequestLength || firstSeparatorIndex == f6RequestLength {
		return DeserializeRequest(rawClaimString)
	}

	// Assume response
	if firstSeparatorIndex == dZeroResponseLength {
		return DeserializeResponse(rawClaimString)
	}

	return nil, fmt.Errorf("unable to determine transaction type")
}

// Parse raw request.
// Good if you know it's a request type.
func DeserializeRequest(rawClaimString string) (interface{}, error) {
	serde.RegisterTypes()

	// Determine transaction type by tran code
	rawClaimString = strings.TrimSpace(rawClaimString)
	rawClaimString = strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), "")

	if len(rawClaimString) < 10 {
		return nil, fmt.Errorf("unable to determine transaction type")
	}

	// D0 request headers lead with the BIN (digits); F6 headers lead with the
	// version (e.g. "F6"), which moves the transaction code to bytes 2-4.
	tranCode := rawClaimString[8:10]
	if rawClaimString[0] == 'F' {
		tranCode = rawClaimString[2:4]
	}

	claimType, err := serde.GetRequestType(tranCode)
	if err != nil {
		return nil, err
	}

	// Create object
	claimObjectRef := reflect.New(claimType).Elem()
	initializeStruct(claimType, claimObjectRef)

	// Parse raw data into type
	return deserializeRaw(rawClaimString, claimType, claimObjectRef)
}

// Parse raw request.
// Good if you know it's a response type.
func DeserializeResponse(rawClaimString string) (interface{}, error) {
	serde.RegisterTypes()

	// Determine transaction type by tran code
	rawClaimString = strings.TrimSpace(rawClaimString)
	rawClaimString = strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), "")

	if len(rawClaimString) < 4 {
		return nil, fmt.Errorf("unable to determine transaction type")
	}

	tranCode := rawClaimString[2:4]

	claimType, err := serde.GetResponseType(tranCode)
	if err != nil {
		return nil, err
	}

	// Create object
	claimObjectRef := reflect.New(claimType).Elem()
	initializeStruct(claimType, claimObjectRef)

	// Parse raw data into type
	return deserializeRaw(rawClaimString, claimType, claimObjectRef)
}

// Parse raw claim.
// Good if you already know the type.
func DeserializeType[V any](rawClaimString string, claim *V) error {
	if claim == nil {
		return fmt.Errorf("'claim' parameter is nil")
	}

	claimType := reflect.TypeOf(claim)
	claimObjectRef := reflect.ValueOf(claim)

	_, err := deserializeRaw(rawClaimString, claimType, claimObjectRef)
	if err != nil {
		return err
	}

	return nil
}

// Parse raw data.
func deserializeRaw(rawClaimString string, claimType reflect.Type, claimObjectRef reflect.Value) (interface{}, error) {
	if strings.TrimSpace(rawClaimString) == serde.Empty {
		return nil, fmt.Errorf("NCPDP data is empty")
	}

	rawClaimString = strings.TrimSpace(strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), ""))

	firstSeparatorIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})
	if firstSeparatorIndex == -1 {
		return nil, fmt.Errorf("improperly formatted data: no NCPDP separators present")
	}

	// Set header
	evaluateHeader(rawClaimString, claimType, claimObjectRef)

	groupIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.GROUP})

	if groupIndex < 0 && strings.HasPrefix(rawClaimString, "F") {
		// vEB+ (e.g. F6) eliminated the group separator and allows only a single
		// transaction per transmission; the claim segments directly follow the
		// shared segments and are located by segment ID instead.
		err := evaluateUngroupedTransaction(rawClaimString, claimType, claimObjectRef)
		if err != nil {
			return nil, err
		}

		return claimObjectRef.Interface(), nil
	}

	// Set group data
	evaluateGrouping(rawClaimString, claimType, claimObjectRef)

	// Set shared segment data
	endIndex := groupIndex
	if endIndex <= 0 {
		endIndex = len(rawClaimString)
	}

	evaluateSegments(rawClaimString[0:endIndex], claimType, claimObjectRef)

	return claimObjectRef.Interface(), nil
}

// Find and parse header.
func evaluateHeader(rawClaimString string, structType reflect.Type, structVal reflect.Value) error {
	headerField, attr, err := serde.GetHeaderField(structType)

	if err != nil {
		return err
	}

	if headerField == nil {
		return nil
	}

	if attr.Deserializer == serde.Empty {
		return fmt.Errorf("deserializer undefined for header")
	}

	if attr.RawField == serde.Empty {
		return fmt.Errorf("raw field undefined for header")
	}

	// Get header field object reference
	baseVal := reflectionutils.GetElementValue(structVal)
	headerValue := baseVal.FieldByName(headerField.Name)

	rawField := headerValue.FieldByName(attr.RawField)
	if !rawField.IsValid() {
		return fmt.Errorf("raw field not found in header definition")
	}

	if rawField.CanSet() {
		rawField.Set(reflect.ValueOf(rawClaimString))
	}

	deserMethod, err := serde.GetMethod(headerValue, attr.Deserializer)
	if err != nil {
		return err
	}

	deserMethod.Call([]reflect.Value{})

	return nil
}

// Evaluate claim groups.
func evaluateGrouping(rawClaimString string, structType reflect.Type, structVal reflect.Value) error {
	rawGroups := stringutils.SplitBySeparator(rawClaimString, ncpdp.GROUP)

	// No groups present
	if len(rawGroups) == 0 {
		return nil
	}

	// Find the slice field with the group tag
	groupField, err := serde.GetGroupSlice(structType)
	if err != nil {
		return err
	}

	if groupField == nil {
		return nil
	}

	baseVal := reflectionutils.GetElementValue(structVal)
	groupSlice := baseVal.FieldByName(groupField.Name)

	for _, rawClaimGroup := range rawGroups {
		// Create new element for group slice
		groupElementType := groupField.Type.Elem()
		groupItem := reflect.New(groupElementType).Elem()

		// Initialize element
		initializeStruct(groupElementType, groupItem)

		// Set field defined as "raw" data tag
		setStructFieldByCodeTag(groupElementType, groupItem, rawGroupTag, reflect.ValueOf(string(ncpdp.GROUP)+rawClaimGroup), 0, 0, false)

		// Set segment data
		evaluateSegments(rawClaimGroup, groupElementType, groupItem)

		// Append element to group slice
		groupSlice.Set(reflect.Append(groupSlice, groupItem))
	}

	return nil
}

// Evaluate an ungrouped (vEB+/F6) transmission. No group separators exist, so
// the single transaction's segments are located by segment ID: the first
// segment that belongs to the claim group struct starts the transaction, and
// it plus everything after it parse as one claim group. Segments before it
// parse as shared segments.
func evaluateUngroupedTransaction(rawClaimString string, structType reflect.Type, structVal reflect.Value) error {
	groupField, err := serde.GetGroupSlice(structType)
	if err != nil {
		return err
	}

	transactionIndex, err := findTransactionIndex(rawClaimString, groupField)
	if err != nil {
		return err
	}

	sharedRaw := rawClaimString
	if transactionIndex != -1 {
		sharedRaw = rawClaimString[:transactionIndex]
	}

	err = evaluateSegments(sharedRaw, structType, structVal)
	if err != nil {
		return err
	}

	if transactionIndex == -1 {
		return nil
	}

	transactionRaw := rawClaimString[transactionIndex:]

	groupElementType := groupField.Type.Elem()
	groupItem := reflect.New(groupElementType).Elem()

	initializeStruct(groupElementType, groupItem)

	setStructFieldByCodeTag(groupElementType, groupItem, rawGroupTag, reflect.ValueOf(transactionRaw), 0, 0, false)

	err = evaluateSegments(transactionRaw, groupElementType, groupItem)
	if err != nil {
		return err
	}

	baseVal := reflectionutils.GetElementValue(structVal)
	groupSlice := baseVal.FieldByName(groupField.Name)
	groupSlice.Set(reflect.Append(groupSlice, groupItem))

	return nil
}

// findTransactionIndex locates the start of the single (vEB+) claim group: the
// first segment whose ID belongs to the group struct. Returns -1 when the
// struct has no group slice or no segment matches.
func findTransactionIndex(rawClaimString string, groupField *reflect.StructField) (int, error) {
	if groupField == nil {
		return -1, nil
	}

	groupSegmentMap, err := serde.GetSegmentDefinitionById(reflectionutils.GetElementType(groupField.Type.Elem()))
	if err != nil {
		return -1, err
	}

	segIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.SEGMENT})
	for segIndex != -1 {
		nextIndex := stringutils.IndexOfAny(rawClaimString, segIndex+1, []byte{ncpdp.SEGMENT})

		endIndex := len(rawClaimString)
		if nextIndex != -1 {
			endIndex = nextIndex
		}

		if _, ok := groupSegmentMap[segmentIdOf(rawClaimString[segIndex+1:endIndex])]; ok {
			return segIndex, nil
		}

		segIndex = nextIndex
	}

	return -1, nil
}

// segmentIdOf returns the segment ID field (e.g. "AM07") of a raw segment.
func segmentIdOf(rawSeg string) string {
	rawFields := stringutils.SplitBySeparator(rawSeg, ncpdp.FIELD)

	idIndex := sliceutils.IndexOfStartsWith(ncpdp.SEGMENT_FIELD_ID, rawFields)
	if idIndex < 0 {
		return serde.Empty
	}

	return rawFields[idIndex]
}

// Add elements to dynamic field list.
func appendDynamicFields(structType reflect.Type, structVal reflect.Value, dyFieldCodes, dyFieldValues []string) error {
	baseType := reflectionutils.GetElementType(structType)
	startOrder := baseType.NumField()

	for i := 0; i < len(dyFieldCodes); i++ {
		dyCode := dyFieldCodes[i]
		dyValue := dyFieldValues[i]

		// Set the order tag
		dyField := serde.NewDynamicStringPointerField(
			serde.DynamicFieldName(dyCode, startOrder),
			serde.FieldTag,
			map[string]string{serde.CodeTag: dyCode, serde.OrderTag: strconv.Itoa(startOrder)})

		// Convert field data to struct field.
		sf := dyField.ToStructField()

		//Create type and object reference.
		typ := reflect.StructOf([]reflect.StructField{sf})
		field := reflect.New(typ).Elem()

		// Initialize object reference
		initializeStruct(typ, field)

		// Set field value
		setStructFieldByCodeTag(typ, field, dyCode, reflect.ValueOf(dyValue), startOrder, 0, true)

		//Create item to append to slice
		ds := dynamic.DynamicStruct{DynamicType: typ}

		iface := field.Interface()
		ds.Value = &iface

		err := appendDynamicSlice(structType, structVal, reflect.ValueOf(ds), serde.FieldTag)
		if err != nil {
			return err
		}

		startOrder++
	}

	return nil
}

// Add element to dynamic segment list
func appendDynamicSlice(structType reflect.Type, structVal reflect.Value, dynamicElementValue reflect.Value, tagType string) error {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	dynamicField, err := serde.GetDynamicSlice(baseType, tagType)
	if err != nil {
		return err
	}

	if dynamicField == nil {
		return nil
	}

	dynamicObj := baseVal.FieldByName(dynamicField.Name)
	if !dynamicObj.IsValid() {
		return fmt.Errorf("dynamic object is invalid")
	}

	if !dynamicObj.CanSet() {
		return fmt.Errorf("dynamic object is read-only")
	}

	// Append element to slice
	dynamicObj.Set(reflect.Append(dynamicObj, dynamicElementValue))

	return nil
}

// Evaluate claim segments.
func evaluateSegments(rawData string, structType reflect.Type, structVal reflect.Value) error {
	rawSegments := stringutils.SplitBySeparator(rawData, ncpdp.SEGMENT)

	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	segmentMap, err := serde.GetSegmentDefinitionById(baseType)
	if err != nil {
		return err
	}

	for _, rawSeg := range rawSegments {
		// Get all fields
		rawFields := stringutils.SplitBySeparator(rawSeg, ncpdp.FIELD)

		// serde.Empty segment, skip
		if len(rawFields) == 0 {
			return nil
		}

		// Find segment ID in field list
		segmentIdIndex := sliceutils.IndexOfStartsWith(ncpdp.SEGMENT_FIELD_ID, rawFields)
		if segmentIdIndex < 0 {
			return fmt.Errorf("segment id element does not exist")
		}

		// Create new segment and initialize
		var elementType reflect.Type
		segmentIdFieldValue := rawFields[segmentIdIndex]
		segmentDef, ok := segmentMap[segmentIdFieldValue]
		dynamicSegment := dynamic.DynamicStruct{}
		isDynamicSegment := false

		if ok {
			elementType = segmentDef.Field.Type
		} else {
			// Segment not found in object definition, create dynamic type
			elementType = serde.NewDynamicSegmentType(rawSeg)
			dynamicSegment.DynamicType = elementType
			isDynamicSegment = true
		}

		segment := reflect.New(elementType).Elem()

		initializeStruct(elementType, segment)

		// Set field defined as "raw" data tag
		setStructFieldByCodeTag(elementType, segment, rawSegmentTag, reflect.ValueOf(string(ncpdp.SEGMENT)+rawSeg), 0, 0, false)

		fieldCountMap := map[string]int{}

		// Set field data
		dyFieldCodes := []string{}
		dyFieldValues := []string{}
		for fieldIndex := 0; fieldIndex < len(rawFields); fieldIndex++ {
			rawField := rawFields[fieldIndex]
			fieldCode := stringutils.Substring(rawField, 0, ncpdp.ID_LENGTH)
			fieldValue := stringutils.Substring(rawField, ncpdp.ID_LENGTH, -1)

			fieldCodeCount := fieldCountMap[fieldCode]

			setCount := setStructFieldByCodeTag(elementType, segment, fieldCode, reflect.ValueOf(fieldValue), fieldIndex+1, fieldCodeCount, isDynamicSegment)
			if setCount == 0 {
				dyFieldCodes = append(dyFieldCodes, fieldCode)
				dyFieldValues = append(dyFieldValues, fieldValue)
			}

			// Increment field code count
			fieldCodeCount++
			fieldCountMap[fieldCode] = fieldCodeCount
		}

		if isDynamicSegment {
			segmentIf := segment.Interface()
			dynamicSegment.Value = &segmentIf
			appendDynamicSlice(baseType, baseVal, reflect.ValueOf(dynamicSegment), serde.SegmentTag)
		} else {
			appendDynamicFields(segment.Type(), segment, dyFieldCodes, dyFieldValues)
			segmentRef := baseVal.FieldByName(segmentDef.Field.Name)
			segmentRef.Set(segment)
		}
	}

	return nil
}

// Set struct field using code tag value to locate appropriate field.
// Ex: `field:"code=EN"`
func setStructFieldByCodeTag(structType reflect.Type, structVal reflect.Value, codeTag string, propVal reflect.Value, orderTag, repeatingFieldIndex int, isDynamic bool) int {
	if structVal.Kind() != reflect.Struct {
		return 0
	}

	setCount := 0

	for i := 0; i < structVal.NumField(); i++ {
		f := structVal.Field(i)
		ft := structType.Field(i)

		switch ft.Type.Kind() {
		case reflect.Struct:
			// Embedded struct
			count := setStructFieldByCodeTag(ft.Type, f, codeTag, propVal, orderTag, repeatingFieldIndex, isDynamic)
			setCount = setCount + count

		case reflect.Slice:
			sliceLen := f.Len()
			sliceItemType := f.Type().Elem()

			// Create new slice element. Only grow this slice when the element
			// type owns the field code directly; a repeating code that lives in
			// a deeper nested slice (e.g. a second 6E reject within a single
			// other payer) must repeat inside the last element instead of
			// splitting it into a new occurrence of this slice.
			if sliceLen <= repeatingFieldIndex && (sliceLen == 0 || isDynamic || typeOwnsCode(sliceItemType, codeTag)) {
				sliceItem := reflect.New(sliceItemType).Elem()

				// Initialize element
				initializeStruct(sliceItemType, sliceItem)

				count := setStructFieldByCodeTag(sliceItemType, sliceItem, codeTag, propVal, orderTag, repeatingFieldIndex, isDynamic)
				if count > 0 {
					f.Set(reflect.Append(f, sliceItem))
				}
				setCount = setCount + count

				continue
			}

			// Update existing slice item
			sliceItem := f.Index(sliceLen - 1)
			count := setStructFieldByCodeTag(sliceItemType, sliceItem, codeTag, propVal, orderTag, repeatingFieldIndex, isDynamic)
			setCount = setCount + count

		default:
			tag := ft.Tag.Get(serde.FieldTag)
			if tag != serde.Empty {
				attr, err := serde.GetFieldAttribute(tag)
				if err != nil {
					return 0
				}

				if isDynamic {
					// Set dynamic fields by field name
					if ft.Name == serde.DynamicFieldName(codeTag, orderTag) {
						setFieldValue(f, ft, propVal, attr)
						setCount++
					}
					continue
				}

				if attr.Code == codeTag {
					setFieldValue(f, ft, propVal, attr)
					setCount++
					continue
				}
			}
		}
	}

	return setCount
}

type typeCodeKey struct {
	structType reflect.Type
	code       string
}

var typeOwnsCodeCache sync.Map

// typeOwnsCode reports whether a struct type has a field with the given code
// tag on the type itself (including embedded structs), excluding fields that
// belong to nested slice element types.
func typeOwnsCode(structType reflect.Type, codeTag string) bool {
	if structType.Kind() != reflect.Struct {
		return false
	}

	key := typeCodeKey{structType: structType, code: codeTag}
	if cached, ok := typeOwnsCodeCache.Load(key); ok {
		return cached.(bool)
	}

	owns := false

	for i := 0; i < structType.NumField(); i++ {
		ft := structType.Field(i)

		switch ft.Type.Kind() {
		case reflect.Struct:
			owns = typeOwnsCode(ft.Type, codeTag)

		case reflect.Slice:
			// Nested slice fields are owned by their own element type.

		default:
			tag := ft.Tag.Get(serde.FieldTag)
			if tag != serde.Empty {
				attr, err := serde.GetFieldAttribute(tag)
				if err == nil && attr.Code == codeTag {
					owns = true
				}
			}
		}

		if owns {
			break
		}
	}

	typeOwnsCodeCache.Store(key, owns)

	return owns
}

// Set base field by data type.
func setField(f reflect.Value, val reflect.Value, attr serde.FieldAttribute) {
	if !f.CanSet() {
		return
	}

	switch f.Interface().(type) {
	case int, int8, int16, int32, int64:
		intVal, err := strconv.Atoi(val.String())
		if err == nil {
			f.Set(reflect.ValueOf(intVal))
		}

	case uint, uint8, uint16, uint32, uint64:
		intVal, err := strconv.Atoi(val.String())
		if err == nil {
			f.Set(reflect.ValueOf(uint(intVal)))
		}

	case float32, float64:
		fs := serde.NewFieldSettings(attr)
		flt, err := fs.Unsign(val.String())
		if err == nil && flt != nil {
			f.Set(reflect.ValueOf(*flt))
		}

	case time.Time:
		if attr.Format != serde.Empty {
			dt, err := time.Parse(attr.Format, val.String())
			if err == nil {
				f.Set(reflect.ValueOf(dt))
			}
		}

	default:
		f.Set(val)
	}
}

// Set field value.
func setFieldValue(f reflect.Value, ft reflect.StructField, val reflect.Value, attr serde.FieldAttribute) {
	switch ft.Type.Kind() {
	case reflect.Pointer:
		fv := reflect.New(ft.Type.Elem())

		setField(fv.Elem(), val, attr)
		if f.CanSet() {
			f.Set(fv)
		}

	default:
		setField(f, val, attr)
	}
}

// Initialize new struct.
func initializeStruct(structType reflect.Type, structVal reflect.Value) {
	if structVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < structVal.NumField(); i++ {
		f := structVal.Field(i)
		ft := structType.Field(i)

		switch ft.Type.Kind() {
		case reflect.Map:
			if f.CanSet() {
				f.Set(reflect.MakeMap(ft.Type))
			}

		case reflect.Slice:
			if f.CanSet() {
				f.Set(reflect.MakeSlice(ft.Type, 0, 0))
			}

		case reflect.Chan:
			if f.CanSet() {
				f.Set(reflect.MakeChan(ft.Type, 0))
			}

		case reflect.Struct:
			initializeStruct(ft.Type, f)

		case reflect.Pointer:
			fieldBaseType := reflectionutils.GetElementType(ft.Type)
			switch fieldBaseType.Kind() {
			case reflect.Struct:
				fv := reflect.New(ft.Type.Elem())
				initializeStruct(ft.Type.Elem(), fv.Elem())
				if f.CanSet() {
					f.Set(fv)
				}
			}

		default:
		}
	}
}
