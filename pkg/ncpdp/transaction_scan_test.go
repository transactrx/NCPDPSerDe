package ncpdp

import (
	"testing"
)

// TestScanRawDataForFieldValues_EmptyField is a regression test: an EMPTY field
// (field code immediately followed by another separator) must return "" — not the
// remainder of the transmission. Before the endIndex >= 0 fix, an empty CM field
// returned "\x1cCNMIAMI..." (control chars + every following field), which
// manufactured garbage values in downstream consumers.
func TestScanRawDataForFieldValues_EmptyField(t *testing.T) {
	fs := string(FIELD)
	ss := string(SEGMENT)
	gs := string(GROUP)
	etx := string(ETX)

	tests := []struct {
		name     string
		rawValue string
		field    string
		want     []string
	}{
		{
			name:     "empty field followed by another field",
			rawValue: fs + "CM" + fs + "CN" + "MIAMI",
			field:    "CM",
			want:     []string{""},
		},
		{
			name:     "empty field at segment boundary",
			rawValue: fs + "CM" + ss + fs + "CP" + "33101",
			field:    "CM",
			want:     []string{""},
		},
		{
			name:     "empty field at group boundary",
			rawValue: fs + "CQ" + gs + fs + "D7" + "00000000001",
			field:    "CQ",
			want:     []string{""},
		},
		{
			name:     "empty field before ETX",
			rawValue: fs + "CP" + etx,
			field:    "CP",
			want:     []string{""},
		},
		{
			name:     "non-empty field is unaffected",
			rawValue: fs + "CM" + "123 MAIN ST" + fs + "CN" + "MIAMI",
			field:    "CM",
			want:     []string{"123 MAIN ST"},
		},
		{
			name:     "field at end of raw data with no trailing separator",
			rawValue: fs + "CM" + "123 MAIN ST",
			field:    "CM",
			want:     []string{"123 MAIN ST"},
		},
		{
			name:     "repeating field with one empty occurrence",
			rawValue: fs + "MJ" + "GROUP1" + fs + "MJ" + fs + "MJ" + "GROUP3",
			field:    "MJ",
			want:     []string{"GROUP1", "", "GROUP3"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tran := NcpdpTransaction[RequestHeader]{RawValue: tc.rawValue}
			got := tran.ScanRawDataForFieldValues(tc.field)

			if len(got) != len(tc.want) {
				t.Fatalf("value count mismatch. Wanted: %d  Got: %d (%q)", len(tc.want), len(got), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("value %d mismatch. Wanted: %q  Got: %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}
