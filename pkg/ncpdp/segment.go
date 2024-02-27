package ncpdp

import (
	"slices"
)

type SegmentId struct {
	Raw *string `field:"code=rawsegment"`
	Id  *string `field:"code=AM,order=1"`
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

// Find first index for the specified ID, -1 when not found
func (seg *NcpdpSegment) FindFirstIndex(id string) int {
	if seg == nil || len(seg.Fields) == 0 {
		return -1
	}

	// Define a function that reutrns true when the field ID matches
	equalFunc := func(f NcpdpField) bool {
		return f.Id == id
	}

	// Get the index of the matching segment
	return slices.IndexFunc(seg.Fields, equalFunc)
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

// Insert field.
func (seg *NcpdpSegment) InsertField(id, value string, index int) {
	if seg == nil || id == Empty || value == Empty {
		return
	}

	seg.Fields = slices.Insert(seg.Fields, index, *NewField(id, value))
}

// Delete field.
func (seg *NcpdpSegment) DeleteField(id string) {
	if seg == nil || id == Empty {
		return
	}

	index := seg.FindFirstIndex(id)
	if index < 0 {
		return
	}

	seg.Fields = append(seg.Fields[:index], seg.Fields[index+1:]...)
}
