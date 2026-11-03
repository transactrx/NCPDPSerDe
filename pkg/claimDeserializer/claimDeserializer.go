package claimdeserializer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

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

var timeType = reflect.TypeOf(time.Time{})

const (
	dZeroRequestLength  = 56
	dZeroResponseLength = 31
	f6RequestLength     = 58
)

// trimClaimFrame removes transport framing from a raw claim: leading
// whitespace (which would throw off header-offset detection) and trailing
// CR/LF line endings. Trailing blanks are left intact — they belong to the
// last field's value (e.g. a blank-padded H6) and must round-trip unchanged.
func trimClaimFrame(raw string) string {
	raw = strings.TrimLeftFunc(raw, unicode.IsSpace)
	return strings.TrimRight(raw, "\r\n")
}

// Parse raw data.
// Good if you don't know anything about the data.
func Deserialize(rawClaimString string) (interface{}, error) {
	serde.RegisterTypes()

	rawClaimString = trimClaimFrame(rawClaimString)

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

	return nil, fmt.Errorf("unable to determine transaction type%s", headerShapeHint(rawClaimString, true))
}

// headerShapeHint describes where the first separator sits relative to the
// known fixed header lengths. A transaction code read from a header of the
// wrong shape is meaningless bytes, so a failed type lookup names the likely
// cause (response fed to the request parser or vice versa, leading garbage, a
// truncated header) instead of only echoing those bytes.
func headerShapeHint(rawClaimString string, parsingRequest bool) string {
	idx := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})

	switch {
	case idx == dZeroRequestLength || idx == f6RequestLength:
		if !parsingRequest {
			return fmt.Sprintf(" (the first separator is at byte %d, the shape of a request header — this data may be a request passed to the response parser)", idx)
		}

	case idx == dZeroResponseLength:
		if parsingRequest {
			return fmt.Sprintf(" (the first separator is at byte %d, the shape of a response header — this data may be a response passed to the request parser)", idx)
		}

	case idx < 0:
		return " (no NCPDP separators present)"

	default:
		return fmt.Sprintf(" (the first separator is at byte %d, which matches no known header length — D0 request 56, F6 request 58, response 31 — so the data may have leading garbage or a truncated header)", idx)
	}

	return ""
}

// Parse raw request.
// Good if you know it's a request type.
func DeserializeRequest(rawClaimString string) (interface{}, error) {
	serde.RegisterTypes()

	// Determine transaction type by tran code. Only framing is removed (ETX,
	// leading whitespace, trailing CR/LF): trailing spaces are significant
	// NCPDP data and must round-trip unchanged.
	rawClaimString = strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), "")
	rawClaimString = trimClaimFrame(rawClaimString)

	if len(rawClaimString) < 10 {
		return nil, fmt.Errorf("unable to determine transaction type: data is %d bytes, shorter than any request header", len(rawClaimString))
	}

	// D0 request headers lead with the BIN (digits); vEB+ headers (e.g. F6)
	// lead with the version, which moves the transaction code to bytes 2-4.
	tranCode := rawClaimString[8:10]
	if ncpdp.HeaderLeadsWithVersion(rawClaimString[:2]) {
		tranCode = rawClaimString[2:4]
	}

	claimType, err := serde.GetRequestType(tranCode)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, headerShapeHint(rawClaimString, true))
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

	// Determine transaction type by tran code. Only framing is removed (ETX,
	// leading whitespace, trailing CR/LF): trailing spaces are significant
	// NCPDP data and must round-trip unchanged.
	rawClaimString = strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), "")
	rawClaimString = trimClaimFrame(rawClaimString)

	if len(rawClaimString) < 4 {
		return nil, fmt.Errorf("unable to determine transaction type: data is %d bytes, shorter than any response header", len(rawClaimString))
	}

	tranCode := rawClaimString[2:4]

	claimType, err := serde.GetResponseType(tranCode)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, headerShapeHint(rawClaimString, false))
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

	// Only framing is removed (ETX, leading whitespace, trailing CR/LF):
	// trailing spaces belong to the last field's value and must round-trip
	// unchanged.
	rawClaimString = strings.ReplaceAll(rawClaimString, string(ncpdp.ETX), "")
	rawClaimString = trimClaimFrame(rawClaimString)

	firstSeparatorIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})
	if firstSeparatorIndex == -1 {
		return nil, fmt.Errorf("improperly formatted data: no NCPDP separators present")
	}

	// Set header
	if err := evaluateHeader(rawClaimString, claimType, claimObjectRef); err != nil {
		return nil, err
	}

	groupIndex := stringutils.IndexOfAny(rawClaimString, 0, []byte{ncpdp.GROUP})

	// The transmission version gates group-separator omission and, during
	// field consumption, version-scoped (sinceVersion) typed fields: a value
	// received for a field scoped to a NEWER version than the transmission
	// must stay unconsumed and preserve as a dynamic field, because the
	// serializer omits such fields for this version and the received bytes
	// would otherwise vanish silently.
	version := ncpdp.DetectTransmissionVersion(rawClaimString)

	// Segments materialized across every part of the transmission. A
	// transmission that parses to zero segments (a header followed only by
	// bare separators or id-only fragments) carries no data and cannot be
	// re-serialized into anything this parser would accept, so it refuses
	// loudly instead of round-tripping to an unparseable husk.
	segmentCount := 0

	if groupIndex < 0 && ncpdp.OmitsGroupSeparator(version) {
		// vEB+ (e.g. F6) eliminated the group separator and allows only a single
		// transaction per transmission; the claim segments directly follow the
		// shared segments and are located by segment ID instead.
		count, err := evaluateUngroupedTransaction(rawClaimString, claimType, claimObjectRef, version)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("improperly formatted data: the transmission contains a header but no parseable segments")
		}

		return claimObjectRef.Interface(), nil
	}

	// Set group data
	groupCount, err := evaluateGrouping(rawClaimString, claimType, claimObjectRef, version)
	if err != nil {
		return nil, err
	}
	segmentCount += groupCount

	// Set shared segment data
	endIndex := groupIndex
	if endIndex <= 0 {
		endIndex = len(rawClaimString)
	}

	sharedCount, err := evaluateSegments(rawClaimString[0:endIndex], claimType, claimObjectRef, version)
	if err != nil {
		return nil, err
	}
	segmentCount += sharedCount

	if segmentCount == 0 {
		return nil, fmt.Errorf("improperly formatted data: the transmission contains a header but no parseable segments")
	}

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

	// The parse method reports malformed headers (short data, misaligned
	// fixed-width fields); swallowing that error would let a garbled header
	// masquerade as a parsed one.
	ret := deserMethod.Call([]reflect.Value{})
	if len(ret) > 0 {
		if parseErr, ok := ret[0].Interface().(error); ok && parseErr != nil {
			return parseErr
		}
	}

	return nil
}

// Evaluate claim groups. Returns the number of segments materialized across
// every group.
func evaluateGrouping(rawClaimString string, structType reflect.Type, structVal reflect.Value, version string) (int, error) {
	rawGroups := stringutils.SplitBySeparator(rawClaimString, ncpdp.GROUP)

	// No groups present
	if len(rawGroups) == 0 {
		return 0, nil
	}

	// Find the slice field with the group tag. A transaction type with no
	// claim section drops any claim groups the transmission carries: per the
	// NCPDP standard an eligibility request has no claim section in D.0 or
	// F6, but non-standard DX data-collection E1 feeds append AM11/AMXX
	// anyway, and the platform intentionally strips those during claim
	// compression before forwarding. Parsing the standard portion instead of
	// failing matches that flow.
	groupField, _ := serde.GetGroupSlice(structType)
	if groupField == nil {
		return 0, nil
	}

	baseVal := reflectionutils.GetElementValue(structVal)
	groupSlice := baseVal.FieldByName(groupField.Name)

	segmentCount := 0

	for _, rawClaimGroup := range rawGroups {
		// Create new element for group slice
		groupElementType := groupField.Type.Elem()
		groupItem := reflect.New(groupElementType).Elem()

		// Initialize element
		initializeStruct(groupElementType, groupItem)

		// Set field defined as "raw" data tag
		setStructFieldByCodeTag(groupElementType, groupItem, fieldTarget{code: rawGroupTag, value: reflect.ValueOf(string(ncpdp.GROUP) + rawClaimGroup), version: version})

		// Set segment data
		count, err := evaluateSegments(rawClaimGroup, groupElementType, groupItem, version)
		if err != nil {
			return segmentCount, err
		}
		segmentCount += count

		// Append element to group slice
		groupSlice.Set(reflect.Append(groupSlice, groupItem))
	}

	return segmentCount, nil
}

// Evaluate an ungrouped (vEB+/F6) transmission. No group separators exist, so
// the single transaction's segments are located by segment ID: the first
// segment that belongs to the claim group struct starts the transaction, and
// it plus everything after it parse as one claim group. Segments before it
// parse as shared segments.
func evaluateUngroupedTransaction(rawClaimString string, structType reflect.Type, structVal reflect.Value, version string) (int, error) {
	// Transaction types without a claim section (e.g. E1 eligibility requests)
	// have no group-tagged slice; the whole transmission parses as shared
	// segments.
	groupField, _ := serde.GetGroupSlice(structType)

	transactionIndex, err := findTransactionIndex(rawClaimString, groupField)
	if err != nil {
		return 0, err
	}

	sharedRaw := rawClaimString
	if transactionIndex != -1 {
		sharedRaw = rawClaimString[:transactionIndex]
	}

	segmentCount, err := evaluateSegments(sharedRaw, structType, structVal, version)
	if err != nil {
		return segmentCount, err
	}

	if transactionIndex == -1 {
		return segmentCount, nil
	}

	transactionRaw := rawClaimString[transactionIndex:]

	groupElementType := groupField.Type.Elem()
	groupItem := reflect.New(groupElementType).Elem()

	initializeStruct(groupElementType, groupItem)

	setStructFieldByCodeTag(groupElementType, groupItem, fieldTarget{code: rawGroupTag, value: reflect.ValueOf(transactionRaw), version: version})

	count, err := evaluateSegments(transactionRaw, groupElementType, groupItem, version)
	if err != nil {
		return segmentCount, err
	}
	segmentCount += count

	baseVal := reflectionutils.GetElementValue(structVal)
	groupSlice := baseVal.FieldByName(groupField.Name)
	groupSlice.Set(reflect.Append(groupSlice, groupItem))

	return segmentCount, nil
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
	rawFields := stringutils.SplitBySeparatorKeepingPrefix(rawSeg, ncpdp.FIELD)

	idIndex := sliceutils.IndexOfStartsWith(ncpdp.SEGMENT_FIELD_ID, rawFields)
	if idIndex < 0 {
		return serde.Empty
	}

	return rawFields[idIndex]
}

// dynamicFieldOrderBase keeps dynamically added field orders above every
// static field order. The serializer keys its per-segment field map by order,
// so a dynamic field whose order collided with a static field's order (or
// with the first order of a nested repeating group) would silently replace
// that field in the rebuilt segment. Static orders are small (< 100).
const dynamicFieldOrderBase = 1 << 16

// Add elements to dynamic field list.
func appendDynamicFields(structType reflect.Type, structVal reflect.Value, dyFieldCodes, dyFieldValues []string) error {
	baseType := reflectionutils.GetElementType(structType)
	startOrder := dynamicFieldOrderBase + baseType.NumField()

	// Offset by the dynamic fields already held so successive calls (one per
	// preserved field of a repeating-group occurrence) keep unique orders —
	// the serializer keys its field map by order, so a collision would
	// silently drop a preserved field.
	if existing, found := findDynamicSliceValue(structType, structVal, serde.FieldTag); found && existing.Kind() == reflect.Slice {
		startOrder += existing.Len()
	}

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
		setStructFieldByCodeTag(typ, field, fieldTarget{code: dyCode, value: reflect.ValueOf(dyValue), order: startOrder, isDynamic: true})

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

// findDynamicSliceValue locates the dynamic slice (the field tagged
// code=dynamic for the given tag type), looking through nested structs: some
// segments keep their DynamicFields inside an embedded group (e.g.
// InsuranceAdditionalInformation.MedicarePartD), and a top-level-only lookup
// silently lost every field preserved on such a segment.
func findDynamicSliceValue(structType reflect.Type, structVal reflect.Value, tagType string) (reflect.Value, bool) {
	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	if baseType.Kind() != reflect.Struct || !baseVal.IsValid() {
		return reflect.Value{}, false
	}

	if dynamicField, err := serde.GetDynamicSlice(baseType, tagType); err == nil && dynamicField != nil {
		if v := baseVal.FieldByName(dynamicField.Name); v.IsValid() {
			return v, true
		}
	}

	for i := 0; i < baseType.NumField(); i++ {
		if baseType.Field(i).Type.Kind() != reflect.Struct {
			continue
		}
		if v, ok := findDynamicSliceValue(baseType.Field(i).Type, baseVal.Field(i), tagType); ok {
			return v, ok
		}
	}

	return reflect.Value{}, false
}

// Add element to dynamic segment list
func appendDynamicSlice(structType reflect.Type, structVal reflect.Value, dynamicElementValue reflect.Value, tagType string) error {
	dynamicObj, found := findDynamicSliceValue(structType, structVal, tagType)
	if !found {
		return fmt.Errorf("dynamic code tag not found in object definition: code=dynamic")
	}

	if !dynamicObj.CanSet() {
		return fmt.Errorf("dynamic object is read-only")
	}

	// Append element to slice
	dynamicObj.Set(reflect.Append(dynamicObj, dynamicElementValue))

	return nil
}

// Evaluate claim segments. Returns the number of segments materialized (empty
// and id-only fragments are skipped and do not count).
func evaluateSegments(rawData string, structType reflect.Type, structVal reflect.Value, version string) (int, error) {
	rawSegments := stringutils.SplitBySeparator(rawData, ncpdp.SEGMENT)

	baseType := reflectionutils.GetElementType(structType)
	baseVal := reflectionutils.GetElementValue(structVal)

	segmentMap, err := serde.GetSegmentDefinitionById(baseType)
	if err != nil {
		return 0, err
	}

	// A non-repeating segment sent twice must not overwrite its first
	// occurrence: the first parse wins and later copies are preserved as
	// dynamic segments, mirroring how duplicated scalar fields are kept.
	parsedSegments := map[string]bool{}

	segmentCount := 0

	for _, rawSeg := range rawSegments {
		materialized, err := evaluateSegment(rawSeg, baseType, baseVal, segmentMap, parsedSegments, version)
		if err != nil {
			return segmentCount, err
		}
		if materialized {
			segmentCount++
		}
	}

	return segmentCount, nil
}

// evaluateSegment parses one raw segment into the transaction model and
// reports whether it materialized. Empty fragments, id-only fragments, and
// segments carrying nothing beyond their identifier are skipped: they render
// empty on rebuild, so materializing them would let a transmission round-trip
// to output the parser itself refuses.
func evaluateSegment(rawSeg string, baseType reflect.Type, baseVal reflect.Value, segmentMap map[string]serde.SegmentDefinition, parsedSegments map[string]bool, version string) (bool, error) {
	// Get all fields, keeping a segment id that arrived without its
	// leading field separator (e.g. "<RS>AM29<FS>...").
	rawFields := stringutils.SplitBySeparatorKeepingPrefix(rawSeg, ncpdp.FIELD)

	// An empty segment ("<RS><RS>") carries no data: skip it and keep
	// parsing — returning here used to silently drop every remaining
	// segment of the part. The stray separator normalizes away on rebuild.
	if len(rawFields) == 0 {
		return false, nil
	}

	// Find segment ID in field list
	segmentIdIndex := sliceutils.IndexOfStartsWith(ncpdp.SEGMENT_FIELD_ID, rawFields)
	if segmentIdIndex < 0 {
		return false, fmt.Errorf("segment id element does not exist: a segment with %d field(s), starting with field code %q, has no AM segment identifier, so its data cannot be attributed", len(rawFields), stringutils.Substring(rawFields[0], 0, ncpdp.ID_LENGTH))
	}

	// A segment that is nothing but its own identifier ("<RS><FS>AM22")
	// carries no data fields and cannot survive serialization (the
	// serializer drops empty segments, and SegmentId is excluded from
	// JSON), so it normalizes away like an empty segment instead of
	// round-tripping inconsistently.
	if len(rawFields) == 1 && len(rawFields[segmentIdIndex]) <= 4 {
		return false, nil
	}

	// Create new segment and initialize. A non-repeating segment already
	// parsed once stays dynamic: the first parse won the typed field.
	var elementType reflect.Type
	segmentIdFieldValue := rawFields[segmentIdIndex]
	segmentDef, ok := segmentMap[segmentIdFieldValue]
	if ok && parsedSegments[segmentIdFieldValue] {
		ok = false
	}
	dynamicSegment := dynamic.DynamicStruct{}
	isDynamicSegment := false

	if ok {
		parsedSegments[segmentIdFieldValue] = true
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
	setStructFieldByCodeTag(elementType, segment, fieldTarget{code: rawSegmentTag, value: reflect.ValueOf(string(ncpdp.SEGMENT) + rawSeg), version: version})

	fields := setSegmentFields(elementType, segment, rawFields, segmentIdIndex, isDynamicSegment, version)
	if !fields.carriesData {
		return false, nil
	}

	if isDynamicSegment {
		segmentIf := segment.Interface()
		dynamicSegment.Value = &segmentIf
		appendDynamicSlice(baseType, baseVal, reflect.ValueOf(dynamicSegment), serde.SegmentTag)
		return true, nil
	}

	// Preservation failure is data loss, never silent: every segment
	// model carries a dynamic slice (possibly in a nested group), so
	// this only fires on a model gap.
	if err := appendDynamicFields(segment.Type(), segment, fields.dyCodes, fields.dyValues); err != nil {
		return false, fmt.Errorf("cannot preserve unconsumed fields on segment %v: %w", segmentIdFieldValue, err)
	}
	baseVal.FieldByName(segmentDef.Field.Name).Set(segment)

	return true, nil
}

// segmentFieldData is what one segment's field walk leaves behind: the codes
// and values no typed field consumed (preserved later as dynamic fields), and
// carriesData — whether anything beyond the segment identifier survives into
// the model. A segment holding only its id and discarded counters renders
// empty on rebuild, so materializing it would let a transmission round-trip
// to output the parser itself refuses.
type segmentFieldData struct {
	dyCodes     []string
	dyValues    []string
	carriesData bool
}

// setSegmentFields routes every raw field of one segment into the segment
// model.
func setSegmentFields(elementType reflect.Type, segment reflect.Value, rawFields []string, segmentIdIndex int, isDynamicSegment bool, version string) segmentFieldData {
	fields := segmentFieldData{dyCodes: []string{}, dyValues: []string{}}
	fieldCountMap := map[string]int{}

	for fieldIndex, rawField := range rawFields {
		fieldCode, fieldValue := stringutils.SplitRunePrefix(rawField, ncpdp.ID_LENGTH)

		setCount, discarded := setStructFieldByCodeTag(elementType, segment, fieldTarget{
			code:                fieldCode,
			value:               reflect.ValueOf(fieldValue),
			order:               fieldIndex + 1,
			repeatingFieldIndex: fieldCountMap[fieldCode],
			isDynamic:           isDynamicSegment,
			version:             version,
		})
		if setCount == 0 && !discarded {
			fields.dyCodes = append(fields.dyCodes, fieldCode)
			fields.dyValues = append(fields.dyValues, fieldValue)
		}
		if fieldIndex != segmentIdIndex && !discarded {
			fields.carriesData = true
		}

		// Increment field code count
		fieldCountMap[fieldCode]++
	}

	return fields
}

// fieldTarget identifies where one raw field value may land: the NCPDP field
// code, the value itself, the field's position within its segment (order),
// which occurrence of a repeating code this is, whether the holding segment
// is dynamic, and the transmission version that gates version-scoped fields.
type fieldTarget struct {
	code                string
	value               reflect.Value
	order               int
	repeatingFieldIndex int
	isDynamic           bool
	version             string
}

// Set struct field using code tag value to locate appropriate field.
// Ex: `field:"code=EN"`
//
// Returns how many fields consumed the value, and whether the value was
// deliberately DISCARDED (an unparseable countfor counter): a discarded value
// must be neither preserved as a dynamic field nor treated as segment data.
func setStructFieldByCodeTag(structType reflect.Type, structVal reflect.Value, target fieldTarget) (int, bool) {
	if structVal.Kind() != reflect.Struct {
		return 0, false
	}

	setCount := 0
	discarded := false

	for i := 0; i < structVal.NumField(); i++ {
		f := structVal.Field(i)
		ft := structType.Field(i)

		var count int
		var disc bool

		switch ft.Type.Kind() {
		case reflect.Struct:
			// Embedded struct
			count, disc = setStructFieldByCodeTag(ft.Type, f, target)

		case reflect.Slice:
			count, disc = setSliceFieldByCodeTag(f, target)

		default:
			var err error
			count, disc, err = setScalarFieldByCodeTag(f, ft, target)
			if err != nil {
				// A malformed field tag on the model aborts the walk, as it
				// always has: it is a programming error, not a data problem.
				return 0, false
			}
		}

		setCount += count
		discarded = discarded || disc
	}

	return setCount, discarded
}

// setSliceFieldByCodeTag routes a field value into a repeating group: either
// into the current (last) occurrence or into a new occurrence appended for it.
func setSliceFieldByCodeTag(f reflect.Value, target fieldTarget) (int, bool) {
	sliceItemType := f.Type().Elem()

	if !sliceStartsNewOccurrence(f, sliceItemType, target) {
		// Update existing slice item
		return setSliceItemField(sliceItemType, f.Index(f.Len()-1), target)
	}

	sliceItem := reflect.New(sliceItemType).Elem()

	// Initialize element
	initializeStruct(sliceItemType, sliceItem)

	count, disc := setSliceItemField(sliceItemType, sliceItem, target)
	if count > 0 {
		f.Set(reflect.Append(f, sliceItem))
	}

	return count, disc
}

// sliceStartsNewOccurrence decides whether this code starts a new occurrence
// of the repeating group or belongs to the current (last) one. NCPDP
// repeating groups are delimited by the reappearance of a code the element
// type owns, never by the code's position in the whole segment — a count
// accumulated across the segment misassigns nested repeats under the second
// occurrence of an outer group (e.g. the HC/DV pair of a second AM05 other
// payer). Only grow this slice when the element type owns the field code
// directly; a repeating code that lives in a deeper nested slice must repeat
// inside the last element instead of splitting it into a new occurrence of
// this slice.
func sliceStartsNewOccurrence(f reflect.Value, sliceItemType reflect.Type, target fieldTarget) bool {
	if target.isDynamic {
		return f.Len() <= target.repeatingFieldIndex
	}
	if f.Len() == 0 {
		return true
	}

	return typeOwnsCode(sliceItemType, target.code) && codeOccupied(sliceItemType, f.Index(f.Len()-1), target.code)
}

// setSliceItemField lands the value inside one occurrence of a repeating
// group, falling back to that occurrence's own dynamic fields when no typed
// field could represent it (see preserveItemDynamicField).
func setSliceItemField(sliceItemType reflect.Type, sliceItem reflect.Value, target fieldTarget) (int, bool) {
	count, disc := setStructFieldByCodeTag(sliceItemType, sliceItem, target)
	if count == 0 && !disc && !target.isDynamic {
		count = preserveItemDynamicField(sliceItemType, sliceItem, target.code, target.value)
	}

	return count, disc
}

// setScalarFieldByCodeTag lands the value in one scalar field when that field
// matches the target's code. The returned error reports a malformed field tag
// on the model, never a data problem.
func setScalarFieldByCodeTag(f reflect.Value, ft reflect.StructField, target fieldTarget) (int, bool, error) {
	if target.isDynamic {
		// Dynamic fields match by NAME, never by tag: the name encodes
		// the code losslessly, while a garbage code containing struct
		// tag syntax (a quote or backslash) renders the tag unreadable
		// and would silently drop the field.
		if ft.Name == serde.DynamicFieldName(target.code, target.order) {
			attr, _ := serde.GetFieldAttribute(ft.Tag.Get(serde.FieldTag))
			if setFieldValue(f, ft, target.value, attr) {
				return 1, false, nil
			}
		}
		return 0, false, nil
	}

	tag := ft.Tag.Get(serde.FieldTag)
	if tag == serde.Empty {
		return 0, false, nil
	}

	attr, err := serde.GetFieldAttribute(tag)
	if err != nil {
		return 0, false, err
	}
	if attr.Code != target.code {
		return 0, false, nil
	}

	// A field scoped to a newer version than the transmission
	// (e.g. an F6-only code arriving in a D0 message) must not
	// consume the value: the serializer omits such fields for
	// this version, so consuming would silently drop the
	// received bytes. Left unconsumed, the raw value preserves
	// as a dynamic field and re-emits verbatim.
	if attr.SinceVersion != serde.Empty && !ncpdp.VersionAtLeast(target.version, attr.SinceVersion) {
		return 0, false, nil
	}

	// A duplicated non-repeating field must not overwrite the
	// first occurrence: leaving it unconsumed routes the raw
	// duplicate into the segment's dynamic fields, so it
	// round-trips instead of silently collapsing into one.
	if fieldOccupied(f) {
		return 0, false, nil
	}

	// A value the typed field cannot represent (an empty or
	// blank-filled numeric/date) leaves the field unset and
	// counts as not consumed, so the caller preserves the raw
	// value as a dynamic field instead of inventing a zero.
	if setFieldValue(f, ft, target.value, attr) {
		return 1, false, nil
	}

	// Except counters: the serializer always derives
	// countfor-tagged values from the actual occurrences (the
	// derived count wins by ruling), so an unparseable counter
	// is consumed and discarded — preserving its raw value would
	// emit it beside the derived counter and spawn a phantom
	// occurrence on reparse.
	return 0, attr.CountFor != serde.Empty, nil
}

// preserveItemDynamicField keeps a field whose code belongs to a repeating
// group's element type but whose value the typed field could not represent
// (an empty or blank-filled numeric/date, or a duplicate within the
// occurrence) in THAT occurrence's dynamic fields. Segment-level preservation
// would lose which occurrence the value belongs to: an empty FU in the first
// DUR item must not re-emit inside the last one. Returns 1 when the value was
// preserved, 0 when it was not (the element type does not own the code, or
// has no dynamic field slice to hold it).
func preserveItemDynamicField(itemType reflect.Type, itemVal reflect.Value, codeTag string, propVal reflect.Value) int {
	if propVal.Kind() != reflect.String || !typeOwnsCode(itemType, codeTag) {
		return 0
	}

	if err := appendDynamicFields(itemType, itemVal, []string{codeTag}, []string{propVal.String()}); err != nil {
		return 0
	}

	return 1
}

// codeOccupied reports whether the struct already carries a value for the
// field with the given code tag, looking through embedded structs but not
// nested slices (their element types own those codes). A repeating group's
// next occurrence starts when its codes reappear, so "occupied" marks the
// boundary of the current occurrence. A value PRESERVED in the occurrence's
// dynamic fields counts as occupied too: the occurrence does carry that code,
// its value just could not land in the typed field — ignoring it would group
// the same message differently on every reparse.
func codeOccupied(structType reflect.Type, structVal reflect.Value, codeTag string) bool {
	if structVal.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < structVal.NumField(); i++ {
		f := structVal.Field(i)
		ft := structType.Field(i)

		switch ft.Type.Kind() {
		case reflect.Struct:
			if codeOccupied(ft.Type, f, codeTag) {
				return true
			}

		case reflect.Slice:
			// Nested slice fields are owned by their own element type.

		default:
			if fieldCodeMatches(ft, codeTag) && fieldOccupied(f) {
				return true
			}
		}
	}

	return dynamicFieldsContainCode(structType, structVal, codeTag)
}

// dynamicFieldsContainCode reports whether the struct's dynamic field slice
// holds a preserved value for the given code.
func dynamicFieldsContainCode(structType reflect.Type, structVal reflect.Value, codeTag string) bool {
	dynSlice, found := findDynamicSliceValue(structType, structVal, serde.FieldTag)
	if !found || dynSlice.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < dynSlice.Len(); i++ {
		ds, ok := dynSlice.Index(i).Interface().(dynamic.DynamicStruct)
		if !ok || ds.DynamicType == nil {
			continue
		}
		for j := 0; j < ds.DynamicType.NumField(); j++ {
			if code, _, ok := dynamic.ParseFieldName(ds.DynamicType.Field(j).Name); ok && code == codeTag {
				return true
			}
		}
	}

	return false
}

// fieldCodeMatches reports whether a struct field carries the given code tag.
func fieldCodeMatches(ft reflect.StructField, codeTag string) bool {
	tag := ft.Tag.Get(serde.FieldTag)
	if tag == serde.Empty {
		return false
	}

	attr, err := serde.GetFieldAttribute(tag)
	return err == nil && attr.Code == codeTag
}

// fieldOccupied reports whether a scalar struct field already holds a parsed
// value: a non-nil pointer, or a non-zero value for plain fields.
func fieldOccupied(f reflect.Value) bool {
	if f.Kind() == reflect.Pointer {
		return !f.IsNil()
	}
	return !f.IsZero()
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

// Set base field by data type. Reports whether the raw value was actually
// parsed and stored; false means the typed field cannot represent the value
// (e.g. an empty or blank-filled numeric/date), so the caller must leave the
// field unset and preserve the raw value instead of keeping a zero value.
func setField(f reflect.Value, val reflect.Value, attr serde.FieldAttribute) bool {
	if !f.CanSet() {
		return false
	}

	switch f.Interface().(type) {
	case int, int8, int16, int32, int64:
		intVal, err := strconv.Atoi(val.String())
		if err != nil {
			return false
		}
		f.Set(reflect.ValueOf(intVal))

	case uint, uint8, uint16, uint32, uint64:
		intVal, err := strconv.Atoi(val.String())
		if err != nil {
			return false
		}
		f.Set(reflect.ValueOf(uint(intVal)))

	case float32, float64:
		fs := serde.NewFieldSettings(attr)
		flt, err := fs.Unsign(val.String())
		if err != nil || flt == nil {
			return false
		}
		f.Set(reflect.ValueOf(*flt))

	case time.Time:
		if attr.Format == serde.Empty {
			return false
		}
		dt, err := time.Parse(attr.Format, val.String())
		if err != nil {
			return false
		}
		// The zero time is this library's "absent" sentinel: the serializer
		// omits zero dates. A received value that parses exactly to it (e.g.
		// the placeholder 00010101) must stay unconsumed and preserve as a
		// dynamic field, or it would silently vanish from the rebuilt message.
		if dt.IsZero() {
			return false
		}
		f.Set(reflect.ValueOf(dt))

	default:
		f.Set(val)
	}

	return true
}

// Set field value. Reports whether a value was stored (see setField); a
// pointer field stays nil when the value cannot be parsed, never pointing at
// a zero value the sender did not transmit.
func setFieldValue(f reflect.Value, ft reflect.StructField, val reflect.Value, attr serde.FieldAttribute) bool {
	switch ft.Type.Kind() {
	case reflect.Pointer:
		fv := reflect.New(ft.Type.Elem())

		if !setField(fv.Elem(), val, attr) {
			return false
		}
		if !f.CanSet() {
			return false
		}
		f.Set(fv)
		return true

	default:
		return setField(f, val, attr)
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
				// *time.Time holds a scalar field value, not a nested group:
				// it must stay nil until a value actually parses, or absent
				// dates read as the zero date (0001-01-01) and are mistaken
				// for populated fields.
				if fieldBaseType == timeType {
					break
				}
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
