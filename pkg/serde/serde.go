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

	// SinceVersion scopes a field to NCPDP versions at or after the named
	// version (e.g. "F6"). The serializer omits such a field when emitting an
	// older transmission, so version-specific fields — notably the repeating
	// group counters added in F6 (KR, RR) whose backing slices are dual-mapped
	// from D0 scalars — never leak into a D0 message. Empty means the field
	// applies to every version.
	SinceVersion string

	// DecimalPlacesSinceVersion/DecimalPlacesSince override DecimalPlaces for
	// transmissions at or after the named version (tag syntax:
	// decimalPlacesSince=F6:7). F6 widened 448-ED Compound Ingredient Quantity
	// from 9(7)v999 to 9(7)v9999999 — the only field whose implied decimal
	// count changed between D0 and F6 — so the same struct field must parse
	// and serialize with version-appropriate precision. Resolve via
	// DecimalPlacesFor; an empty version name means no override.
	DecimalPlacesSinceVersion string
	DecimalPlacesSince        int

	// UnsignedSinceVersion turns Overpunch off for transmissions at or after
	// the named version (tag syntax: unsignedSince=F6). F6 made the percentage
	// tax rates 483-HE and 560-AY unsigned (s9(3)v9999 became 9(3)v9999), so
	// serializing them for F6 must emit plain digits, never an overpunch
	// character. Parsing is unaffected: overpunch removal already passes plain
	// digits through. Resolve via OverpunchFor; empty means no override.
	UnsignedSinceVersion string
}

// DecimalPlacesFor resolves the implied decimal count for the given
// transmission version, honoring the decimalPlacesSince override. Version
// ordering comes from ncpdp.VersionAtLeast, so an override scoped to F6
// applies to F6 and every later version, while an unknown or blank version
// ranks oldest and keeps the base DecimalPlaces.
func (attr FieldAttribute) DecimalPlacesFor(version string) int {
	if attr.DecimalPlacesSinceVersion != Empty && ncpdp.VersionAtLeast(version, attr.DecimalPlacesSinceVersion) {
		return attr.DecimalPlacesSince
	}

	return attr.DecimalPlaces
}

// OverpunchFor resolves whether serialization applies overpunch signing for
// the given transmission version, honoring the unsignedSince override. An
// unknown or blank version ranks oldest and keeps the base Overpunch.
func (attr FieldAttribute) OverpunchFor(version string) bool {
	if attr.UnsignedSinceVersion != Empty && ncpdp.VersionAtLeast(version, attr.UnsignedSinceVersion) {
		return false
	}

	return attr.Overpunch
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
	decimalSinceTag = "decimalPlacesSince"
	unsignedTag     = "unsignedSince"
	overpunchTag    = "overpunch"
	versionTag      = "version"
	sinceVersionTag = "sinceVersion"
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
		key, value, found := strings.Cut(tagField, "=")
		if !found {
			continue
		}

		switch key {
		case CodeTag:
			attribute.Code = value
			attribute.Dynamic = value == dynamicTag
		case formatTag:
			attribute.Format = inferFormat(value)
		case decimalTag:
			attribute.DecimalPlaces, _ = strconv.Atoi(value)
		case decimalSinceTag:
			if version, places, ok := strings.Cut(value, ":"); ok {
				attribute.DecimalPlacesSinceVersion = version
				attribute.DecimalPlacesSince, _ = strconv.Atoi(places)
			}
		case unsignedTag:
			attribute.UnsignedSinceVersion = value
		case OrderTag:
			attribute.Order, _ = strconv.Atoi(value)
		case overpunchTag:
			attribute.Overpunch, _ = strconv.ParseBool(value)
		case countForTag:
			attribute.CountFor = value
		case sinceVersionTag:
			attribute.SinceVersion = value
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

// Create field settings from field attribute data, using the base (oldest
// version) decimal places. Prefer NewVersionedFieldSettings when the
// transmission version is known.
func NewFieldSettings(attr FieldAttribute) ncpdp.FieldSettings {
	return NewVersionedFieldSettings(attr, Empty)
}

// Create field settings resolved for the given transmission version:
// decimalPlacesSince and unsignedSince overrides take effect at or after
// their named version.
func NewVersionedFieldSettings(attr FieldAttribute, version string) ncpdp.FieldSettings {
	return ncpdp.FieldSettings{
		DecimalPlaces: attr.DecimalPlacesFor(version),
		Format:        attr.Format,
		Overpunch:     attr.OverpunchFor(version),
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
	return dynamic.FieldName(fieldCode, order)
}

// Create dynamic segment type.
//
// Note: Structs are immutable and cannot be edited after creation.
// All known elements must be defined up front.
func NewDynamicSegmentType(rawData string) reflect.Type {
	dyFieldList := dynamic.DynamicFieldList{}

	// Keep a prefix field (a segment id without its leading separator) so the
	// dynamic type's field indices line up with the deserializer's field walk.
	rawFields := stringutils.SplitBySeparatorKeepingPrefix(rawData, ncpdp.FIELD)

	// Iterate fields, defining type for each. The code split must stay
	// rune-aware and aligned with the deserializer's field walk.
	for i := 0; i < len(rawFields); i++ {
		rawField := rawFields[i]
		fieldCode, _ := stringutils.SplitRunePrefix(rawField, ncpdp.ID_LENGTH)

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
