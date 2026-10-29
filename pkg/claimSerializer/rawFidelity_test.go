package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// A duplicated non-repeating field (seen in production: one AM07 carrying C8
// and D3 twice each) must not collapse into a single occurrence. The first
// occurrence lands in the typed field; the duplicate is preserved as a
// dynamic field and re-emits, so the message round-trips without loss.
func Test_RoundTripPreservesDuplicateScalarFields(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "880151D0B2TESTTEST1 1011234567893     20231219          " +
		ss + fs + "AM04" + fs + "C2TEST" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D2100001" + fs + "E103" + fs + "D700002125401" +
		fs + "C800" + fs + "D300" + fs + "D300" + fs + "C800" + fs + "28EA"

	obj := request.Reversal{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	claim := obj.Claims[0].Claim
	if claim.OtherCoverageCode == nil || *claim.OtherCoverageCode != "00" {
		t.Errorf("Other coverage code mismatch. Wanted: 00   Got: %v", claim.OtherCoverageCode)
	}
	if claim.FillNumber == nil || *claim.FillNumber != "00" {
		t.Errorf("Fill number mismatch. Wanted: 00   Got: %v", claim.FillNumber)
	}

	for _, serialize := range []struct {
		name string
		run  func() (string, error)
	}{
		{"direct", func() (string, error) { return Serialize(&obj) }},
		{"via JSON", func() (string, error) { return jsonRoundTripSerialize(obj) }},
	} {
		serialized, err := serialize.run()
		if err != nil {
			t.Fatalf("%v serialization failed: %v", serialize.name, err)
		}

		assertFieldValues(t, serialized, "C8", []string{"00", "00"})
		assertFieldValues(t, serialized, "D3", []string{"00", "00"})

		// Reparsing the rebuilt message and serializing again must stay
		// stable: the duplicate keeps round-tripping, not multiplying.
		reparsed := request.Reversal{}
		if err := claimdeserializer.DeserializeType(serialized, &reparsed); err != nil {
			t.Fatalf("%v: failed to reparse: %v", serialize.name, err)
		}
		again, err := Serialize(&reparsed)
		if err != nil {
			t.Fatalf("%v: reserialization failed: %v", serialize.name, err)
		}
		assertFieldValues(t, again, "C8", []string{"00", "00"})
		assertFieldValues(t, again, "D3", []string{"00", "00"})
	}
}

// Header bytes re-emit exactly as received while the parsed fields are
// unchanged (seen in production: a PCN and a software vendor certification id
// right-justified with leading spaces). Editing a header field rebuilds the
// header normalized instead.
func Test_RoundTripPreservesHeaderPadding(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	origHeader := "880151D0B1     TQR  1011234567893     20231219     TQR  "
	raw := origHeader +
		ss + fs + "AM04" + fs + "C2TEST" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D2100001" + fs + "E103" + fs + "D700002125401"

	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	if obj.Header.Value.Pcn != "TQR" {
		t.Errorf("Parsed PCN mismatch. Wanted: TQR   Got: %q", obj.Header.Value.Pcn)
	}

	for _, serialize := range []struct {
		name string
		run  func() (string, error)
	}{
		{"direct", func() (string, error) { return Serialize(&obj) }},
		{"via JSON", func() (string, error) { return jsonRoundTripSerialize(obj) }},
	} {
		serialized, err := serialize.run()
		if err != nil {
			t.Fatalf("%v serialization failed: %v", serialize.name, err)
		}
		if !strings.HasPrefix(serialized, origHeader) {
			t.Errorf("%v: header bytes not preserved.\nWanted prefix: %q\nGot: %q", serialize.name, origHeader, serialized[:len(origHeader)])
		}
	}

	// A consumer edit invalidates the raw bytes: the header rebuilds
	// normalized (left-justified, space-padded fields).
	obj.Header.Value.Pcn = "NEWPCN"
	rebuiltHeader := "880151D0B1NEWPCN    1011234567893     20231219TQR       "
	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(serialized, rebuiltHeader) {
		t.Errorf("Edited header did not rebuild normalized.\nWanted prefix: %q\nGot: %q", rebuiltHeader, serialized[:len(rebuiltHeader)])
	}
}
