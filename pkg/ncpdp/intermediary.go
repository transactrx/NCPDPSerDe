package ncpdp

import "strings"

// One repetition of the F6 intermediary segment (AM19).
type intermediaryEntry struct {
	typeCode  string // 8H
	qualifier string // 8K
	id        string // 8M
}

// IsServiceBypassed reports whether the customer bypassed the service, which
// the request indicates by carrying an intermediary bypass code matching the
// supplied ID.
//
// D.0 encodes this in the claim segment (AM07): EW (qualifier) = 01 and
// EX (auth code) = id. F6 moved it to the intermediary segment (AM19), where
// entries repeat: 8H (type code) = 01, 8K (ID qualifier) = 06 and
// 8M (ID) = id must all appear on the same entry.
//
// Both encodings are checked, so callers do not need to dispatch on version.
// Returns true when any record (D.0 batch) or any AM19 entry matches.
// Numeric code variants are accepted ("1" matches "01", "6" matches "06").
func (req *NcpdpTransaction[V]) IsServiceBypassed(id string) bool {
	id = strings.TrimSpace(id)

	if req == nil || id == Empty {
		return false
	}

	// RecordCount and FindSegmentInRecord fall back to transaction-level
	// segments when there are no records (F6 removed the group separator).
	for i := 0; i < req.RecordCount(); i++ {
		if claimHasIntermediaryAuth(req.FindSegmentInRecord(i, CLAIM_SEGMENT_ID), id) {
			return true
		}

		if intermediarySegmentHasId(req.FindSegmentInRecord(i, INTERMEDIARY_SEGMENT_ID), id) {
			return true
		}
	}

	return false
}

// D.0: EX carries the authorization code; EW must qualify it as an
// intermediary authorization ID.
func claimHasIntermediaryAuth(claim *NcpdpSegment, id string) bool {
	if claim == nil {
		return false
	}

	authCode := claim.FindFirstField(INTERMEDIARY_AUTH_CODE_FIELD_ID)
	if authCode == nil || strings.TrimSpace(authCode.Value) != id {
		return false
	}

	qualifier := claim.FindFirstField(INTERMEDIARY_AUTH_CODE_QUALIFIER_FIELD_ID)

	return qualifier != nil && codesEqual(qualifier.Value, INTERMEDIARY_AUTH_QUALIFIER_D0)
}

// F6: the type code, qualifier and ID must match on the same AM19 entry.
func intermediarySegmentHasId(seg *NcpdpSegment, id string) bool {
	for _, entry := range groupIntermediaryEntries(seg) {
		if codesEqual(entry.typeCode, INTERMEDIARY_TYPE_CODE_F6) &&
			codesEqual(entry.qualifier, INTERMEDIARY_ID_QUALIFIER_F6) &&
			strings.TrimSpace(entry.id) == id {
			return true
		}
	}

	return false
}

// Walk AM19 fields in wire order and split them into repeating entries.
// A new entry starts when a field ID repeats, so entries that omit leading
// optional fields still group correctly.
func groupIntermediaryEntries(seg *NcpdpSegment) []intermediaryEntry {
	if seg == nil || len(seg.Fields) == 0 {
		return nil
	}

	entries := []intermediaryEntry{}
	current := intermediaryEntry{}
	seen := map[string]bool{}

	flush := func() {
		if len(seen) > 0 {
			entries = append(entries, current)
			current = intermediaryEntry{}
			seen = map[string]bool{}
		}
	}

	for _, field := range seg.Fields {
		target := entryField(&current, field.Id)
		if target == nil {
			continue
		}

		if seen[field.Id] {
			flush()
			target = entryField(&current, field.Id)
		}

		seen[field.Id] = true
		*target = field.Value
	}

	flush()

	return entries
}

// Map a field ID to its slot in the entry, nil for fields that do not
// participate in matching (8G count, 8J, 8N, 8U).
func entryField(entry *intermediaryEntry, fieldId string) *string {
	switch fieldId {
	case INTERMEDIARY_TYPE_CODE_FIELD_ID:
		return &entry.typeCode
	case INTERMEDIARY_ID_QUALIFIER_FIELD_ID:
		return &entry.qualifier
	case INTERMEDIARY_ID_FIELD_ID:
		return &entry.id
	default:
		return nil
	}
}

// Compare NCPDP code values, accepting numeric variants ("1" matches "01").
func codesEqual(value, canonical string) bool {
	value = strings.TrimSpace(value)

	if value == canonical {
		return true
	}

	trimmed := strings.TrimLeft(value, "0")

	return trimmed != Empty && trimmed == strings.TrimLeft(canonical, "0")
}
