package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
)

// FuzzRoundTrip drives arbitrary bytes through the full
// NCPDP -> model -> JSON -> model -> NCPDP chain and holds the library to
// three invariants, independent of how malformed the input is:
//
//  1. never panic — refuse loudly or parse, but never crash;
//  2. self-consistency — any message the parser accepts must serialize, and
//     the serialized output must be accepted by our own parser (we must not
//     emit messages we would refuse);
//  3. convergence — repeated round trips of our own output reach a byte-stable
//     fixed point within a few rebuilds instead of drifting or growing
//     unboundedly. Counter-less repeating codes with empty values (a bare
//     "UG,UG,UG,UH" fragment) group ambiguously — canonical field order moves
//     the stray code one occurrence per pass until it settles — so the fixed
//     point may take up to one pass per occurrence; real traffic with
//     counters and values is stable after the first rebuild.
//
// Byte-equality with the INPUT is deliberately not an invariant: documented
// normalizations (implied-decimal/overpunch minimal forms, derived counts,
// separator normalization, non-UTF-8 replacement on the JSON leg) rewrite
// equivalent forms.
//
// Run with: go test -fuzz=FuzzRoundTrip -fuzztime=120s ./pkg/claimSerializer
func FuzzRoundTrip(f *testing.F) {
	fs := "\x1c"
	ss := "\x1e"
	gs := "\x1d"

	// Well-formed traffic of each routed shape.
	f.Add(REQUEST_B1)
	f.Add(REQUEST_B2)
	f.Add(REQUEST_B3)
	f.Add(REQUEST_E1)
	f.Add(RESPONSE_B2)
	f.Add(RESPONSE_B3)
	f.Add(RESPONSE_E1)

	// Malformed shapes from the ncpdpServices reports: empty segment, segment
	// id without its leading separator, duplicated segment, repeating-group
	// items with unparseable values, garbage field codes.
	f.Add("999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" +
		ss + fs + "AM04" + fs + "C2ID12345" +
		ss + ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "E103" + fs + "DJ1")
	f.Add("D0E11A073350321        20260820" +
		ss + fs + "AM25" + fs + "FODU332263" +
		ss + "AM29" + fs + "C419590930" + fs + "CAMARTIN" + fs + "CBLAW" +
		gs + ss + fs + "AM21" + fs + "ANA" + fs + "F300000000000000")
	f.Add("999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" +
		ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		ss + fs + "AM01" + fs + "C419990101" + fs + "CAJANE" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "DJ1")
	f.Add("D0B11A011234567893     20261001" +
		ss + fs + "AM20" + fs + "F4TEXT" +
		gs + ss + fs + "AM21" + fs + "ANP" + fs + "F3AUTH123" +
		ss + fs + "AM24" + fs + "J61" + fs + "E4DC" + fs + "FT0" + fs + "FU00000000" + fs + "FV" +
		fs + "J62" + fs + "E4DD" + fs + "FT1" + fs + "FU20261008" + fs + "FV0000060000")
	f.Add("880151D0B2TESTTEST1 1011234567893     20231219          " +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D2100001" + fs + "@A123" + fs + "\x04Z9" + fs + "C800" + fs + "C800")

	f.Fuzz(func(t *testing.T, raw string) {
		obj, err := claimdeserializer.Deserialize(raw)
		if err != nil {
			return // loud refusal is a valid outcome
		}

		out1, err := jsonRoundTripSerialize(obj)
		if err != nil {
			if strings.HasPrefix(err.Error(), "unknown type") {
				t.Skip("transaction type not covered by the JSON round-trip helper")
			}
			t.Fatalf("accepted input failed to serialize: %v\ninput: %q", err, raw)
		}

		assertRoundTripConvergence(t, raw, out1)
	})
}

// assertRoundTripConvergence re-parses and re-serializes the library's own
// output until it reaches a byte-stable fixed point, failing when a rebuild
// is refused by our own parser or the bytes still drift after nine rebuilds.
func assertRoundTripConvergence(t *testing.T, raw, first string) {
	t.Helper()

	prev := first
	for pass := 2; pass <= 10; pass++ {
		obj, err := claimdeserializer.Deserialize(prev)
		if err != nil {
			t.Fatalf("rebuilt message refused by our own parser (pass %d): %v\ninput: %q\nrebuilt: %q", pass, err, raw, prev)
		}

		next, err := jsonRoundTripSerialize(obj)
		if err != nil {
			t.Fatalf("round trip %d failed to serialize: %v\nrebuilt: %q", pass, err, prev)
		}

		if next == prev {
			return // reached the fixed point
		}
		prev = next
	}

	t.Fatalf("serialization did not converge within nine rebuilds:\ninput: %q\nlast: %q", raw, prev)
}
