package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// fieldValuesOf returns the raw values of every occurrence of a field code in
// a serialized claim, in transmission order.
func fieldValuesOf(serialized, code string) []string {
	values := []string{}

	rest := serialized
	for {
		i := strings.Index(rest, string(ncpdp.FIELD)+code)
		if i < 0 {
			return values
		}
		rest = rest[i+1+len(code):]

		value := rest
		if j := strings.IndexAny(rest, string([]byte{ncpdp.FIELD, ncpdp.SEGMENT, ncpdp.GROUP})); j >= 0 {
			value = rest[:j]
		}
		values = append(values, value)
	}
}

func assertFieldValues(t *testing.T, serialized, code string, want []string) {
	t.Helper()

	got := fieldValuesOf(serialized, code)
	if len(got) != len(want) {
		t.Errorf("Field %v occurrence mismatch. Wanted: %q   Got: %q", code, want, got)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Field %v #%v mismatch. Wanted: %q   Got: %q", code, i+1, want[i], got[i])
		}
	}
}

// An AM05 with two other payers must rebuild with each payer's own HB/NR
// counts and its qualifier/amount pairs intact. The per-segment occurrence
// counts used to split payer 2's pairs and double its rebuilt counts.
func Test_RoundTripCobTwoPayers(t *testing.T) {
	fs := string(ncpdp.FIELD)

	oldAM05 := "AM05" + fs + "4C1" + fs + "5C01" + fs + "6C03" + fs + "7C015581    " + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000000{" + fs + "5E02" + fs + "6E70 " + fs + "6EA5 "
	newAM05 := "AM05" + fs + "4C2" +
		fs + "5C01" + fs + "6C03" + fs + "7C01111" + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000100{" + fs + "NR01" + fs + "NP06" + fs + "NQ0000050{" +
		fs + "5C02" + fs + "6C03" + fs + "7C02222" + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000200{" + fs + "NR01" + fs + "NP06" + fs + "NQ0000060{"

	raw := strings.Replace(REQUEST_B1, oldAM05, newAM05, 1)
	if raw == REQUEST_B1 {
		t.Fatal("AM05 replacement did not apply to REQUEST_B1")
	}

	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
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

		assertFieldValues(t, serialized, "HB", []string{"1", "1"})
		assertFieldValues(t, serialized, "HC", []string{"07", "07"})
		assertFieldValues(t, serialized, "DV", []string{"100{", "200{"})
		assertFieldValues(t, serialized, "NR", []string{"1", "1"})
		assertFieldValues(t, serialized, "NP", []string{"06", "06"})
		assertFieldValues(t, serialized, "NQ", []string{"50{", "60{"})
	}
}

// Empty FU (date) and FV (quantity) in a DUR item must stay nil in the model
// and round-trip as present-but-empty fields: no invented 0001-01-01 or 0000,
// and no dropped field.
func Test_RoundTripPreservesEmptyDurFields(t *testing.T) {
	fs := string(ncpdp.FIELD)
	rs := string(ncpdp.SEGMENT)

	raw := RESPONSE_B1 + rs + fs + "AM24" + fs + "J61" + fs + "E4PA" + fs + "FT0" + fs + "FU" + fs + "FV" + fs + "FW2"

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	items := obj.Claims[0].Dur.Items
	if len(items) != 1 {
		t.Fatalf("DUR item count mismatch. Wanted: 1   Got: %v", len(items))
	}
	if items[0].PreviousFillDate != nil {
		t.Errorf("PreviousFillDate mismatch. Wanted: nil   Got: %v", *items[0].PreviousFillDate)
	}
	if items[0].PreviousFillDateQuantity != nil {
		t.Errorf("PreviousFillDateQuantity mismatch. Wanted: nil   Got: %v", *items[0].PreviousFillDateQuantity)
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

		assertFieldValues(t, serialized, "FU", []string{""})
		assertFieldValues(t, serialized, "FV", []string{""})
	}
}

// An empty C4 (date of birth) must stay nil in the model and round-trip as a
// present-but-empty field instead of becoming 0001-01-01 or disappearing.
func Test_RoundTripPreservesEmptyBirthDate(t *testing.T) {
	raw := strings.Replace(REQUEST_E1, "C419420501", "C4", 1)
	if raw == REQUEST_E1 {
		t.Fatal("C4 replacement did not apply to REQUEST_E1")
	}

	obj := request.Eligibility{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	if obj.Patient.BirthDate != nil {
		t.Errorf("BirthDate mismatch. Wanted: nil   Got: %v", *obj.Patient.BirthDate)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertFieldValues(t, serialized, "C4", []string{""})
}

// Blank-filled numeric fields (how DX data-collection feeds mark "not
// populated") must round-trip byte for byte instead of becoming zeros.
func Test_RoundTripPreservesBlankNumericFields(t *testing.T) {
	raw := strings.Replace(REQUEST_B2, "#E00028", "#E     ", 1)
	raw = strings.Replace(raw, "&U1481", "&U    ", 1)
	if strings.Count(raw, "#E     ") != 1 || strings.Count(raw, "&U    ") != 1 {
		t.Fatal("blank-fill replacements did not apply to REQUEST_B2")
	}

	obj := request.Reversal{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	dataCollection := obj.Claims[0].DataCollection
	if dataCollection.DecimalQuantity != nil {
		t.Errorf("DecimalQuantity mismatch. Wanted: nil   Got: %v", *dataCollection.DecimalQuantity)
	}
	if dataCollection.AcquisitionCost != nil {
		t.Errorf("AcquisitionCost mismatch. Wanted: nil   Got: %v", *dataCollection.AcquisitionCost)
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

		assertFieldValues(t, serialized, "#E", []string{"     "})
		assertFieldValues(t, serialized, "&U", []string{"    "})
	}
}
