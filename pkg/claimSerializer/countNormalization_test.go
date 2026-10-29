package claimserializer

import (
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// A received count that contradicts the actual occurrences is deliberately
// corrected on rebuild (seen in production: AM21 UF=4 over two additional
// messages, UF=11 over five). The message content itself round-trips in full;
// only the count is normalized to the true occurrence count.
func Test_RoundTripNormalizesMismatchedCounts(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "D0B11A011234567893     20231219" +
		gs + ss + fs + "AM21" + fs + "ANP" + fs + "F3AUTH123456" +
		fs + "UF4" +
		fs + "UH01" + fs + "FQMember redetermination reqd as of 2026" +
		fs + "UG+" +
		fs + "UH02" + fs + "FQ11." +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001"

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	status := obj.Claims[0].Status
	if len(status.AdditionalMessages) != 2 {
		t.Fatalf("Additional message count mismatch. Wanted: 2   Got: %v", len(status.AdditionalMessages))
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

		assertFieldValues(t, serialized, "UF", []string{"2"})
		assertFieldValues(t, serialized, "UH", []string{"01", "02"})
		assertFieldValues(t, serialized, "FQ", []string{"Member redetermination reqd as of 2026", "11."})
	}
}
