package ncpdp

import (
	"slices"
	"strings"
)

type NcpdpSegmentSpecification struct {
	SegmentCode string                    `json:"segmentCode"`
	Fields      []NcpdpFieldSpecification `json:"fields"`
}

type NcpdpFieldSpecification struct {
	FieldCode string `json:"fieldCode"`
	SortOrder int    `json:"sortOrder"`
}

// Find the specified field by ID
func (seg *NcpdpSegmentSpecification) FindField(fieldCode string) *NcpdpFieldSpecification {
	if seg == nil || len(seg.Fields) == 0 {
		return nil
	}

	// Define a function that reutrns true when the field ID matches
	equalFunc := func(f NcpdpFieldSpecification) bool {
		return strings.HasSuffix(f.FieldCode, fieldCode)
	}

	// Get the index of the matching element
	i := slices.IndexFunc(seg.Fields, equalFunc)

	if i == -1 {
		return nil
	}

	return &seg.Fields[i]
}
