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

// A later payer with several qualifier/amount pairs must keep each pair
// intact and in order: the second HC within payer 2 starts a new nested
// element because HC is already occupied in the current pair, never because
// of its position in the segment.
func Test_RoundTripCobSecondPayerMultipleAmounts(t *testing.T) {
	fs := string(ncpdp.FIELD)

	oldAM05 := "AM05" + fs + "4C1" + fs + "5C01" + fs + "6C03" + fs + "7C015581    " + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000000{" + fs + "5E02" + fs + "6E70 " + fs + "6EA5 "
	newAM05 := "AM05" + fs + "4C2" +
		fs + "5C01" + fs + "6C03" + fs + "7C01111" + fs + "HB1" + fs + "HC07" + fs + "DV0000100{" +
		fs + "5C02" + fs + "6C03" + fs + "7C02222" + fs + "HB2" + fs + "HC07" + fs + "DV0000200{" + fs + "HC08" + fs + "DV0000300{"

	raw := strings.Replace(REQUEST_B1, oldAM05, newAM05, 1)
	if raw == REQUEST_B1 {
		t.Fatal("AM05 replacement did not apply to REQUEST_B1")
	}

	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	payers := obj.Claims[0].CoordinationOfBenefits.OtherPayer
	if len(payers) != 2 {
		t.Fatalf("OtherPayer count mismatch. Wanted: 2   Got: %v", len(payers))
	}
	if len(payers[0].OtherPayerAmountsPaid) != 1 || len(payers[1].OtherPayerAmountsPaid) != 2 {
		t.Fatalf("amount pair distribution mismatch. Wanted: 1/2   Got: %v/%v",
			len(payers[0].OtherPayerAmountsPaid), len(payers[1].OtherPayerAmountsPaid))
	}
	for i, want := range []struct {
		qualifier string
		amount    float64
	}{
		{"07", 20.00},
		{"08", 30.00},
	} {
		pair := payers[1].OtherPayerAmountsPaid[i]
		if pair.Qualifier == nil || *pair.Qualifier != want.qualifier || pair.AmountPaid == nil || *pair.AmountPaid != want.amount {
			t.Errorf("payer 2 pair %v mismatch. Wanted: {%v %v}   Got: {%v %v}", i+1, want.qualifier, want.amount, pair.Qualifier, pair.AmountPaid)
		}
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertFieldValues(t, serialized, "HB", []string{"1", "2"})
	assertFieldValues(t, serialized, "HC", []string{"07", "07", "08"})
	assertFieldValues(t, serialized, "DV", []string{"100{", "200{", "300{"})
}

// Empty fields inside the second occurrence of a repeating group must stay
// nil there (not bleed into the first occurrence or invent values), the
// first occurrence's values must survive, an empty string field round-trips
// in place, and re-serializing a reparsed message is stable.
func Test_RoundTripEmptyFieldsInSecondDurItem(t *testing.T) {
	fs := string(ncpdp.FIELD)
	rs := string(ncpdp.SEGMENT)

	raw := RESPONSE_B1 + rs + fs + "AM24" +
		fs + "J61" + fs + "E4PA" + fs + "FS" + fs + "FT0" + fs + "FU20260101" + fs + "FV0001000" + fs + "FW2" +
		fs + "J62" + fs + "E4TD" + fs + "FT1" + fs + "FU" + fs + "FV" + fs + "FW3"

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	items := obj.Claims[0].Dur.Items
	if len(items) != 2 {
		t.Fatalf("DUR item count mismatch. Wanted: 2   Got: %v", len(items))
	}
	if items[0].PreviousFillDate == nil || items[0].PreviousFillDate.Format("20060102") != "20260101" {
		t.Errorf("item 1 PreviousFillDate mismatch. Wanted: 20260101   Got: %v", items[0].PreviousFillDate)
	}
	if items[0].PreviousFillDateQuantity == nil || *items[0].PreviousFillDateQuantity != 1.0 {
		t.Errorf("item 1 PreviousFillDateQuantity mismatch. Wanted: 1   Got: %v", items[0].PreviousFillDateQuantity)
	}
	if items[0].ClinicalSignificanceCode == nil || *items[0].ClinicalSignificanceCode != "" {
		t.Errorf("item 1 ClinicalSignificanceCode mismatch. Wanted: empty string   Got: %v", items[0].ClinicalSignificanceCode)
	}
	if items[1].PreviousFillDate != nil {
		t.Errorf("item 2 PreviousFillDate mismatch. Wanted: nil   Got: %v", *items[1].PreviousFillDate)
	}
	if items[1].PreviousFillDateQuantity != nil {
		t.Errorf("item 2 PreviousFillDateQuantity mismatch. Wanted: nil   Got: %v", *items[1].PreviousFillDateQuantity)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	// Item 1's values inline; item 2's preserved empties re-emitted at the
	// segment's end. The empty string FS is a typed field and stays in place.
	assertFieldValues(t, serialized, "FU", []string{"20260101", ""})
	assertFieldValues(t, serialized, "FV", []string{"1000", ""})
	assertFieldValues(t, serialized, "FS", []string{""})
	assertFieldValues(t, serialized, "J6", []string{"1", "2"})

	// Reparsing the rebuilt message and serializing again must be stable.
	reparsed := response.Billing{}
	if err := claimdeserializer.DeserializeType(serialized, &reparsed); err != nil {
		t.Fatal(err)
	}
	serialized2, err := Serialize(&reparsed)
	if err != nil {
		t.Fatal(err)
	}
	if serialized2 != serialized {
		t.Errorf("re-serialization is not stable.\nFirst:  %q\nSecond: %q", serialized, serialized2)
	}
}
