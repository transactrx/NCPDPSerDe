package ncpdp

import (
	"slices"
)

type SegmentId struct {
	Raw string  `field:"code=rawsegment"`
	Id  *string `field:"code=AM"`
}

type NcpdpSegment struct {
	Id         string
	RawValue   string
	Fields     []NcpdpField
	StartIndex int
}

// Find field for the specified ID.
func (seg *NcpdpSegment) FindFirstField(id string) *NcpdpField {
	if seg == nil || len(seg.Fields) == 0 {
		return nil
	}

	// Define a function that reutrns true when the field ID matches
	equalFunc := func(f NcpdpField) bool {
		return f.Id == id
	}

	// Get the index of the matching segment
	i := slices.IndexFunc(seg.Fields, equalFunc)

	if i == -1 {
		return nil
	}

	return &seg.Fields[i]
}

// Find all fields for the specified ID.
func (seg *NcpdpSegment) FindAllFields(id string) []NcpdpField {
	fields := []NcpdpField{}

	if seg == nil || len(seg.Fields) == 0 {
		return fields
	}

	for _, field := range seg.Fields {
		if field.Id == id {
			fields = append(fields, field)
		}
	}

	return fields
}

// Append field.
func (seg *NcpdpSegment) AppendField(id, value string) {
	if seg == nil || id == Empty || value == Empty {
		return
	}

	seg.Fields = append(seg.Fields, *NewField(id, value))
}

// Append field copy.
func (seg *NcpdpSegment) AppendFieldCopy(field *NcpdpField) {
	if seg == nil || field == nil {
		return
	}

	seg.AppendField(field.Id, field.Value)

	seg.Fields = append(seg.Fields, *NewField(field.Id, field.Value))
}
