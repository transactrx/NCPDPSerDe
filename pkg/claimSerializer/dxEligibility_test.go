package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// buildDxEligibilityRequest builds a non-standard DX data-collection E1
// request: a D.0-era eligibility transaction followed by a claim group
// carrying AM11 pricing and the proprietary AMXX data-collection segment.
func buildDxEligibilityRequest() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	return "880151DXE122222222221011730433129     20261020          " +
		ss + fs + "AM04" + fs + "C2CARD4M94KP02" + fs + "CCJOHN" + fs + "CDDOE" + fs + "C61" +
		ss + fs + "AM01" + fs + "CX01" + fs + "CY123456789" + fs + "C419440801" + fs + "C51" + fs + "CAJOHN" + fs + "CBDOE" +
		gs +
		ss + fs + "AM11" + fs + "D90000000{" + fs + "DQ0000000{" + fs + "DU0000000{" +
		ss + fs + "AMXX" + fs + "&GG" + fs + "&N1558764159" + fs + "&U        " + fs + "#E     "
}

// Per the NCPDP standard an eligibility request carries no claim section in
// D.0 or F6, so the model defines no claim group. Non-standard DX
// data-collection E1 feeds append AM11/AMXX anyway; the platform
// intentionally strips those during claim compression before forwarding, and
// the library matches that flow: the standard portion parses, the group is
// dropped, and the rebuilt message carries no group separator.
func Test_EligibilityRequestDropsNonStandardClaimGroup(t *testing.T) {
	raw := buildDxEligibilityRequest()

	obj := request.Eligibility{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	assertString(t, "Patient last name (CB)", obj.Patient.LastName, "DOE")
	assertString(t, "Insurance cardholder id (C2)", obj.Insurance.Cardholder.Id, "CARD4M94KP02")

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertNoGroupSeparators(t, serialized)
	for _, code := range []string{"AM11", "AMXX"} {
		if strings.Contains(serialized, code) {
			t.Errorf("rebuilt eligibility request still carries dropped segment %v", code)
		}
	}
}

// A standard E1 request (no group separator) must keep parsing and
// serializing without a group separator.
func Test_StandardEligibilityHasNoClaimGroup(t *testing.T) {
	obj := request.Eligibility{}
	if err := claimdeserializer.DeserializeType(REQUEST_E1, &obj); err != nil {
		t.Fatal(err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertNoGroupSeparators(t, serialized)
}
