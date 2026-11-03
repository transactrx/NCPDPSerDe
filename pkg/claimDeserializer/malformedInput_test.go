package claimdeserializer

import (
	"strings"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// A malformed fixed-width header cannot be parsed correctly — every byte of
// length drift shifts the remaining header fields — so it must refuse loudly
// instead of silently garbling them. All shapes below were seen in production:
// a 46-byte header missing the vendor id, a 55-byte header one byte short, a
// PCN padded with multi-byte non-breaking spaces, and 154-byte 5X switch
// headers. The check is structural (the byte after the header must be a
// segment or group separator); header field VALUES are deliberately not
// validated — this library is not a validator.
func Test_MalformedHeaderRefused(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	body := ss + fs + "AM04" + fs + "C2ID12345" +
		ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "E103" + fs + "D700172240780" + fs + "DJ1"

	tests := []struct {
		name    string
		header  string
		wantErr bool
	}{
		{"well-formed 56-byte header", "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001", false},
		{"46-byte header (no vendor id)", "999999D0B1TESTPCN   1011234567893     20261001", true},
		{"55-byte header (provider id one short)", "999999D0B1TESTPCN   1011234567893    20261001VENDOR0001", true},
		{"57-byte header", "999999D0B1TESTPCN   1011234567893     20261001VENDOR00012", true},
		{"non-breaking spaces in PCN", "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001", true},
		{"154-byte 5X switch header", "0047665XB1USFLU     1011811921893     20261020          SHAUN     18104    0001180000037660380777040160 RB55YVUG  202710203EFG6471899   00001PA           ", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj := request.Billing{}
			err := DeserializeType(test.header+body, &obj)

			if test.wantErr && err == nil {
				t.Errorf("Malformed header accepted silently. Wanted: error   Got: nil")
			}
			if !test.wantErr && err != nil {
				t.Errorf("Well-formed header refused. Wanted: nil   Got: %v", err)
			}
		})
	}
}

// A transmission that parses to zero segments (a header followed only by bare
// separators or id-only fragments) cannot be re-serialized into anything this
// parser would accept — the rebuild would be a header-only string the parser
// refuses — so it refuses loudly up front. Found by FuzzRoundTrip.
func Test_TransmissionWithoutSegmentsRefused(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	tests := []struct {
		name string
		raw  string
	}{
		{"request header and segment separator only", "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" + ss},
		{"response header and segment separator only", "D0B11A011234567893     20261001" + ss},
		{"request header and group separator only", "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" + gs + ss},
		{"id-only segment is the whole transmission", "D0B11A011234567893     20261001" + ss + fs + "AM22"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Deserialize(test.raw); err == nil {
				t.Errorf("Transmission without segments accepted silently. Wanted: error   Got: nil")
			}
		})
	}
}

// A segment with no AM segment id at all cannot be attributed to anything:
// its data would be lost or misassigned, so it must refuse loudly instead of
// silently dropping the rest of the part.
func Test_SegmentWithoutIdRefused(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" +
		ss + fs + "C2ID12345" + fs + "CCJOHN" +
		ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "E103" + fs + "D700172240780" + fs + "DJ1"

	obj := request.Billing{}
	err := DeserializeType(raw, &obj)
	if err == nil {
		t.Fatal("Segment without an id accepted silently. Wanted: error   Got: nil")
	}
	if !strings.Contains(err.Error(), "segment id") {
		t.Errorf("Error does not name the problem. Got: %v", err)
	}
}
