package claimserializer

import (
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// Production B1 response carrying ECL-coded 132-UH qualifiers (12, and 38 =
// Max Fair Price) alongside their 526-FQ values.
const RESPONSE_B1_CODED_UH = "D0B11A011003004904     20261007\x1e\x1cAM25\x1cC1000031Y0000\x1c2FBMC\x1cC23HZN17390110\x1e\x1cAM29\x1cCAKEITH\x1cCBARVIDSON\x1cC419560712\x1d\x1e\x1cAM21\x1cANP\x1cF3262802970378004997\x1c5F2\x1c6F035\x1c6F061\x1cUF02\x1cUH12\x1cFQMAXIMUM DAILY DOSE OF 1.000000\x1cUH38\x1cFQ405.66\x1e\x1cAM22\x1cEM1\x1cD27205908\x1e\x1cAM23\x1cF5432C\x1cF62061F\x1cF7100{\x1cF91729C\x1cFM46\x1cFC999999I\x1cFE999999I\x1c4U432C\x1c4V01\x1cMU1\x1cMV02\x1cMW2161F"

// Same response after appending a gateway note: the note is chunked into
// free-text qualifiers 01-03 (with UG continuations) inserted ahead of the
// coded 12/38 values, which must be retained verbatim. UF is re-emitted
// without zero padding (counts round-trip as integers).
const RESPONSE_B1_CODED_UH_APPENDED = "D0B11A011003004904     20261007\x1e\x1cAM25\x1cC1000031Y0000\x1c2FBMC\x1cC23HZN17390110\x1e\x1cAM29\x1cCAKEITH\x1cCBARVIDSON\x1cC419560712\x1d\x1e\x1cAM21\x1cANP\x1cF3262802970378004997\x1c5F2\x1c6F035\x1c6F061\x1cUF5\x1cUH01\x1cFQRSE: INGREDIENT COST UPDATED.  OLD INGRE\x1cUG+\x1cUH02\x1cFQDIENT COST: 717.66. NEW INGREDIENT COST:\x1cUG+\x1cUH03\x1cFQ 734.18\x1cUH12\x1cFQMAXIMUM DAILY DOSE OF 1.000000\x1cUH38\x1cFQ405.66\x1e\x1cAM22\x1cEM1\x1cD27205908\x1e\x1cAM23\x1cF5432C\x1cF62061F\x1cF7100{\x1cF91729C\x1cFM46\x1cFC999999I\x1cFE999999I\x1c4U432C\x1c4V01\x1cMU1\x1cMV02\x1cMW2161F"

func TestAppendMessageRetainsCodedQualifiersEndToEnd(t *testing.T) {
	obj, err := claimdeserializer.DeserializeResponse(RESPONSE_B1_CODED_UH)
	if err != nil {
		t.Fatalf("DeserializeResponse: %v", err)
	}

	resp, ok := obj.(response.Billing)
	if !ok {
		t.Fatalf("unexpected deserialized type %T", obj)
	}

	if len(resp.Claims) != 1 {
		t.Fatalf("claim count mismatch. Wanted: 1   Got: %d", len(resp.Claims))
	}

	status := &resp.Claims[0].Status

	if len(status.AdditionalMessages) != 2 {
		t.Fatalf("additional message count mismatch. Wanted: 2   Got: %d", len(status.AdditionalMessages))
	}
	for i, wantQfr := range []string{"12", "38"} {
		if got := status.AdditionalMessages[i].Qualifier; got == nil || *got != wantQfr {
			t.Errorf("message[%d] qualifier mismatch. Wanted: %q   Got: %v", i, wantQfr, got)
		}
	}

	err = status.AppendMessage("RSE: Ingredient cost updated.  Old ingredient cost: 717.66. New ingredient cost: 734.18")
	if err != nil {
		t.Fatalf("AppendMessage: %v", err)
	}

	got, err := Serialize(&resp)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}

	if got != RESPONSE_B1_CODED_UH_APPENDED {
		t.Errorf("serialized output mismatch.\nWanted: %q\nGot:    %q", RESPONSE_B1_CODED_UH_APPENDED, got)
	}
}
