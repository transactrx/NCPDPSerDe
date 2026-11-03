package claimserializer

import (
	"encoding/json"
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// F6 eliminated the group separator, so segment placement cannot depend on
// position: payers send transmission-level segments (AM25 insurance, AM29
// patient) after the claim segments, and those must still land in the typed
// shared fields instead of being stranded in the claim's dynamic preservation.

// Mirrors the reported production shape: status and claim first, insurance
// and patient at the end of the transmission.
func buildF6ResponseWithTrailingSharedSegments() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	return "F6B11A011234567893     20261026" +
		ss + fs + "AM21" + fs + "ANP" + fs + "F3AUTH99887" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001" +
		ss + fs + "AM25" + fs + "2FNETRX01" +
		ss + fs + "AM29" + fs + "CAJANE" + fs + "CBTESTER" + fs + "C419800115"
}

func Test_F6SharedSegmentsAfterClaimSegmentsMapToSharedFields(t *testing.T) {
	raw := buildF6ResponseWithTrailingSharedSegments()

	obj, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	resp, ok := obj.(response.Billing)
	if !ok {
		t.Fatalf("unexpected type %T", obj)
	}

	assertStrPtr(t, "trailing shared segments", "Insurance.NetworkReimbursementId", strPtr("NETRX01"), resp.Insurance.NetworkReimbursementId)
	assertStrPtr(t, "trailing shared segments", "Patient.FirstName", strPtr("JANE"), resp.Patient.FirstName)
	assertStrPtr(t, "trailing shared segments", "Patient.LastName", strPtr("TESTER"), resp.Patient.LastName)
	if resp.Patient.BirthDate == nil || resp.Patient.BirthDate.Format("20060102") != "19800115" {
		t.Errorf("Patient.BirthDate mismatch. Wanted: 19800115   Got: %v", resp.Patient.BirthDate)
	}

	if len(resp.Claims) != 1 {
		t.Fatalf("claim count mismatch. Wanted: 1   Got: %v", len(resp.Claims))
	}
	if len(resp.Claims[0].DynamicSegments) != 0 {
		t.Errorf("shared segments were preserved dynamically on the claim instead of mapping: %v entries", len(resp.Claims[0].DynamicSegments))
	}
	if len(resp.DynamicSegments) != 0 {
		t.Errorf("shared segments were preserved dynamically on the transmission instead of mapping: %v entries", len(resp.DynamicSegments))
	}
}

// The rebuild emits the shared segments at their catalog position (before the
// claim), and that reordering is a pure canonicalization: the rebuild parses
// back to exactly the same data and is byte-stable from then on.
func Test_F6TrailingSharedSegmentRebuildConverges(t *testing.T) {
	raw := buildF6ResponseWithTrailingSharedSegments()

	obj, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	resp := obj.(response.Billing)

	out, err := Serialize(&resp)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	if strings.Index(out, "AM25") > strings.Index(out, "AM21") {
		t.Errorf("rebuild did not move AM25 to its catalog position before the claim: %q", out)
	}

	reparsed, err := claimdeserializer.DeserializeResponse(out)
	if err != nil {
		t.Fatalf("rebuild does not reparse: %v", err)
	}
	assertSameJson(t, obj, reparsed)

	again, err := func() (string, error) {
		r := reparsed.(response.Billing)
		return Serialize(&r)
	}()
	if err != nil {
		t.Fatalf("second serialize: %v", err)
	}
	if again != out {
		t.Errorf("rebuild is not byte-stable:\nfirst:  %q\nsecond: %q", out, again)
	}
}

// Unknown segment IDs have no catalog home, so they keep positional
// attribution: before any claim segment they belong to the transmission,
// after one they belong to the claim.
func Test_F6UnknownSegmentsKeepPositionalAttribution(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	raw := "F6B11A011234567893     20261026" +
		ss + fs + "AM98" + fs + "XX1" +
		ss + fs + "AM21" + fs + "ANP" + fs + "F3AUTH99887" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001" +
		ss + fs + "AM99" + fs + "ZZ2" +
		ss + fs + "AM25" + fs + "2FNETRX01"

	obj, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	resp := obj.(response.Billing)

	if len(resp.DynamicSegments) != 1 {
		t.Errorf("transmission dynamic segment count mismatch. Wanted: 1 (AM98)   Got: %v", len(resp.DynamicSegments))
	}
	if len(resp.Claims) != 1 || len(resp.Claims[0].DynamicSegments) != 1 {
		t.Fatalf("claim dynamic segment count mismatch. Wanted: 1 (AM99)   Got: %+v", resp.Claims)
	}
	assertStrPtr(t, "unknown segment attribution", "Insurance.NetworkReimbursementId", strPtr("NETRX01"), resp.Insurance.NetworkReimbursementId)
}

func assertSameJson(t *testing.T, a, b interface{}) {
	t.Helper()

	ja, errA := json.Marshal(a)
	jb, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		t.Fatalf("marshal: %v / %v", errA, errB)
	}
	if string(ja) != string(jb) {
		t.Errorf("parses differ:\noriginal: %s\nrebuild:  %s", ja, jb)
	}
}
