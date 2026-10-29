package claimdeserializer

import (
	"strings"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// buildTwoPayerCobRequest swaps REQUEST_B1's single-payer AM05 for one with
// two other payers, each carrying its own amount-paid (HB/HC/DV) and patient
// responsibility (NR/NP/NQ) group.
func buildTwoPayerCobRequest(t *testing.T) string {
	t.Helper()

	fs := string(ncpdp.FIELD)

	oldAM05 := "AM05" + fs + "4C1" + fs + "5C01" + fs + "6C03" + fs + "7C015581    " + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000000{" + fs + "5E02" + fs + "6E70 " + fs + "6EA5 "
	newAM05 := "AM05" + fs + "4C2" +
		fs + "5C01" + fs + "6C03" + fs + "7C01111" + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000100{" + fs + "NR01" + fs + "NP06" + fs + "NQ0000050{" +
		fs + "5C02" + fs + "6C03" + fs + "7C02222" + fs + "E820231219" + fs + "HB1" + fs + "HC07" + fs + "DV0000200{" + fs + "NR01" + fs + "NP06" + fs + "NQ0000060{"

	raw := strings.Replace(REQUEST_B1, oldAM05, newAM05, 1)
	if raw == REQUEST_B1 {
		t.Fatal("AM05 replacement did not apply to REQUEST_B1")
	}

	return raw
}

// A nested repeating group under the second occurrence of an outer repeating
// group must be assigned within that occurrence. Field occurrence counts span
// the whole segment, so payer 2's first HC/DV (each the segment's second
// occurrence of its code) used to split into two half-filled amount elements,
// and the rebuilt HB/NR counts doubled.
func Test_CobSecondPayerKeepsAmountPairsTogether(t *testing.T) {
	obj := request.Billing{}
	if err := DeserializeType(buildTwoPayerCobRequest(t), &obj); err != nil {
		t.Fatal(err)
	}

	cob := obj.Claims[0].CoordinationOfBenefits

	if len(cob.OtherPayer) != 2 {
		t.Fatalf("OtherPayer count mismatch. Wanted: 2   Got: %v", len(cob.OtherPayer))
	}

	for i, want := range []struct {
		amount         float64
		responsibility float64
	}{
		{10.00, 5.00},
		{20.00, 6.00},
	} {
		payer := cob.OtherPayer[i]

		if len(payer.OtherPayerAmountsPaid) != 1 {
			t.Fatalf("Payer %v OtherPayerAmountsPaid count mismatch. Wanted: 1   Got: %+v", i+1, payer.OtherPayerAmountsPaid)
		}
		paid := payer.OtherPayerAmountsPaid[0]
		if paid.Qualifier == nil || *paid.Qualifier != "07" || paid.AmountPaid == nil || *paid.AmountPaid != want.amount {
			t.Errorf("Payer %v amount pair mismatch. Wanted: {07 %v}   Got: %+v", i+1, want.amount, paid)
		}
		if payer.OtherPayerAmountPaidCount == nil || *payer.OtherPayerAmountPaidCount != 1 {
			t.Errorf("Payer %v HB count mismatch. Wanted: 1   Got: %v", i+1, payer.OtherPayerAmountPaidCount)
		}

		if len(payer.OtherPayerPatientResponsibilityAmounts) != 1 {
			t.Fatalf("Payer %v responsibility count mismatch. Wanted: 1   Got: %+v", i+1, payer.OtherPayerPatientResponsibilityAmounts)
		}
		resp := payer.OtherPayerPatientResponsibilityAmounts[0]
		if resp.Qualifier == nil || *resp.Qualifier != "06" || resp.Amount == nil || *resp.Amount != want.responsibility {
			t.Errorf("Payer %v responsibility pair mismatch. Wanted: {06 %v}   Got: %+v", i+1, want.responsibility, resp)
		}
	}
}
