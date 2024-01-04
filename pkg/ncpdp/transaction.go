package ncpdp

import (
	"fmt"
	"slices"
	"strings"
	"time"

	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

type NcpdpTransaction[V RequestHeader | ResponseHeader] struct {
	Header   NcpdpHeader[V]
	RawValue string
	Segments []NcpdpSegment
	Records  []NcpdpRecord
	Created  time.Time
}

func NewTransactionRequest(rawData string) NcpdpTransaction[RequestHeader] {
	tran := NcpdpTransaction[RequestHeader]{
		RawValue: rawData,
		Created:  time.Now().UTC(),
	}

	return tran
}

func NewTransactionResponse(rawData string) NcpdpTransaction[ResponseHeader] {
	tran := NcpdpTransaction[ResponseHeader]{
		RawValue: rawData,
		Created:  time.Now().UTC(),
	}

	return tran
}

// Rebuid the complete NCPDP raw value
func (tran *NcpdpTransaction[V]) BuildNcpdp() error {
	if tran == nil {
		return fmt.Errorf("transaction is null")
	}

	err := tran.BuildNcpdpHeader()
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	builder.WriteString(tran.Header.RawValue)
	builder.WriteString(tran.buildBody())

	tran.RawValue = builder.String()

	return nil
}

func (tran *NcpdpTransaction[V]) buildBody() string {
	if tran == nil {
		return Empty
	}

	builder := strings.Builder{}

	// Build shared segments
	builder.WriteString(buildSegments(tran.Segments))

	//Build records
	builder.WriteString(buildRecords(tran.Records))

	//ETX on the end
	builder.WriteByte(ETX)

	return builder.String()
}

// Build raw value for segments.
func buildSegments(segments []NcpdpSegment) string {
	if len(segments) == 0 {
		return Empty
	}

	builder := strings.Builder{}

	for _, segment := range segments {
		// Write segment separator
		builder.WriteByte(SEGMENT)

		// Write fields
		for _, field := range segment.Fields {
			builder.WriteByte(FIELD)
			builder.WriteString(fmt.Sprintf("%v%v", field.Id, field.Value))
		}
	}

	return builder.String()
}

// Build raw value for segments.
func buildRecords(records []NcpdpRecord) string {
	if len(records) == 0 {
		return Empty
	}

	builder := strings.Builder{}

	for _, record := range records {
		// Write group separator
		builder.WriteByte(GROUP)

		// Write segments
		builder.WriteString(buildSegments(record.Segments))
	}

	return builder.String()
}

// Rebuid the header's raw value
func (tran *NcpdpTransaction[V]) BuildNcpdpHeader() error {
	if tran == nil {
		return nil
	}

	oldHeader := tran.Header.RawValue

	err := tran.Header.BuildNcpdpHeader()
	if err != nil {
		return err
	}

	newHeader := tran.Header.RawValue
	tran.RawValue = strings.Replace(tran.RawValue, oldHeader, newHeader, 1)

	return nil
}

// Find segment for the specified ID.
func (tran *NcpdpTransaction[V]) FindSegment(id string) *NcpdpSegment {
	if tran == nil || len(tran.Segments) == 0 {
		return nil
	}

	// Define a function that reutrns true when the segment ID matches
	equalFunc := func(s NcpdpSegment) bool {
		return s.Id == id
	}

	// Get the index of the matching segment
	i := slices.IndexFunc(tran.Segments, equalFunc)

	if i == -1 {
		return nil
	}

	return &tran.Segments[i]
}

// Find segment for the specified ID.
// Record index of -1 will search shared segments.
func (tran *NcpdpTransaction[V]) FindFirstField(segmentId, fieldId string, recordIndex int) *NcpdpField {
	if tran == nil {
		return nil
	}

	var segment *NcpdpSegment = nil

	if recordIndex < 0 {
		// Shared segment
		segment = tran.FindSegment(segmentId)
	} else if recordIndex < len(tran.Records) {
		// Non-shared segment
		segment = tran.Records[recordIndex].FindSegment(segmentId)
	}

	if segment == nil {
		return nil
	}

	return segment.FindFirstField(fieldId)
}

// Insert field into claim body.
func (tran *NcpdpTransaction[V]) InsertField(recordIndex int, spec *NcpdpSegmentSpecification, fieldId string, fieldValue interface{}, settings *FieldSettings) error {
	if tran == nil || spec == nil || fieldId == Empty || fieldValue == nil {
		return nil
	}

	strField := fmt.Sprint(fieldValue)
	if strField == Empty {
		return nil
	}

	if settings == nil {
		settings = &FieldSettings{}
	}

	// Get segment to insert field into
	var segment *NcpdpSegment = nil
	if recordIndex < 0 {
		// Shared segment
		segment = tran.FindSegment(spec.SegmentCode)
	} else if recordIndex < len(tran.Records) {
		// Non-shared segment
		segment = tran.Records[recordIndex].FindSegment(spec.SegmentCode)
	}

	if segment == nil {
		//Need to add missing segments? later?
		return fmt.Errorf("unable to insert field. Segment %v not found in record %v", spec.SegmentCode, recordIndex)
	}

	// If field does not repeat, ensure it's not already there before inserting
	if !settings.Repeating {
		existingField := segment.FindFirstField(fieldId)
		if existingField != nil {
			return nil
		}
	}

	//Get NCPDP field specification to determine insert order
	fieldSpec := spec.FindField(fieldId)
	if fieldSpec == nil {
		return nil
	}

	insertOrder := fieldSpec.SortOrder

	//Loop through existing fields to determine insertion point based on NCPDP spec sort order
	inserted := false
	for _, field := range segment.Fields {
		specField := spec.FindField(field.Id)
		if specField == nil {
			continue
		}

		currentOrder := specField.SortOrder
		if currentOrder < insertOrder {
			continue
		}

		//Insert field when insert order >= to current order
		tran.insertField(fieldId, fieldValue, settings, field.StartIndex)
		inserted = true
		break
	}

	//Insert at end of segment
	if !inserted {
		tran.insertField(fieldId, fieldValue, settings, segment.StartIndex+len(segment.RawValue))
	}

	return nil
}

// Insert field value into raw data
func (tran *NcpdpTransaction[V]) insertField(fieldId string, fieldValue interface{}, settings *FieldSettings, insertAt int) {
	if tran == nil {
		return
	}

	rawFieldVal := fmt.Sprintf("%v%v%v",
		FIELD,
		fieldId,
		settings.convertFieldValueToString(fieldValue))

	newVal := fmt.Sprintf("%v%v%v",
		stringutils.Substring(tran.RawValue, 0, insertAt),
		rawFieldVal,
		stringutils.Substring(tran.RawValue, insertAt, -1))

	tran.RawValue = newVal
}

// Parse NCPDP claim data.
func (tran *NcpdpTransaction[V]) ParseNcpdp() error {
	if tran == nil {
		return nil
	}

	// Parse header
	header := NcpdpHeader[V]{
		RawValue: tran.RawValue,
	}

	err := header.ParseNcpdpHeader()
	if err != nil {
		return err
	}

	tran.Header = header

	// Parse the rest
	tran.parseBody()

	return nil
}

// Parse NCPDP claim body.
func (tran *NcpdpTransaction[V]) parseBody() {
	if tran == nil {
		return
	}

	for i := tran.Header.Size - 1; i < len(tran.RawValue); i++ {
		ch := tran.RawValue[i]

		switch ch {
		case SEGMENT:
			rawSegment := stringutils.SplitFirst(tran.RawValue, i, []byte{SEGMENT, GROUP})
			segment := parseSegment(rawSegment, i)
			if segment != nil {
				tran.Segments = append(tran.Segments, *segment)
			}

			//Adjust index to skip to the end of the segment.
			i += len(rawSegment) - 1

		case GROUP:
			rawGroup := stringutils.SplitFirst(tran.RawValue, i, []byte{GROUP})
			group := parseGroup(rawGroup, i)
			if group != nil {
				tran.Records = append(tran.Records, *group)
			}

			//Adjust index to skip to the end of the record.
			i += len(rawGroup) - 1
		}
	}
}

// Parse NCPDP group.
func parseGroup(rawData string, startIndex int) *NcpdpRecord {
	record := NcpdpRecord{
		RawValue:   rawData,
		StartIndex: startIndex,
	}

	for i := 0; i < len(rawData); i++ {
		ch := rawData[i]

		switch ch {
		case SEGMENT:
			rawSegment := stringutils.SplitFirst(rawData, i, []byte{SEGMENT, GROUP})
			segment := parseSegment(rawSegment, record.StartIndex+i)

			if segment != nil {
				record.Segments = append(record.Segments, *segment)
			}

			//Adjust index to skip to the end of the segment.
			i += len(rawSegment) - 1
		}
	}

	return &record
}

// Parse NCPDP segment.
func parseSegment(rawData string, startIndex int) *NcpdpSegment {
	//Handle empty segment by checking for all separators instead of just FS.
	firstFs := stringutils.IndexOfAny(rawData, ID_LENGTH, []byte{FIELD, SEGMENT, GROUP, ETX})

	//Skip empty segment
	if firstFs < 0 {
		return nil
	}

	//Segment ID is the first field in the segment.
	segmentId := stringutils.Substring(rawData, ID_LENGTH, firstFs-ID_LENGTH)

	segment := NcpdpSegment{
		Id:         segmentId,
		RawValue:   rawData,
		StartIndex: startIndex,
	}

	for i := 0; i < len(rawData); i++ {
		ch := rawData[i]

		switch ch {
		case FIELD:
			rawField := stringutils.SplitFirst(rawData, i, []byte{FIELD, SEGMENT, GROUP, ETX})

			if rawField != Empty {
				field := NcpdpField{
					Id:         stringutils.Substring(rawField, 1, ID_LENGTH),
					Value:      stringutils.Substring(rawField, ID_LENGTH+1, -1),
					RawValue:   rawField,
					StartIndex: segment.StartIndex + i,
				}

				segment.Fields = append(segment.Fields, field)
			}

			//Adjust index to skip to the end of the field.
			i += len(rawField) - 1
		}
	}

	return &segment
}
