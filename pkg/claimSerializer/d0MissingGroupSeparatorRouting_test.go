package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// The group separator is what delimits claims in a pre-vEB (D0) transmission,
// but some senders omit it. The claim group is still locatable by segment ID
// (the same routing vEB+ uses by design), so the claim must parse into the
// typed Claims slice — consumers read those fields — instead of stranding
// every claim segment in transmission-level dynamic preservation.

func buildD0BillingWithoutGroupSeparator() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	// 56-byte D0 request header, then shared and claim segments with no 0x1D.
	return "123456" + "D0" + "B1" + "TESTPCN   " + "1" + "01" + "1234567893     " + "20261103" + "CERT123456" +
		ss + fs + "AM04" + fs + "C2CARD001" +
		ss + fs + "AM01" + fs + "C419800115" + fs + "CAJOHN" +
		ss + fs + "AM07" + fs + "EM1" + fs + "D2123456789012" +
		ss + fs + "AM11" + fs + "D9150{" + fs + "DQ150{"
}

func Test_D0MissingGroupSeparatorStillPopulatesClaim(t *testing.T) {
	raw := buildD0BillingWithoutGroupSeparator()

	obj, err := claimdeserializer.Deserialize(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	req, ok := obj.(request.Billing)
	if !ok {
		t.Fatalf("unexpected type %T", obj)
	}

	assertStrPtr(t, "missing group separator", "Insurance.Cardholder.Id", strPtr("CARD001"), req.Insurance.Cardholder.Id)
	assertStrPtr(t, "missing group separator", "Patient.FirstName", strPtr("JOHN"), req.Patient.FirstName)

	if len(req.Claims) != 1 {
		t.Fatalf("claim count mismatch. Wanted: 1   Got: %v", len(req.Claims))
	}
	claim := req.Claims[0]

	assertStrPtr(t, "missing group separator", "PrescriptionServiceReference.Number", strPtr("123456789012"), claim.Claim.PrescriptionServiceReference.Number)
	if claim.Pricing.IngredientCostSubmitted == nil || *claim.Pricing.IngredientCostSubmitted != 15.00 {
		t.Errorf("Pricing.IngredientCostSubmitted mismatch. Wanted: 15.00   Got: %v", claim.Pricing.IngredientCostSubmitted)
	}

	if len(req.DynamicSegments) != 0 {
		t.Errorf("claim segments were preserved dynamically on the transmission instead of mapping: %v entries", len(req.DynamicSegments))
	}
	if len(claim.DynamicSegments) != 0 {
		t.Errorf("claim segments were preserved dynamically on the claim instead of mapping: %v entries", len(claim.DynamicSegments))
	}
}

// The rebuild emits the group separator at its canonical position, and that
// repair is a pure canonicalization: the rebuild parses back to exactly the
// same data and is byte-stable from then on.
func Test_D0MissingGroupSeparatorRebuildConverges(t *testing.T) {
	raw := buildD0BillingWithoutGroupSeparator()

	obj, err := claimdeserializer.Deserialize(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	req := obj.(request.Billing)

	out, err := Serialize(&req)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	gsIndex := strings.Index(out, string(ncpdp.GROUP))
	if gsIndex < 0 {
		t.Fatalf("rebuild did not restore the group separator: %q", out)
	}
	if am07 := strings.Index(out, "AM07"); am07 < gsIndex {
		t.Errorf("group separator not emitted before the claim segments: %q", out)
	}

	reparsed, err := claimdeserializer.Deserialize(out)
	if err != nil {
		t.Fatalf("rebuild does not reparse: %v", err)
	}
	assertSameJson(t, obj, reparsed)

	r2 := reparsed.(request.Billing)
	again, err := Serialize(&r2)
	if err != nil {
		t.Fatalf("second serialize: %v", err)
	}
	if again != out {
		t.Errorf("rebuild is not byte-stable:\nfirst:  %q\nsecond: %q", out, again)
	}
}
