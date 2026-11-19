package claimserializer_test

import (
	"encoding/json"
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	claimserializer "github.com/transactrx/NCPDPSerDe/pkg/claimSerializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// 448-ED Compound Ingredient Quantity is the only field whose implied decimal
// count changed between D0 and F6: 9(7)v999 (3 decimals) became 9(7)v9999999
// (7 decimals). 483-HE/560-AY Percentage Tax Rate kept their 4 decimals but
// lost their sign in F6: s9(3)v9999 became 9(3)v9999, so F6 output must not
// carry an overpunch character. These tests pin both behaviors per version.

func d0CompoundBillingRequest(ed, he string) string {
	return "880151D0B1          1011234567893     20231219          " +
		seg + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		gs + seg + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" +
		seg + fs + "AM11" + fs + "D9105E" + fs + "HE" + he +
		seg + fs + "AM10" + fs + "EF03" + fs + "EG2" + fs + "EC1" +
		fs + "RE03" + fs + "TE51927177800" + fs + "ED" + ed + fs + "EE0000500{" + fs + "UE01"
}

func f6CompoundBillingRequest(ed, he string) string {
	return "F6B100880151TEST      1011234567893     20260611SVCID     " +
		seg + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		seg + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" +
		seg + fs + "AM11" + fs + "D9105E" + fs + "HE" + he +
		seg + fs + "AM10" + fs + "EF03" + fs + "EG2" + fs + "EC1" +
		fs + "RE03" + fs + "TE51927177800" + fs + "ED" + ed + fs + "EE0000500{" + fs + "UE01"
}

func deserializeBillingRequest(t *testing.T, raw string) request.Billing {
	t.Helper()

	req, err := claimdeserializer.Deserialize(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}

	billing, ok := req.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing, got: %T", req)
	}
	if len(billing.Claims) == 0 {
		t.Fatal("expected a claim group to be populated")
	}

	return billing
}

// TestD0CompoundQuantityParsesWith3Decimals pins the D0 interpretation of
// 448-ED: 0001500 with 3 implied decimals is 1.5, and a D0 round trip emits
// the 3-decimal shape again.
func TestD0CompoundQuantityParsesWith3Decimals(t *testing.T) {
	billing := deserializeBillingRequest(t, d0CompoundBillingRequest("0001500", "0070{"))

	ingredients := billing.Claims[0].Compound.Ingredients
	if len(ingredients) != 1 {
		t.Fatalf("expected 1 compound ingredient, got %d", len(ingredients))
	}
	if q := ingredients[0].Quantity; q == nil || *q != 1.5 {
		t.Fatalf("D0 ED=0001500 must parse as 1.5, got %v", q)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"ED1500"+fs) {
		t.Errorf("D0 output must carry ED with 3 implied decimals (1500).\nout: %q", out)
	}
}

// TestF6CompoundQuantityParsesWith7Decimals pins the F6 widening of 448-ED:
// 15000000 with 7 implied decimals is 1.5 — NOT 15000.000 as the D0 divisor
// would read it — and an F6 round trip emits the 7-decimal shape.
func TestF6CompoundQuantityParsesWith7Decimals(t *testing.T) {
	billing := deserializeBillingRequest(t, f6CompoundBillingRequest("15000000", "00700"))

	ingredients := billing.Claims[0].Compound.Ingredients
	if len(ingredients) != 1 {
		t.Fatalf("expected 1 compound ingredient, got %d", len(ingredients))
	}
	if q := ingredients[0].Quantity; q == nil || *q != 1.5 {
		t.Fatalf("F6 ED=15000000 must parse as 1.5, got %v", q)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"ED15000000"+fs) {
		t.Errorf("F6 output must carry ED with 7 implied decimals (15000000).\nout: %q", out)
	}
}

// TestF6CompoundQuantityFractionalPrecision proves F6 precision beyond three
// decimals survives a round trip: 0.0000416 is representable only with the
// 7-decimal divisor.
func TestF6CompoundQuantityFractionalPrecision(t *testing.T) {
	billing := deserializeBillingRequest(t, f6CompoundBillingRequest("00000416", "00700"))

	if q := billing.Claims[0].Compound.Ingredients[0].Quantity; q == nil || *q != 0.0000416 {
		t.Fatalf("F6 ED=00000416 must parse as 0.0000416, got %v", q)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"ED00000416"+fs) {
		t.Errorf("F6 output must preserve the 7-decimal fraction (00000416).\nout: %q", out)
	}
}

// TestD0TaxRateKeepsOverpunch pins the D0 behavior of 483-HE: the field is
// signed (s9(3)v9999), so serialization keeps the overpunch character.
func TestD0TaxRateKeepsOverpunch(t *testing.T) {
	billing := deserializeBillingRequest(t, d0CompoundBillingRequest("0001500", "0070{"))

	if r := billing.Claims[0].Pricing.PercentageSalesTaxRateSubmitted; r == nil || *r != 0.07 {
		t.Fatalf("D0 HE=0070{ must parse as 0.07, got %v", r)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"HE0070{") {
		t.Errorf("D0 output must keep the overpunched HE value (0070{).\nout: %q", out)
	}
}

// TestF6TaxRateIsUnsigned pins the F6 change to 483-HE: the field became
// unsigned (9(3)v9999), so serialization must emit plain digits and never an
// overpunch character.
func TestF6TaxRateIsUnsigned(t *testing.T) {
	billing := deserializeBillingRequest(t, f6CompoundBillingRequest("15000000", "00700"))

	if r := billing.Claims[0].Pricing.PercentageSalesTaxRateSubmitted; r == nil || *r != 0.07 {
		t.Fatalf("F6 HE=00700 must parse as 0.07, got %v", r)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"HE00700") {
		t.Errorf("F6 output must carry HE as plain unsigned digits (00700).\nout: %q", out)
	}
	if strings.Contains(out, fs+"HE0070{") {
		t.Errorf("F6 output must not overpunch the unsigned HE field.\nout: %q", out)
	}
}

// TestF6ResponseTaxRateIsUnsigned covers the response-side twin 560-AY, and
// proves lenient parsing: an overpunched value from a nonconforming F6 sender
// still parses, but re-serializes unsigned.
func TestF6ResponseTaxRateIsUnsigned(t *testing.T) {
	raw := "F6B11A011234567893     20260611" +
		seg + fs + "AM21" + fs + "ANP" +
		seg + fs + "AM23" + fs + "F5107B" + fs + "AY0070{"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	if r := billing.Claims[0].Pricing.PercentageSalesTaxRatePaid; r == nil || *r != 0.07 {
		t.Fatalf("F6 AY=0070{ must still parse leniently as 0.07, got %v", r)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"AY00700") {
		t.Errorf("F6 output must carry AY as plain unsigned digits (00700).\nout: %q", out)
	}
	if strings.Contains(out, fs+"AY0070{") {
		t.Errorf("F6 output must not overpunch the unsigned AY field.\nout: %q", out)
	}
}

// TestF6CompoundQuantityBoundaryPrecision pins float64 fidelity at the edges
// of the widened F6 format: the 14-digit maximum (9999999.9999999) and the
// smallest representable fraction (0.0000001) must both survive a parse →
// serialize round trip digit for digit.
func TestF6CompoundQuantityBoundaryPrecision(t *testing.T) {
	tests := []struct {
		name  string
		wire  string
		value float64
	}{
		{"maximum 9(7)v9999999", "99999999999999", 9999999.9999999},
		{"smallest fraction", "00000001", 0.0000001},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			billing := deserializeBillingRequest(t, f6CompoundBillingRequest(test.wire, "00700"))

			if q := billing.Claims[0].Compound.Ingredients[0].Quantity; q == nil || *q != test.value {
				t.Fatalf("F6 ED=%s must parse as %v, got %v", test.wire, test.value, q)
			}

			out, err := claimserializer.Serialize(&billing)
			if err != nil {
				t.Fatalf("serialize: %v", err)
			}
			if !strings.Contains(out, fs+"ED"+test.wire+fs) {
				t.Errorf("F6 output must re-emit ED%s digit for digit.\nout: %q", test.wire, out)
			}
		})
	}
}

// TestD0TaxRateNegativeOverpunchRoundTrip proves the unsignedSince gate did
// not disturb D0's negative handling: an overpunched negative rate still
// parses negative and re-serializes byte-identically.
func TestD0TaxRateNegativeOverpunchRoundTrip(t *testing.T) {
	billing := deserializeBillingRequest(t, d0CompoundBillingRequest("0001500", "0070N"))

	if r := billing.Claims[0].Pricing.PercentageSalesTaxRateSubmitted; r == nil || *r != -0.0705 {
		t.Fatalf("D0 HE=0070N must parse as -0.0705, got %v", r)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"HE0070N") {
		t.Errorf("D0 output must keep the negative overpunched HE value (0070N).\nout: %q", out)
	}
}

// TestVersionedAttributesDoNotLeakAcrossVersions parses the SAME wire digits
// alternately as D0 and F6: 0001500 is 1.5 under D0's 3-decimal divisor and
// 0.00015 under F6's 7-decimal one. Interleaving guards any future caching of
// resolved field settings from bleeding one version's divisor into the other.
func TestVersionedAttributesDoNotLeakAcrossVersions(t *testing.T) {
	for i := 0; i < 3; i++ {
		d0 := deserializeBillingRequest(t, d0CompoundBillingRequest("0001500", "0070{"))
		if q := d0.Claims[0].Compound.Ingredients[0].Quantity; q == nil || *q != 1.5 {
			t.Fatalf("pass %d: D0 ED=0001500 must parse as 1.5, got %v", i, q)
		}

		f6 := deserializeBillingRequest(t, f6CompoundBillingRequest("0001500", "00700"))
		if q := f6.Claims[0].Compound.Ingredients[0].Quantity; q == nil || *q != 0.00015 {
			t.Fatalf("pass %d: F6 ED=0001500 must parse as 0.00015, got %v", i, q)
		}

		d0Out, err := claimserializer.Serialize(&d0)
		if err != nil {
			t.Fatalf("pass %d: serialize D0: %v", i, err)
		}
		f6Out, err := claimserializer.Serialize(&f6)
		if err != nil {
			t.Fatalf("pass %d: serialize F6: %v", i, err)
		}
		if !strings.Contains(d0Out, fs+"ED1500"+fs) {
			t.Errorf("pass %d: D0 output must use the 3-decimal shape.\nout: %q", i, d0Out)
		}
		if !strings.Contains(f6Out, fs+"ED00001500"+fs) {
			t.Errorf("pass %d: F6 output must use the 7-decimal shape.\nout: %q", i, f6Out)
		}
	}
}

// TestF6VersionedFieldsSurviveJsonRoundTrip covers the consumer path that
// moves claims through JSON (ncpdpServices): the 7-decimal quantity and the
// unsigned tax rate must survive marshal → unmarshal → serialize and produce
// the same bytes as a direct serialization.
func TestF6VersionedFieldsSurviveJsonRoundTrip(t *testing.T) {
	billing := deserializeBillingRequest(t, f6CompoundBillingRequest("15000000", "00700"))

	direct, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}

	jsonBytes, err := json.Marshal(billing)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	fresh := request.Billing{}
	if err := json.Unmarshal(jsonBytes, &fresh); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if q := fresh.Claims[0].Compound.Ingredients[0].Quantity; q == nil || *q != 1.5 {
		t.Fatalf("quantity must survive the JSON leg as 1.5, got %v", q)
	}

	rebuilt, err := claimserializer.Serialize(&fresh)
	if err != nil {
		t.Fatalf("serialize after JSON round trip: %v", err)
	}
	if rebuilt != direct {
		t.Errorf("JSON round trip changed the serialized bytes.\ndirect:  %q\nrebuilt: %q", direct, rebuilt)
	}
	if !strings.Contains(rebuilt, fs+"ED15000000"+fs) || !strings.Contains(rebuilt, fs+"HE00700") {
		t.Errorf("rebuilt F6 message lost a versioned field shape.\nout: %q", rebuilt)
	}
}

// TestF6QuestionPercentResponseIsUnsigned covers 379-4D, which the F6
// dictionary defines as unsigned 9(3)v99: F6 output must carry plain digits.
func TestF6QuestionPercentResponseIsUnsigned(t *testing.T) {
	raw := "F6B100880151TEST      1011234567893     20260611SVCID     " +
		seg + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		seg + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" +
		seg + fs + "AM14" + fs + "2Q01" + fs + "2Z1" + fs + "4B001" + fs + "4D05000"

	billing := deserializeBillingRequest(t, raw)

	questions := billing.Claims[0].AdditionalDocumentation.Questions
	if len(questions) != 1 || questions[0].PercentResponse == nil || *questions[0].PercentResponse != 50.0 {
		t.Fatalf("F6 4D=05000 must parse as 50.0, got %+v", questions)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"4D5000") {
		t.Errorf("F6 output must carry 4D as plain unsigned digits (5000).\nout: %q", out)
	}
	if strings.Contains(out, fs+"4D500{") {
		t.Errorf("F6 output must not overpunch the unsigned 4D field.\nout: %q", out)
	}
}

// TestD0QuestionPercentResponseKeepsOverpunch pins 379-4D's historical D0
// serialization (overpunched) unchanged.
func TestD0QuestionPercentResponseKeepsOverpunch(t *testing.T) {
	raw := "880151D0B1          1011234567893     20231219          " +
		seg + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		gs + seg + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" +
		seg + fs + "AM14" + fs + "2Q01" + fs + "2Z1" + fs + "4B001" + fs + "4D0500{"

	billing := deserializeBillingRequest(t, raw)

	questions := billing.Claims[0].AdditionalDocumentation.Questions
	if len(questions) != 1 || questions[0].PercentResponse == nil || *questions[0].PercentResponse != 50.0 {
		t.Fatalf("D0 4D=0500{ must parse as 50.0, got %+v", questions)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"4D500{") {
		t.Errorf("D0 output must keep the overpunched 4D value (500{).\nout: %q", out)
	}
}

// TestD0ResponseTaxRateKeepsOverpunch pins 560-AY's D0 behavior unchanged.
func TestD0ResponseTaxRateKeepsOverpunch(t *testing.T) {
	raw := "D0B11A011234567893     20260611" +
		seg + fs + "AM25" + fs + "C1TESTGROUP" +
		gs + seg + fs + "AM21" + fs + "ANP" +
		seg + fs + "AM23" + fs + "F5107B" + fs + "AY0070{"

	rsp, err := claimdeserializer.DeserializeResponse(raw)
	if err != nil {
		t.Fatalf("deserialize: %v", err)
	}
	billing := rsp.(response.Billing)

	if r := billing.Claims[0].Pricing.PercentageSalesTaxRatePaid; r == nil || *r != 0.07 {
		t.Fatalf("D0 AY=0070{ must parse as 0.07, got %v", r)
	}

	out, err := claimserializer.Serialize(&billing)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if !strings.Contains(out, fs+"AY0070{") {
		t.Errorf("D0 output must keep the overpunched AY value (0070{).\nout: %q", out)
	}
}
