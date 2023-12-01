package ncpdp

import (
	"fmt"
	"slices"
	"strings"

	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

// Find group code in claim request
func (req *NcpdpTransaction[V]) GetGroupCode() string {
	if req == nil {
		return Empty
	}

	groupCode := Empty

	groupField := req.FindFirstField(INSURANCE_SEGMENT_ID, GROUP_CODE_FIELD_ID, -1)
	if groupField != nil {
		groupCode = groupField.Value
	}

	return groupCode
}

// Generate reversal transmission.
func (req *NcpdpTransaction[V]) GenerateReversal(requiredFields []string) (*NcpdpTransaction[RequestHeader], error) {
	if req == nil {
		return nil, fmt.Errorf("request is null")
	}

	// Copy the raw data from the original request
	reversal := NewTransactionRequest(req.RawValue)

	// Parse the copy
	err := reversal.ParseNcpdp()
	if err != nil {
		return nil, err
	}

	// Update transaction code and rebuild header
	reversal.Header.Value.TransactionCode = REVERSAL
	reversal.BuildNcpdpHeader()

	// Compress raw data
	compressed := compress(reversal.RawValue, requiredFields)
	reversal.RawValue = compressed

	// Parse complete reversal
	err = reversal.ParseNcpdp()
	if err != nil {
		return nil, err
	}

	return &reversal, nil
}

// Compress empty and unnecessary values based on field list.
func compress(rawValue string, requiredFields []string) string {
	cleanFields := cleanFieldIds(requiredFields)

	rawValue = compressFieldRecursive(rawValue, cleanFields, 0)
	rawValue = removeEmptyFields(rawValue)
	rawValue = removeEmptySegments(rawValue)
	rawValue = removeEmptyGroups(rawValue)

	return rawValue
}

// If fields are formatted with dashes, only keep the part after the dash.
//
//	Ex: 201-B1 -> B1
//	Ex: N6 -> N6
func cleanFieldIds(requiredFields []string) []string {
	cleanedFields := []string{}

	for _, field := range requiredFields {
		dashIndex := strings.Index(field, "-")
		fieldCode := field[dashIndex+1:]

		if !slices.Contains(cleanedFields, fieldCode) {
			cleanedFields = append(cleanedFields, fieldCode)
		}
	}

	slices.Sort(cleanedFields)

	return cleanedFields
}

// Compress fields from raw data.
func compressFieldRecursive(rawValue string, requiredFields []string, startIndex int) string {
	if len(rawValue) == 0 || len(requiredFields) == 0 || startIndex >= len(rawValue) {
		return rawValue
	}

	fieldStartIndex := stringutils.IndexOfAny(rawValue, startIndex, []byte{FIELD})

	// No more fields
	if fieldStartIndex == -1 {
		return rawValue
	}

	fieldId := stringutils.Substring(rawValue, fieldStartIndex+1, ID_LENGTH)

	keepIt := slices.Contains(requiredFields, fieldId)

	// Keep field if it's in the required list or it's a segment identifier
	if keepIt || fieldId == SEGMENT_FIELD_ID {
		return compressFieldRecursive(rawValue, requiredFields, fieldStartIndex+1)
	}

	fieldEndIndex := stringutils.IndexOfAny(rawValue, fieldStartIndex+1, []byte{FIELD, SEGMENT, GROUP, ETX})
	if fieldEndIndex == -1 {
		fieldEndIndex = len(rawValue)
	}

	// Remove the field
	newVal := fmt.Sprintf("%v%v", stringutils.Substring(rawValue, 0, fieldStartIndex), stringutils.Substring(rawValue, fieldEndIndex, -1))

	return compressFieldRecursive(newVal, requiredFields, fieldStartIndex)
}

// Remove empty fields from raw data.
func removeEmptyFields(rawValue string) string {
	if len(rawValue) == 0 {
		return rawValue
	}

	startIndex := 0
	finished := false

	for ok := true; ok; ok = !finished {
		fieldStartIndex := stringutils.IndexOfAny(rawValue, startIndex, []byte{FIELD})

		// No more fields
		if fieldStartIndex == -1 {
			finished = true
		}

		fieldEndIndex := stringutils.IndexOfAny(rawValue, fieldStartIndex+1, []byte{FIELD, SEGMENT, GROUP, ETX})
		if fieldEndIndex == -1 {
			fieldEndIndex = len(rawValue)
		}

		fieldValueLength := fieldEndIndex - fieldStartIndex + 3

		contents := strings.TrimSpace(stringutils.Substring(rawValue, fieldStartIndex+3, fieldValueLength))
		contents = stringutils.TrimAll(contents, []byte{FIELD, SEGMENT, GROUP, ETX})

		if len(contents) == 0 {
			// Remove empty field, only contains field id
			rawValue = fmt.Sprintf("%v%v", stringutils.Substring(rawValue, 0, fieldStartIndex), stringutils.Substring(rawValue, fieldEndIndex, -1))
		} else {
			startIndex = fieldStartIndex + 1
		}
	}

	return rawValue
}

// Remove empty segments from raw data.
func removeEmptySegments(rawValue string) string {
	if len(rawValue) == 0 {
		return rawValue
	}

	startIndex := 0
	finished := false

	for ok := true; ok; ok = !finished {
		segmentStartIndex := stringutils.IndexOfAny(rawValue, startIndex, []byte{SEGMENT})

		// No more segments
		if segmentStartIndex == -1 {
			finished = true
		}

		segmentEndIndex := stringutils.IndexOfAny(rawValue, segmentStartIndex+1, []byte{SEGMENT, GROUP, ETX})
		if segmentEndIndex == -1 {
			segmentEndIndex = len(rawValue)
		}

		segmentValueLength := segmentEndIndex - segmentStartIndex + 1

		contents := strings.TrimSpace(stringutils.Substring(rawValue, segmentStartIndex+1, segmentValueLength))
		contents = stringutils.TrimAll(contents, []byte{FIELD, SEGMENT, GROUP, ETX})

		if len(contents) <= 4 {
			// Remove empty segment, or segment containing only the segment ID
			rawValue = fmt.Sprintf("%v%v", stringutils.Substring(rawValue, 0, segmentStartIndex), stringutils.Substring(rawValue, segmentEndIndex, -1))
		} else {
			startIndex = segmentStartIndex + 1
		}
	}

	return rawValue
}

// Remove empty record groups from raw data.
func removeEmptyGroups(rawValue string) string {
	if len(rawValue) == 0 {
		return rawValue
	}

	startIndex := 0
	finished := false

	for ok := true; ok; ok = !finished {
		groupStartIndex := stringutils.IndexOfAny(rawValue, startIndex, []byte{GROUP})

		// No more groups
		if groupStartIndex == -1 {
			finished = true
		}

		groupEndIndex := stringutils.IndexOfAny(rawValue, groupStartIndex+1, []byte{GROUP, ETX})
		if groupEndIndex == -1 {
			groupEndIndex = len(rawValue)
		}

		groupValueLength := groupEndIndex - groupStartIndex + 1

		contents := strings.TrimSpace(stringutils.Substring(rawValue, groupStartIndex+1, groupValueLength))
		contents = stringutils.TrimAll(contents, []byte{FIELD, SEGMENT, GROUP, ETX})

		if len(contents) <= 1 {
			// Remove empty group
			rawValue = fmt.Sprintf("%v%v", stringutils.Substring(rawValue, 0, groupStartIndex), stringutils.Substring(rawValue, groupEndIndex, -1))
		} else {
			startIndex = groupStartIndex + 1
		}
	}

	return rawValue
}
