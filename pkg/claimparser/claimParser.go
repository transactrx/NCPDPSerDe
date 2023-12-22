package claimparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
	sliceutils "github.com/transactrx/NCPDPSerDe/pkg/sliceUtils"
	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

const (
	empty          = ""
	fieldTag       = "field"
	headerTag      = "header"
	groupTag       = "group"
	segmentTag     = "segment"
	rawGroupTag    = "rawgroup"
	rawSegmentTag  = "rawsegment"
	codeTag        = "code"
	formatTag      = "format"
	decimalTag     = "decimalPlaces"
	overpunchTag   = "overpunch"
	versionTag     = "version"
	parseMethodTag = "parseMethod"
	rawFieldTag    = "rawField"
)

const (
	dZeroRequestLength  = 56
	dZeroResponseLength = 31
)

type headerAttribute struct {
	version     string
	parseMethod string
	rawField    string
}

type fieldAttribute struct {
	code          string
	format        string
	decimalPlaces int
	overpunch     bool
}

type segmentAttribute struct {
	code string
}

type segmentDefinition struct {
	field reflect.StructField
	tag   segmentAttribute
}

// Formatting map.
var formats = map[string]string{
	"YYYYMMdd": "20060102",
	"HHmmss":   "150405",
	"HHmm":     "1504",
}

var transactionTypeRegistry = map[string]reflect.Type{}
var registerOnce sync.Once

// Add type to registry map.
func registerType(key string, typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	transactionTypeRegistry[key] = t
}

// Get type from registry.
func getRegisteredType(key string) (reflect.Type, error) {
	tranType, ok := transactionTypeRegistry[key]
	if !ok {
		return nil, fmt.Errorf("type not found: %v", key)
	}

	return tranType, nil
}

func getRequestType(tranCode string) (reflect.Type, error) {
	key := fmt.Sprintf("%v|request", tranCode)
	return getRegisteredType(key)
}

func getResponseType(tranCode string) (reflect.Type, error) {
	key := fmt.Sprintf("%v|response", tranCode)
	return getRegisteredType(key)
}

// Register types we need to dynamically create.
func registerTypes() {
	registerOnce.Do(func() {
		// Request types
		registerType("B1|request", (*request.Billing)(nil))
		registerType("B2|request", (*request.Reversal)(nil))
		registerType("B3|request", (*request.Rebill)(nil))
		registerType("E1|request", (*request.Eligibility)(nil))
		registerType("D1|request", (*request.PredeterminationOfBenefits)(nil))
		registerType("S1|request", (*request.ServiceBilling)(nil))
		registerType("S2|request", (*request.ServiceReversal)(nil))
		registerType("S3|request", (*request.ServiceRebill)(nil))
		registerType("P1|request", (*request.PriorAuthorization)(nil))
		registerType("P2|request", (*request.PriorAuthorizationReversal)(nil))
		registerType("P3|request", (*request.PriorAuthorizationInquiry)(nil))
		registerType("P4|request", (*request.PriorAuthorizationRequestOnly)(nil))
		registerType("C1|request", (*request.ControlledSubstanceReporting)(nil))
		registerType("C2|request", (*request.ControlledSubstanceReportingReversal)(nil))
		registerType("C3|request", (*request.ControlledSubstanceReportingRebill)(nil))
		registerType("N1|request", (*request.Information)(nil))
		registerType("N2|request", (*request.InformationReversal)(nil))
		registerType("N3|request", (*request.InformationRebill)(nil))

		// Response types
		registerType("B1|response", (*response.Billing)(nil))
		registerType("B2|response", (*response.Reversal)(nil))
		registerType("B3|response", (*response.Rebill)(nil))
		registerType("E1|response", (*response.Eligibility)(nil))
		registerType("D1|response", (*response.PredeterminationOfBenefits)(nil))
		registerType("S1|response", (*response.ServiceBilling)(nil))
		registerType("S2|response", (*response.ServiceReversal)(nil))
		registerType("S3|response", (*response.ServiceRebill)(nil))
		registerType("P1|response", (*response.PriorAuthorization)(nil))
		registerType("P2|response", (*response.PriorAuthorizationReversal)(nil))
		registerType("P3|response", (*response.PriorAuthorizationInquiry)(nil))
		registerType("P4|response", (*response.PriorAuthorizationRequestOnly)(nil))
		registerType("C1|response", (*response.ControlledSubstanceReporting)(nil))
		registerType("C2|response", (*response.ControlledSubstanceReportingReversal)(nil))
		registerType("C3|response", (*response.ControlledSubstanceReportingRebill)(nil))
		registerType("N1|response", (*response.Information)(nil))
		registerType("N2|response", (*response.InformationReversal)(nil))
		registerType("N3|response", (*response.InformationRebill)(nil))
	})
}

// Parse raw data.
// Good if you don't know anything about the data.
func ParseDynamic(rawClaimString string) (interface{}, error) {
	registerTypes()

	rawClaimString = strings.TrimSpace(rawClaimString)

	// Determine type by first separator index
	firstSeparatorIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})
	if firstSeparatorIndex == -1 {
		return nil, fmt.Errorf("improperly formatted data: no NCPDP separators present")
	}

	// Assume request
	if firstSeparatorIndex == dZeroRequestLength {
		return ParseRequestDynamic(rawClaimString)
	}

	// Assume response
	if firstSeparatorIndex == dZeroResponseLength {
		return ParseResponseDynamic(rawClaimString)
	}

	return nil, fmt.Errorf("unable to determine transaction type")
}

// Parse raw request.
// Good if you know it's a request type.
func ParseRequestDynamic(rawClaimString string) (interface{}, error) {
	registerTypes()

	// Determine transaction type by tran code
	rawClaimString = strings.TrimSpace(rawClaimString)

	if len(rawClaimString) < 10 {
		return nil, fmt.Errorf("unable to determine transaction type")
	}

	tranCode := rawClaimString[8:10]

	claimType, err := getRequestType(tranCode)
	if err != nil {
		return nil, err
	}

	// Dynamically create type
	claimObjectRef := reflect.New(claimType).Elem()
	initializeStruct(claimType, claimObjectRef)

	// Parse raw data into type
	return parseRaw(rawClaimString, claimType, claimObjectRef)
}

// Parse raw request.
// Good if you know it's a response type.
func ParseResponseDynamic(rawClaimString string) (interface{}, error) {
	registerTypes()

	// Determine transaction type by tran code
	rawClaimString = strings.TrimSpace(rawClaimString)

	if len(rawClaimString) < 4 {
		return nil, fmt.Errorf("unable to determine transaction type")
	}

	tranCode := rawClaimString[2:4]

	claimType, err := getResponseType(tranCode)
	if err != nil {
		return nil, err
	}

	// Dynamically create type
	claimObjectRef := reflect.New(claimType).Elem()
	initializeStruct(claimType, claimObjectRef)

	// Parse raw data into type
	return parseRaw(rawClaimString, claimType, claimObjectRef)
}

// Parse raw claim.
// Good if you already know the type.
func ParseRawClaim[V any](rawClaimString string, claim *V) error {
	claimType := reflect.TypeOf(claim)
	claimObjectRef := reflect.ValueOf(claim)

	_, err := parseRaw(rawClaimString, claimType, claimObjectRef)
	if err != nil {
		return err
	}

	return nil
}

// Parse raw data.
func parseRaw(rawClaimString string, claimType reflect.Type, claimObjectRef reflect.Value) (interface{}, error) {
	if strings.TrimSpace(rawClaimString) == empty {
		return nil, fmt.Errorf("NCPDP data is empty")
	}

	rawClaimString = strings.TrimSpace(strings.TrimSuffix(rawClaimString, string(ncpdp.ETX)))

	firstSeparatorIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})
	if firstSeparatorIndex == -1 {
		return nil, fmt.Errorf("improperly formatted data: no NCPDP separators present")
	}

	// Set header
	evaluateHeader(rawClaimString, claimType, claimObjectRef)

	// Set group data
	evaluateGrouping(rawClaimString, claimType, claimObjectRef)

	// Set shared segment data
	endIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.GROUP})
	if endIndex <= 0 {
		endIndex = len(rawClaimString)
	}

	evaluateSegments(rawClaimString[0:endIndex], claimType, claimObjectRef)

	return claimObjectRef.Interface(), nil
}

// Find and parse header.
func evaluateHeader(rawClaimString string, structType reflect.Type, structVal reflect.Value) error {
	headerField, attr, err := getHeaderField(structType)

	if err != nil {
		return err
	}

	if headerField == nil {
		return nil
	}

	if attr.parseMethod == empty {
		return fmt.Errorf("parse method undefined for header")
	}

	if attr.rawField == empty {
		return fmt.Errorf("raw field undefined for header")
	}

	// Get header field object reference
	baseVal := getBaseValue(structVal)
	headerValue := baseVal.FieldByName(headerField.Name)

	rawField := headerValue.FieldByName(attr.rawField)
	if !rawField.IsValid() {
		return fmt.Errorf("raw field not found in header definition")
	}

	if rawField.CanSet() {
		rawField.Set(reflect.ValueOf(rawClaimString))
	}

	parseMethod, err := getMethod(headerValue, attr.parseMethod)
	if err != nil {
		return err
	}

	parseMethod.Call([]reflect.Value{})

	return nil
}

// Get method by name.
func getMethod(structVal reflect.Value, methodName string) (reflect.Value, error) {
	method := structVal.MethodByName(methodName)
	if method.IsValid() {
		return method, nil
	}

	if structVal.CanAddr() {
		method = structVal.Addr().MethodByName(methodName)
		if method.IsValid() {
			return method, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("%v method not found in header definition", methodName)
}

// Evaluate claim groups.
func evaluateGrouping(rawClaimString string, structType reflect.Type, structVal reflect.Value) error {
	rawGroups := stringutils.SplitBySeparator(rawClaimString, ncpdp.GROUP)

	// No groups present
	if len(rawGroups) == 0 {
		return nil
	}

	// Find the slice field with the group tag
	groupField, err := getGroupSlice(structType)
	if err != nil {
		return err
	}

	if groupField == nil {
		return nil
	}

	baseVal := getBaseValue(structVal)
	groupSlice := baseVal.FieldByName(groupField.Name)

	for _, rawClaimGroup := range rawGroups {
		// Create new element for group slice
		groupElementType := groupField.Type.Elem()
		groupItem := reflect.New(groupElementType).Elem()

		// Initialize element
		initializeStruct(groupElementType, groupItem)

		// Set field defined as "raw" data tag
		setStructFieldByCodeTag(groupElementType, groupItem, rawGroupTag, reflect.ValueOf(fmt.Sprint(string(ncpdp.GROUP), rawClaimGroup)), 0)

		// Set segment data
		evaluateSegments(rawClaimGroup, groupElementType, groupItem)

		// Append element to group slice
		groupSlice.Set(reflect.Append(groupSlice, groupItem))
	}

	return nil
}

// Evaluate claim segments.
func evaluateSegments(rawData string, structType reflect.Type, structVal reflect.Value) error {
	rawSegments := stringutils.SplitBySeparator(rawData, ncpdp.SEGMENT)

	baseType := getBaseType(structType)
	baseVal := getBaseValue(structVal)

	for _, rawSeg := range rawSegments {
		segmentMap, err := getSegmentDefinition(baseType)
		if err != nil {
			return err
		}

		// Get all fields
		rawFields := stringutils.SplitBySeparator(rawSeg, ncpdp.FIELD)

		// Empty segment, skip
		if len(rawFields) == 0 {
			return nil
		}

		// Find segment ID in field list
		segmentIdIndex := sliceutils.IndexOfStartsWith(ncpdp.SEGMENT_FIELD_ID, rawFields)
		if segmentIdIndex < 0 {
			return fmt.Errorf("segment id element does not exist")
		}

		// Create new segment and initialize
		segmentIdFieldValue := rawFields[segmentIdIndex]
		segmentDef, ok := segmentMap[segmentIdFieldValue]
		if !ok {
			// segment not found in object definition
			continue
		}
		elementType := segmentDef.field.Type
		segment := reflect.New(elementType).Elem()

		initializeStruct(elementType, segment)

		// Set field defined as "raw" data tag
		setStructFieldByCodeTag(elementType, segment, rawSegmentTag, reflect.ValueOf(fmt.Sprint(string(ncpdp.SEGMENT), rawSeg)), 0)

		fieldCountMap := map[string]int{}

		// Set field data
		for _, rawField := range rawFields {
			fieldCode := stringutils.Substring(rawField, 0, ncpdp.ID_LENGTH)
			fieldValue := stringutils.Substring(rawField, ncpdp.ID_LENGTH, -1)

			fieldCodeCount := fieldCountMap[fieldCode]

			setStructFieldByCodeTag(elementType, segment, fieldCode, reflect.ValueOf(strings.TrimSpace(fieldValue)), fieldCodeCount)

			// Increment field code count
			fieldCodeCount++
			fieldCountMap[fieldCode] = fieldCodeCount
		}

		segmentRef := baseVal.FieldByName(segmentDef.field.Name)
		segmentRef.Set(segment)
	}

	return nil
}

// Set struct field using code tag value to locate appropriate field.
// Ex: `field:"code=EN"`
func setStructFieldByCodeTag(structType reflect.Type, structVal reflect.Value, tagValue string, propVal reflect.Value, repeatingFieldIndex int) int {
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
			count := setStructFieldByCodeTag(ft.Type, f, tagValue, propVal, repeatingFieldIndex)
			setCount = setCount + count

		case reflect.Slice:
			sliceLen := f.Len()
			sliceItemType := f.Type().Elem()

			// Create new slice element
			if sliceLen <= repeatingFieldIndex {
				sliceItem := reflect.New(sliceItemType).Elem()

				// Initialize element
				initializeStruct(sliceItemType, sliceItem)

				count := setStructFieldByCodeTag(sliceItemType, sliceItem, tagValue, propVal, repeatingFieldIndex)
				if count > 0 {
					f.Set(reflect.Append(f, sliceItem))
				}
				setCount = setCount + count

				continue
			}

			// Update existing slice item
			sliceItem := f.Index(sliceLen - 1)
			count := setStructFieldByCodeTag(sliceItemType, sliceItem, tagValue, propVal, repeatingFieldIndex)
			setCount = setCount + count

		default:
			tag := ft.Tag.Get(fieldTag)
			if tag != empty {
				attr, err := getFieldAttribute(tag)
				if err != nil {
					return 0
				}
				if attr.code == tagValue {
					setFieldValue(f, ft, propVal, attr)
					setCount++
				}
			}
		}
	}

	return setCount
}

// Set base field by data type.
func setField(f reflect.Value, val reflect.Value, attr fieldAttribute) {
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
		fs := newFieldSettings(attr)
		flt, err := fs.Unsign(val.String())
		if err == nil && flt != nil {
			f.Set(reflect.ValueOf(*flt))
		}

	case time.Time:
		if attr.format != empty {
			dt, err := time.Parse(attr.format, val.String())
			if err == nil {
				f.Set(reflect.ValueOf(dt))
			}
		}

	default:
		f.Set(val)
	}
}

// Set field value.
func setFieldValue(f reflect.Value, ft reflect.StructField, val reflect.Value, attr fieldAttribute) {
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

// Create field settings from field attribute data.
func newFieldSettings(attr fieldAttribute) ncpdp.FieldSettings {
	return ncpdp.FieldSettings{
		DecimalPlaces: attr.decimalPlaces,
		Format:        attr.format,
		Overpunch:     attr.overpunch,
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
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			if f.CanSet() {
				f.Set(fv)
			}
		default:
		}
	}
}

func getBaseValue(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return getBaseValue(val.Elem())
	default:
		return val
	}
}

func getBaseType(tType reflect.Type) reflect.Type {
	switch tType.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return getBaseType(tType.Elem())
	default:
		return tType
	}
}

// Get header field.
func getHeaderField(tType reflect.Type) (*reflect.StructField, *headerAttribute, error) {
	baseType := getBaseType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		//Return the first header item. Multiples will be ignored.
		tag := field.Tag.Get(headerTag)
		if tag != empty {
			attr, err := getHeaderAttribute(tag)
			if err != nil {
				return nil, nil, err
			}
			return &field, &attr, nil
		}
	}

	return nil, nil, fmt.Errorf("header tag not found in object definition")
}

// Get slice field with the "group" tag.
// Ex: `group:"max=4"`
func getGroupSlice(tType reflect.Type) (*reflect.StructField, error) {
	baseType := getBaseType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		//Return the first group slice found. Multiples will be ignored.
		tag := field.Tag.Get(groupTag)
		if tag != empty {
			if field.Type.Kind() == reflect.Slice {
				return &field, nil
			}
		}
	}

	return nil, fmt.Errorf("group tag not found in object definition")
}

// Get map of segment definitions by segment ID.
func getSegmentDefinition(tType reflect.Type) (map[string]segmentDefinition, error) {
	segmentMap := make(map[string]segmentDefinition)
	baseType := getBaseType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		tag := field.Tag.Get(segmentTag)
		if tag != empty {
			segmentAttribute, err := getSegmentAttribute(tag)
			if err != nil {
				return nil, err
			}

			segmentMap[segmentAttribute.code] = segmentDefinition{tag: segmentAttribute, field: field}
		}
	}

	return segmentMap, nil
}

// Get segment attribute data.
func getSegmentAttribute(tag string) (segmentAttribute, error) {
	attribute := segmentAttribute{}

	attribute.code = parseTag(tag, "code")

	return attribute, nil
}

// Get header attribute data.
func getHeaderAttribute(tag string) (headerAttribute, error) {
	attribute := headerAttribute{}

	tagFields := strings.Split(tag, ",")

	for _, tagField := range tagFields {
		if strings.HasPrefix(tagField, versionTag) {
			attribute.version = parseTag(tagField, versionTag)
		}
		if strings.HasPrefix(tagField, parseMethodTag) {
			attribute.parseMethod = parseTag(tagField, parseMethodTag)
		}
		if strings.HasPrefix(tagField, rawFieldTag) {
			attribute.rawField = parseTag(tagField, rawFieldTag)
		}
	}

	return attribute, nil
}

// Get field attribute data.
func getFieldAttribute(tag string) (fieldAttribute, error) {
	attribute := fieldAttribute{}

	tagFields := strings.Split(tag, ",")

	for _, tagField := range tagFields {
		if strings.HasPrefix(tagField, codeTag) {
			attribute.code = parseTag(tagField, codeTag)
		}
		if strings.HasPrefix(tagField, formatTag) {
			attribute.format = inferFormat(parseTag(tagField, formatTag))
		}
		if strings.HasPrefix(tagField, decimalTag) {
			attribute.decimalPlaces, _ = strconv.Atoi(parseTag(tagField, decimalTag))
		}
		if strings.HasPrefix(tagField, overpunchTag) {
			attribute.overpunch, _ = strconv.ParseBool(parseTag(tagField, overpunchTag))
		}
	}

	return attribute, nil
}

// Get format key and return appropriate value
func inferFormat(key string) string {
	return formats[key]
}

// Parse tag data for specified key/value pair.
func parseTag(tagField, tagKey string) string {
	if !strings.HasPrefix(tagField, tagKey) {
		return empty
	}

	tagFormat := fmt.Sprintf("%v=", tagKey) + "%v"
	var tagValue string

	fmt.Sscanf(tagField, tagFormat, &tagValue)

	return tagValue
}
