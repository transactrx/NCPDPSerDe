package claimserializer_test

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	claimserializer "github.com/transactrx/NCPDPSerDe/pkg/claimSerializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

const (
	fs  = "\x1c" // field separator
	gs  = "\x1d" // group separator
	seg = "\x1e" // segment separator
)

// A real D.0 discount-card response whose Insurance segment carries a single
// J7/J8 payer pair and NO payer-ID count (KR). KR is an F6-only field.
func d0InsuranceResponse() string {
	return "D0B11A011184739229     20260729" + seg +
		fs + "AM25" + fs + "C1RXPB21" + fs + "2FMAG17CN2" + fs + "J703" + fs + "J8019363" + gs + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582" + seg +
		fs + "AM22" + fs + "EM1" + fs + "D26430498"
}

// TestD0RoundTripOmitsF6PayerCount is the regression test for the bug where a
// no-op deserialize->serialize of a D.0 response injected an F6-only KR
// (Payer/Health Plan ID Count) field that the payer never sent.
func TestD0RoundTripOmitsF6PayerCount(t *testing.T) {
	rsp, err := claimdeserializer.DeserializeResponse(d0InsuranceResponse())
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if strings.Contains(out, fs+"KR") {
		t.Errorf("D0 output must not contain the F6-only KR field.\nout: %q", out)
	}
	// The payer identity itself must survive the round-trip.
	if !strings.Contains(out, fs+"J703") || !strings.Contains(out, fs+"J8019363") {
		t.Errorf("D0 output dropped the J7/J8 payer pair.\nout: %q", out)
	}
}

// TestF6RoundTripKeepsPayerCount proves the version gate is not a blanket
// suppression: an F6 transmission still auto-derives KR from the Payers slice.
func TestF6RoundTripKeepsPayerCount(t *testing.T) {
	// F6 response header leads with the version value "F6"; no group separator.
	raw := "F6B11A011184739229     20260729" + seg +
		fs + "AM25" + fs + "C1RXPB21" + fs + "2FMAG17CN2" + fs + "J703" + fs + "J8019363" + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	if len(billing.Insurance.Payers) == 0 {
		t.Fatalf("expected Payers slice to be populated for F6 input")
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if !strings.Contains(out, fs+"KR1") {
		t.Errorf("F6 output must auto-derive KR1 from the Payers slice.\nout: %q", out)
	}
}

// TestF6MultiPayerCountDerived guards the countfor behavior the user wants
// preserved: an F6 message with repeating J7/J8 pairs auto-derives the true
// count on serialization, regardless of any count supplied in the source.
func TestF6MultiPayerCountDerived(t *testing.T) {
	// F6 response, two repeating payer pairs, source count deliberately absent.
	raw := "F6B11A011184739229     20260729" + seg +
		fs + "AM25" + fs + "C1RXPB21" + fs + "2FMAG17CN2" +
		fs + "J703" + fs + "J8019363" + fs + "J703" + fs + "J8024368" + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	if got := len(billing.Insurance.Payers); got != 2 {
		t.Fatalf("expected 2 payers parsed, got %d", got)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"KR2") {
		t.Errorf("expected derived KR2 for two payers, got: %q", out)
	}
}

// TestD0ResponsePatientOmitsRR covers the second dual-mapped group: the
// response Patient segment (AM29). RR (618-RR Patient ID Count) is F6-only, so
// a D0 response carrying only patient names must not gain an RR field.
func TestD0ResponsePatientOmitsRR(t *testing.T) {
	// Shared AM29 patient segment precedes the claim group (AM21).
	raw := "D0B11A011184739229     20260729" + seg +
		fs + "AM29" + fs + "CAJESSICA" + fs + "CBSILLAMAN" + gs + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if strings.Contains(out, fs+"RR") {
		t.Errorf("D0 response must not contain the F6-only RR field.\nout: %q", out)
	}
	if !strings.Contains(out, fs+"CAJESSICA") || !strings.Contains(out, fs+"CBSILLAMAN") {
		t.Errorf("D0 response dropped patient name fields.\nout: %q", out)
	}
}

// TestD0OmitsHandPopulatedF6OnlyLeaf covers the general sinceVersion gate on
// leaf fields: an F6-only field set on the struct by application code must not
// leak into a D0 transmission, where the field does not exist.
func TestD0OmitsHandPopulatedF6OnlyLeaf(t *testing.T) {
	rsp, err := claimdeserializer.DeserializeResponse(d0InsuranceResponse())
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	recon := "RECON123"
	billing.Claims[0].Status.ReconciliationId = &recon

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if strings.Contains(out, fs+"34RECON123") {
		t.Errorf("D0 output must not contain the F6-only Reconciliation ID (34) field.\nout: %q", out)
	}
}

// TestF6KeepsHandPopulatedF6OnlyLeaf proves the same field is emitted for F6.
func TestF6KeepsHandPopulatedF6OnlyLeaf(t *testing.T) {
	raw := "F6B11A011184739229     20260729" + seg +
		fs + "AM25" + fs + "C1RXPB21" + fs + "2FMAG17CN2" + fs + "J703" + fs + "J8019363" + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	recon := "RECON123"
	billing.Claims[0].Status.ReconciliationId = &recon

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if !strings.Contains(out, fs+"34RECON123") {
		t.Errorf("F6 output must contain the Reconciliation ID (34) field.\nout: %q", out)
	}
}

// TestF6ResponsePatientDerivesRR proves RR is still auto-derived for F6, where
// the patient-ID group (CX/CY) genuinely repeats.
func TestF6ResponsePatientDerivesRR(t *testing.T) {
	// Shared AM29 patient segment precedes the claim group (AM21); no group
	// separator in F6.
	raw := "F6B11A011184739229     20260729" + seg +
		fs + "AM29" + fs + "CX01" + fs + "CYAAA111" + fs + "CX02" + fs + "CYBBB222" + seg +
		fs + "AM21" + fs + "ANP" + fs + "F3655380953234993582"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	if got := len(billing.Patient.Ids); got != 2 {
		t.Fatalf("expected 2 patient IDs parsed, got %d", got)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"RR2") {
		t.Errorf("F6 response must auto-derive RR2 from the Ids slice.\nout: %q", out)
	}
}
