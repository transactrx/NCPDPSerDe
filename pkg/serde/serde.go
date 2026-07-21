package serde

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

type HeaderAttribute struct {
	Version      string
	Deserializer string
	Serializer   string
	RawField     string
}

type FieldAttribute struct {
	Code          string
	Format        string
	DecimalPlaces int
	Overpunch     bool
	Order         int
	Dynamic       bool

	// CountFor marks a counter field that is derived automatically during
	// serialization: either the name of a sibling slice field whose length
	// this field reports, or CountForIndex for a 1-based position counter
	// inside a repeating item.
	CountFor string
}

type SegmentAttribute struct {
	Code    string
	Order   int
	Dynamic bool
}

type SegmentDefinition struct {
	Field reflect.StructField
	Tag   SegmentAttribute
}

type FieldDefinition struct {
	Field reflect.StructField
	Tag   FieldAttribute
}

// Formatting map.
//
// For readability only.
// Ex: YYMMdd is more readable than numbers without meaning.
var Formats = map[string]string{
	"YYYYMMdd": "20060102",
	"HHmmss":   "150405",
	"HHmm":     "1504",
}

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

const (
	Empty           = ""
	FieldTag        = "field"
	headerTag       = "header"
	groupTag        = "group"
	SegmentTag      = "segment"
	rawGroupTag     = "rawgroup"
	rawSegmentTag   = "rawsegment"
	CodeTag         = "code"
	formatTag       = "format"
	decimalTag      = "decimalPlaces"
	overpunchTag    = "overpunch"
	versionTag      = "version"
	deserializerTag = "deserializer"
	serializerTag   = "serializer"
	rawFieldTag     = "rawField"
	OrderTag        = "order"
	dynamicTag      = "dynamic"
	countForTag     = "countfor"
)

// CountForIndex is the reserved countfor tag value marking a per-item
// counter (1-based position within the repeating slice).
const CountForIndex = "index"

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

func GetRequestType(tranCode string) (reflect.Type, error) {
	return getRegisteredType(tranCode + "|request")
}

// TransactionTypes returns a copy of the registered transaction type map,
// keyed by "<transactionCode>|request" or "<transactionCode>|response".
func TransactionTypes() map[string]reflect.Type {
	RegisterTypes()

	types := make(map[string]reflect.Type, len(transactionTypeRegistry))
	for k, v := range transactionTypeRegistry {
		types[k] = v
	}

	return types
}

func GetResponseType(tranCode string) (reflect.Type, error) {
	return getRegisteredType(tranCode + "|response")
}

// Register types we need to dynamically create.
func RegisterTypes() {
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

var segmentAttributeCache sync.Map

// Get segment attribute data.
func GetSegmentAttribute(tag string) (SegmentAttribute, error) {
	if cached, ok := segmentAttributeCache.Load(tag); ok {
		return cached.(SegmentAttribute), nil
	}

	attribute := SegmentAttribute{}

	tagFields := strings.Split(tag, ",")

	for _, tagField := range tagFields {
		if strings.HasPrefix(tagField, CodeTag) {
			segmentCode := parseTag(tagField, CodeTag)
			attribute.Code = segmentCode

			if segmentCode == dynamicTag {
				attribute.Dynamic = true
			}
		}
		if strings.HasPrefix(tagField, OrderTag) {
			attribute.Order, _ = strconv.Atoi(parseTag(tagField, OrderTag))
		}
	}

	segmentAttributeCache.Store(tag, attribute)

	return attribute, nil
}

// Get header attribute data.
func GetHeaderAttribute(tag string) (HeaderAttribute, error) {
	attribute := HeaderAttribute{}

	tagFields := strings.Split(tag, ",")

	for _, tagField := range tagFields {
		if strings.HasPrefix(tagField, versionTag) {
			attribute.Version = parseTag(tagField, versionTag)
		}
		if strings.HasPrefix(tagField, deserializerTag) {
			attribute.Deserializer = parseTag(tagField, deserializerTag)
		}
		if strings.HasPrefix(tagField, serializerTag) {
			attribute.Serializer = parseTag(tagField, serializerTag)
		}
		if strings.HasPrefix(tagField, rawFieldTag) {
			attribute.RawField = parseTag(tagField, rawFieldTag)
		}
	}

	return attribute, nil
}

var fieldAttributeCache sync.Map

// Get field attribute data.
func GetFieldAttribute(tag string) (FieldAttribute, error) {
	if cached, ok := fieldAttributeCache.Load(tag); ok {
		return cached.(FieldAttribute), nil
	}

	attribute := FieldAttribute{}

	tagFields := strings.Split(tag, ",")

	for _, tagField := range tagFields {
		if strings.HasPrefix(tagField, CodeTag) {
			code := parseTag(tagField, CodeTag)
			attribute.Code = code

			if code == dynamicTag {
				attribute.Dynamic = true
			}
		}
		if strings.HasPrefix(tagField, formatTag) {
			attribute.Format = inferFormat(parseTag(tagField, formatTag))
		}
		if strings.HasPrefix(tagField, decimalTag) {
			attribute.DecimalPlaces, _ = strconv.Atoi(parseTag(tagField, decimalTag))
		}
		if strings.HasPrefix(tagField, OrderTag) {
			attribute.Order, _ = strconv.Atoi(parseTag(tagField, OrderTag))
		}
		if strings.HasPrefix(tagField, overpunchTag) {
			attribute.Overpunch, _ = strconv.ParseBool(parseTag(tagField, overpunchTag))
		}
		if strings.HasPrefix(tagField, countForTag) {
			attribute.CountFor = parseTag(tagField, countForTag)
		}
	}

	fieldAttributeCache.Store(tag, attribute)

	return attribute, nil
}

// Get format key and return appropriate value
func inferFormat(key string) string {
	return Formats[key]
}

// Parse tag data for specified key/value pair.
func parseTag(tagField, tagKey string) string {
	tagValue, found := strings.CutPrefix(tagField, tagKey+"=")
	if !found {
		return Empty
	}

	return tagValue
}

var segmentDefinitionByIdCache sync.Map

// Get map of segment definitions by segment ID. The returned map is cached and must not be modified.
func GetSegmentDefinitionById(tType reflect.Type) (map[string]SegmentDefinition, error) {
	baseType := reflectionutils.GetElementType(tType)

	if cached, ok := segmentDefinitionByIdCache.Load(baseType); ok {
		return cached.(map[string]SegmentDefinition), nil
	}

	segmentMap := make(map[string]SegmentDefinition)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		tag := field.Tag.Get(SegmentTag)
		if tag != Empty {
			segmentAttribute, err := GetSegmentAttribute(tag)
			if err != nil {
				return nil, err
			}

			segmentMap[segmentAttribute.Code] = SegmentDefinition{Tag: segmentAttribute, Field: field}
		}
	}

	segmentDefinitionByIdCache.Store(baseType, segmentMap)

	return segmentMap, nil
}

var segmentDefinitionByOrderCache sync.Map

// Get map of segment definitions by segment order. The returned map is cached and must not be modified.
func GetSegmentDefinitionByOrder(tType reflect.Type) (map[int]SegmentDefinition, error) {
	baseType := reflectionutils.GetElementType(tType)

	if cached, ok := segmentDefinitionByOrderCache.Load(baseType); ok {
		return cached.(map[int]SegmentDefinition), nil
	}

	segmentMap := make(map[int]SegmentDefinition)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		tag := field.Tag.Get(SegmentTag)
		if tag != Empty {
			segmentAttribute, err := GetSegmentAttribute(tag)
			if err != nil {
				return nil, err
			}

			if segmentAttribute.Order != 0 {
				segmentMap[segmentAttribute.Order] = SegmentDefinition{Tag: segmentAttribute, Field: field}
			}
		}
	}

	segmentDefinitionByOrderCache.Store(baseType, segmentMap)

	return segmentMap, nil
}

// Get header field.
func GetHeaderField(tType reflect.Type) (*reflect.StructField, *HeaderAttribute, error) {
	baseType := reflectionutils.GetElementType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		//Return the first header item. Multiples will be ignored.
		tag := field.Tag.Get(headerTag)
		if tag != Empty {
			attr, err := GetHeaderAttribute(tag)
			if err != nil {
				return nil, nil, err
			}
			return &field, &attr, nil
		}
	}

	return nil, nil, fmt.Errorf("header tag not found in object definition")
}

// Get method by name.
func GetMethod(structVal reflect.Value, methodName string) (reflect.Value, error) {
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

// Create field settings from field attribute data.
func NewFieldSettings(attr FieldAttribute) ncpdp.FieldSettings {
	return ncpdp.FieldSettings{
		DecimalPlaces: attr.DecimalPlaces,
		Format:        attr.Format,
		Overpunch:     attr.Overpunch,
	}
}

// Get slice field with the "group" tag.
// Ex: `group:"max=4"`
func GetGroupSlice(tType reflect.Type) (*reflect.StructField, error) {
	baseType := reflectionutils.GetElementType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		//Return the first group slice found. Multiples will be ignored.
		tag := field.Tag.Get(groupTag)
		if tag != Empty {
			if field.Type.Kind() == reflect.Slice {
				return &field, nil
			}
		}
	}

	return nil, fmt.Errorf("group tag not found in object definition")
}

// Get slice field where the tag refects code=dynamic.
// Supported tagTypes: segment, field.
//
// Ex: `segment:"code=dynamic"`
// Ex: `field:"code=dynamic"`
func GetDynamicSlice(tType reflect.Type, tagType string) (*reflect.StructField, error) {
	baseType := reflectionutils.GetElementType(tType)

	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)

		tag := field.Tag.Get(tagType)
		if tag != Empty {
			switch tagType {
			case SegmentTag:
				segmentAttribute, err := GetSegmentAttribute(tag)
				if err != nil {
					return nil, err
				}
				if segmentAttribute.Dynamic && field.Type.Kind() == reflect.Slice {
					return &field, nil
				}

			case FieldTag:
				fieldAttribute, err := GetFieldAttribute(tag)
				if err != nil {
					return nil, err
				}
				if fieldAttribute.Dynamic && field.Type.Kind() == reflect.Slice {
					return &field, nil
				}
			}

		}
	}

	return nil, fmt.Errorf("dynamic code tag not found in object definition: code=dynamic")
}

// Create dynamic field name. Handle repeating fields by
// including the element index within the name.
func DynamicFieldName(fieldCode string, order int) string {
	return "Field_" + replaceIllegalCharacters(fieldCode) + "_" + strconv.Itoa(order)
}

// Replace illegal characters in field names.
func replaceIllegalCharacters(str string) string {
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

// Create dynamic segment type.
//
// Note: Structs are immutable and cannot be edited after creation.
// All known elements must be defined up front.
func NewDynamicSegmentType(rawData string) reflect.Type {
	dyFieldList := dynamic.DynamicFieldList{}

	rawFields := stringutils.SplitBySeparator(rawData, ncpdp.FIELD)

	// Iterate fields, defining type for each
	for i := 0; i < len(rawFields); i++ {
		rawField := rawFields[i]
		fieldCode := stringutils.Substring(rawField, 0, ncpdp.ID_LENGTH)

		order := i + 1
		name := DynamicFieldName(fieldCode, order)

		dynamicField := NewDynamicStringPointerField(name, FieldTag, map[string]string{CodeTag: fieldCode, OrderTag: strconv.Itoa(order)})
		dyFieldList = append(dyFieldList, dynamicField)
	}

	// Convert dynamic type list to struct field list
	sf := reflect.StructField{
		Name: "Raw",
		Type: reflectionutils.StringPointerType(),
		Tag:  reflect.StructTag(`field:"code=rawsegment"`),
	}
	sfList := []reflect.StructField{sf}
	sfList = append(sfList, dyFieldList.ToStructFieldList()...)

	segmentType := reflect.StructOf(sfList)

	return segmentType
}

// Create dynamic string pointer field
func NewDynamicStringPointerField(fieldName, tagName string, tags map[string]string) dynamic.DynamicFieldData {
	return dynamic.DynamicFieldData{
		Name:    fieldName,
		Type:    reflectionutils.StringPointerType(),
		Tags:    tags,
		TagName: tagName,
	}
}
