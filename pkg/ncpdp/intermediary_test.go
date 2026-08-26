package ncpdp

import "testing"

const requestHeaderD0 = "880151D0B1          1011234567893     20260611          "
const requestHeaderD0Batch = "880151D0B1          3011234567893     20260611          "
const requestHeaderF6 = "F6B188015600PCN1234567101 1234567893    20260611VENDORCERT"

// D.0 claim segments (AM07). Each segment starts with SEGMENT(0x1E) +
// FIELD(0x1C) + segment ID, then FIELD(0x1C) + fieldId + value per field.
const claimSegmentPlain = "\x1e\x1cAM07\x1cEM1\x1cD26000001"
const claimSegmentBypass = "\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cEW01\x1cEXAUTH123"

func parseRequest(t *testing.T, raw string) *NcpdpTransaction[RequestHeader] {
	t.Helper()
	tran := NewTransactionRequest(raw)
	if err := tran.ParseNcpdp(); err != nil {
		t.Fatalf("ParseNcpdp: %v", err)
	}
	return &tran
}

type intermediaryAuthCase struct {
	name string
	raw  string
	id   string
	want bool
}

var intermediaryAuthD0Cases = []intermediaryAuthCase{
	{
		name: "D0 qualifier 01 with matching auth code",
		raw:  requestHeaderD0 + "\x1d" + claimSegmentBypass,
		id:   "AUTH123",
		want: true,
	},
	{
		name: "D0 numeric qualifier variant",
		raw:  requestHeaderD0 + "\x1d\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cEW1\x1cEXAUTH123",
		id:   "AUTH123",
		want: true,
	},
	{
		name: "D0 wrong qualifier",
		raw:  requestHeaderD0 + "\x1d\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cEW99\x1cEXAUTH123",
		id:   "AUTH123",
		want: false,
	},
	{
		name: "D0 auth code mismatch",
		raw:  requestHeaderD0 + "\x1d" + claimSegmentBypass,
		id:   "OTHER",
		want: false,
	},
	{
		name: "D0 missing qualifier",
		raw:  requestHeaderD0 + "\x1d\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cEXAUTH123",
		id:   "AUTH123",
		want: false,
	},
	{
		name: "D0 missing auth code",
		raw:  requestHeaderD0 + "\x1d\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cEW01",
		id:   "AUTH123",
		want: false,
	},
	{
		// Guards per-record pairing: record 2 carries a different qualified
		// auth code, only record 3 matches. A flat first-value scan fails this.
		name: "D0 batch matches on later record",
		raw: requestHeaderD0Batch +
			"\x1d" + claimSegmentPlain +
			"\x1d\x1e\x1cAM07\x1cEM1\x1cD26000002\x1cEW01\x1cEXOTHER" +
			"\x1d" + claimSegmentBypass,
		id:   "AUTH123",
		want: true,
	},
	{
		name: "D0 no claim segment",
		raw:  requestHeaderD0 + "\x1d\x1e\x1cAM01\x1cCAFIRST\x1cCBLAST",
		id:   "AUTH123",
		want: false,
	},
}

// F6 places AM19 at the transaction level (no group separator) and allows the
// 8H/8K/8M fields to repeat once per intermediary entry.
var intermediaryAuthF6Cases = []intermediaryAuthCase{
	{
		name: "F6 single matching entry",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H01\x1c8K06\x1c8MAUTH123",
		id:   "AUTH123",
		want: true,
	},
	{
		name: "F6 numeric code variants",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H1\x1c8K6\x1c8MAUTH123",
		id:   "AUTH123",
		want: true,
	},
	{
		name: "F6 second entry matches",
		raw: requestHeaderF6 + claimSegmentPlain +
			"\x1e\x1cAM19\x1c8G2\x1c8H01\x1c8KQQ\x1c8MOTHER\x1c8H01\x1c8K06\x1c8MAUTH123",
		id:   "AUTH123",
		want: true,
	},
	{
		// Qualifier 06 and the target ID both exist, but on different entries.
		// A flat scan that ignores entry boundaries would falsely match.
		name: "F6 fields must match on the same entry",
		raw: requestHeaderF6 + claimSegmentPlain +
			"\x1e\x1cAM19\x1c8G2\x1c8H01\x1c8K06\x1c8MOTHER\x1c8H02\x1c8KQQ\x1c8MAUTH123",
		id:   "AUTH123",
		want: false,
	},
	{
		name: "F6 wrong qualifier",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H01\x1c8KQQ\x1c8MAUTH123",
		id:   "AUTH123",
		want: false,
	},
	{
		name: "F6 wrong type code",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H02\x1c8K06\x1c8MAUTH123",
		id:   "AUTH123",
		want: false,
	},
	{
		name: "F6 no intermediary segment",
		raw:  requestHeaderF6 + claimSegmentPlain,
		id:   "AUTH123",
		want: false,
	},
	{
		name: "F6 id argument is trimmed",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H01\x1c8K06\x1c8MAUTH123",
		id:   "  AUTH123  ",
		want: true,
	},
	{
		name: "empty id never matches",
		raw:  requestHeaderF6 + claimSegmentPlain + "\x1e\x1cAM19\x1c8G1\x1c8H01\x1c8K06\x1c8MAUTH123",
		id:   "   ",
		want: false,
	},
}

func TestIsServiceBypassedAcrossVersions(t *testing.T) {
	cases := append(append([]intermediaryAuthCase{}, intermediaryAuthD0Cases...), intermediaryAuthF6Cases...)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tran := parseRequest(t, tc.raw)
			if got := tran.IsServiceBypassed(tc.id); got != tc.want {
				t.Errorf("IsServiceBypassed(%q) mismatch. Wanted: %v   Got: %v", tc.id, tc.want, got)
			}
		})
	}
}

// Defensive: the helper must not panic on nil or empty transactions.
func TestIsServiceBypassedDefensive(t *testing.T) {
	var nilTran *NcpdpTransaction[RequestHeader]
	if nilTran.IsServiceBypassed("AUTH123") {
		t.Error("nil transaction should not match")
	}

	empty := parseRequest(t, requestHeaderD0)
	if empty.IsServiceBypassed("AUTH123") {
		t.Error("transaction without segments should not match")
	}
}
