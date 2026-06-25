package ncpdp

import (
	"sort"
	"testing"
)

const responseHeaderPaidD0 = "D0B11A011234567893     20260611"
const responseHeaderRejectD0 = "D0B11A011234567893     20260611"
const responseHeaderPaidF6 = "F6B11A011234567893     20260611"
const responseHeaderRejectF6 = "F6B11A011234567893     20260611"

// AM21 status segment for a paid response. Each segment starts with
// SEGMENT(0x1E) + FIELD(0x1C) + segment ID, then FIELD(0x1C) + fieldId + value
// for every field.
const statusSegmentPaid = "\x1e\x1cAM21\x1cANP"

// AM21 status segment with two reject codes and two qualified additional
// messages. Built so D0 / F6 share the same segment content — only the
// surrounding group separator differs.
const statusSegmentReject = "\x1e\x1cAM21\x1cANR\x1cFB85\x1cFB99\x1cUH01\x1cFQfirst message\x1cUH02\x1cFQsecond message"

func parseResponse(t *testing.T, raw string) *NcpdpTransaction[ResponseHeader] {
	t.Helper()
	tran := NewTransactionResponse(raw)
	if err := tran.ParseNcpdp(); err != nil {
		t.Fatalf("ParseNcpdp: %v", err)
	}
	return &tran
}

func sortedCopy(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	a, b = sortedCopy(a), sortedCopy(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type responseHelperCase struct {
	name              string
	raw               string
	wantStatus        string
	wantIsPaid        bool
	wantIsRejected    bool
	wantRejectCodes   []string
	wantMessages      map[string]string
	wantRecordsLen    int
	wantSharedSegLen  int // includes AM21 only when transaction has no records
	statusInSharedSeg bool
}

// D0 places AM21 inside a record (after a 0x1D group separator); F6 emits no
// group separator so AM21 ends up at the transaction level. The helpers must
// work for both shapes.
var responseHelperCases = []responseHelperCase{
	{
		name:              "D0 paid",
		raw:               responseHeaderPaidD0 + "\x1d" + statusSegmentPaid,
		wantStatus:        PAID_STATUS,
		wantIsPaid:        true,
		wantIsRejected:    false,
		wantRejectCodes:   []string{},
		wantMessages:      map[string]string{},
		wantRecordsLen:    1,
		wantSharedSegLen:  0,
		statusInSharedSeg: false,
	},
	{
		name:              "D0 rejected with two reject codes and two messages",
		raw:               responseHeaderRejectD0 + "\x1d" + statusSegmentReject,
		wantStatus:        REJECTED_STATUS,
		wantIsPaid:        false,
		wantIsRejected:    true,
		wantRejectCodes:   []string{"85", "99"},
		wantMessages:      map[string]string{"01": "first message", "02": "second message"},
		wantRecordsLen:    1,
		wantSharedSegLen:  0,
		statusInSharedSeg: false,
	},
	{
		name:              "F6 paid",
		raw:               responseHeaderPaidF6 + statusSegmentPaid,
		wantStatus:        PAID_STATUS,
		wantIsPaid:        true,
		wantIsRejected:    false,
		wantRejectCodes:   []string{},
		wantMessages:      map[string]string{},
		wantRecordsLen:    0,
		wantSharedSegLen:  1,
		statusInSharedSeg: true,
	},
	{
		name:              "F6 rejected with two reject codes and two messages",
		raw:               responseHeaderRejectF6 + statusSegmentReject,
		wantStatus:        REJECTED_STATUS,
		wantIsPaid:        false,
		wantIsRejected:    true,
		wantRejectCodes:   []string{"85", "99"},
		wantMessages:      map[string]string{"01": "first message", "02": "second message"},
		wantRecordsLen:    0,
		wantSharedSegLen:  1,
		statusInSharedSeg: true,
	},
}

func TestResponseHelpersAcrossVersions(t *testing.T) {
	for _, tc := range responseHelperCases {
		t.Run(tc.name, func(t *testing.T) {
			tran := parseResponse(t, tc.raw)
			verifyParsedShape(t, tran, tc)
			verifyStatusFlags(t, tran, tc)
			verifyRejectsAndMessages(t, tran, tc)
		})
	}
}

func verifyParsedShape(t *testing.T, tran *NcpdpTransaction[ResponseHeader], tc responseHelperCase) {
	t.Helper()
	if len(tran.Records) != tc.wantRecordsLen {
		t.Errorf("Records length mismatch. Wanted: %v   Got: %v", tc.wantRecordsLen, len(tran.Records))
	}
	if len(tran.Segments) != tc.wantSharedSegLen {
		t.Errorf("Shared segments length mismatch. Wanted: %v   Got: %v", tc.wantSharedSegLen, len(tran.Segments))
	}
	if tc.statusInSharedSeg && tran.FindSegment(RESPONSE_STATUS_SEGMENT_ID) == nil {
		t.Error("Expected AM21 to live in shared segments for F6 transmissions")
	}
}

func verifyStatusFlags(t *testing.T, tran *NcpdpTransaction[ResponseHeader], tc responseHelperCase) {
	t.Helper()
	if got := tran.Status(); got != tc.wantStatus {
		t.Errorf("Status mismatch. Wanted: %q   Got: %q", tc.wantStatus, got)
	}
	if got := tran.IsPaid(); got != tc.wantIsPaid {
		t.Errorf("IsPaid mismatch. Wanted: %v   Got: %v", tc.wantIsPaid, got)
	}
	if got := tran.IsRejected(); got != tc.wantIsRejected {
		t.Errorf("IsRejected mismatch. Wanted: %v   Got: %v", tc.wantIsRejected, got)
	}
}

func verifyRejectsAndMessages(t *testing.T, tran *NcpdpTransaction[ResponseHeader], tc responseHelperCase) {
	t.Helper()
	if got := tran.GetRejectCodes(); !equalStringSlices(got, tc.wantRejectCodes) {
		t.Errorf("GetRejectCodes mismatch. Wanted: %v   Got: %v", tc.wantRejectCodes, got)
	}

	gotMessages := tran.GetAdditionalMessages()
	if len(gotMessages) != len(tc.wantMessages) {
		t.Errorf("GetAdditionalMessages length mismatch. Wanted: %v   Got: %v", tc.wantMessages, gotMessages)
	}
	for k, want := range tc.wantMessages {
		got, ok := gotMessages[k]
		if !ok || got != want {
			t.Errorf("GetAdditionalMessages[%q] mismatch. Wanted: %q   Got: %q (present=%v)", k, want, got, ok)
		}
	}
}

// IsStatusOf is the shared shortcut behind IsPaid/IsRejected. Exercising it
// directly guards against future refactors that move the F6 fallback elsewhere.
func TestResponseIsStatusOfAcrossVersions(t *testing.T) {
	cases := []struct {
		name  string
		raw   string
		probe string
		want  bool
	}{
		{"D0 paid matches P", responseHeaderPaidD0 + "\x1d" + statusSegmentPaid, PAID_STATUS, true},
		{"D0 paid does not match R", responseHeaderPaidD0 + "\x1d" + statusSegmentPaid, REJECTED_STATUS, false},
		{"F6 paid matches P", responseHeaderPaidF6 + statusSegmentPaid, PAID_STATUS, true},
		{"F6 paid does not match R", responseHeaderPaidF6 + statusSegmentPaid, REJECTED_STATUS, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tran := parseResponse(t, tc.raw)
			if got := tran.IsStatusOf(tc.probe); got != tc.want {
				t.Errorf("IsStatusOf(%q) mismatch. Wanted: %v   Got: %v", tc.probe, tc.want, got)
			}
		})
	}
}

// Defensive: helpers must not panic when there is no AM21 segment at all.
func TestResponseHelpersWithoutStatusSegment(t *testing.T) {
	for _, version := range []string{D0, F6} {
		t.Run(version, func(t *testing.T) {
			raw := version + "B11A011234567893     20260611"
			tran := parseResponse(t, raw)
			verifyEmptyResponseHelpers(t, tran)
		})
	}
}

func verifyEmptyResponseHelpers(t *testing.T, tran *NcpdpTransaction[ResponseHeader]) {
	t.Helper()
	if got := tran.Status(); got != Empty {
		t.Errorf("Status mismatch. Wanted empty   Got: %q", got)
	}
	if tran.IsPaid() {
		t.Error("IsPaid should be false when no status segment is present")
	}
	if tran.IsRejected() {
		t.Error("IsRejected should be false when no status segment is present")
	}
	if got := tran.GetRejectCodes(); len(got) != 0 {
		t.Errorf("GetRejectCodes mismatch. Wanted empty   Got: %v", got)
	}
	if got := tran.GetAdditionalMessages(); len(got) != 0 {
		t.Errorf("GetAdditionalMessages mismatch. Wanted empty   Got: %v", got)
	}
}
