package ncpdp

import "slices"

type NcpdpRecord struct {
	RawValue   string
	Segments   []NcpdpSegment
	StartIndex int
}

// Find segment for the specified ID.
func (record *NcpdpRecord) FindSegment(id string) *NcpdpSegment {
	if record == nil || len(record.Segments) == 0 {
		return nil
	}

	// Define a function that reutrns true when the segment ID matches
	equalFunc := func(s NcpdpSegment) bool {
		return s.Id == id
	}

	// Get the index of the matching segment
	i := slices.IndexFunc(record.Segments, equalFunc)

	if i == -1 {
		return nil
	}

	return &record.Segments[i]
}
