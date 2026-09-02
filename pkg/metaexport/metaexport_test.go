package metaexport

import (
	"testing"
)

func buildOrFail(t *testing.T) *Metadata {
	t.Helper()

	meta, err := Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	return meta
}

func findField(t *testing.T, meta *Metadata, typeName, fieldName string) Field {
	t.Helper()

	typeMeta, ok := meta.Types[typeName]
	if !ok {
		t.Fatalf("type %v not exported", typeName)
	}

	for _, f := range typeMeta.Fields {
		if f.Name == fieldName {
			return f
		}
	}

	t.Fatalf("field %v.%v not exported", typeName, fieldName)
	return Field{}
}

func TestBuildExportsAllRegisteredTransactions(t *testing.T) {
	meta := buildOrFail(t)

	if len(meta.Transactions) != 36 {
		t.Errorf("expected 36 registered transactions, got %d", len(meta.Transactions))
	}

	if meta.Transactions["B1|request"] != "request.Billing" {
		t.Errorf("B1|request mapped to %q", meta.Transactions["B1|request"])
	}

	if meta.Transactions["B1|response"] != "response.Billing" {
		t.Errorf("B1|response mapped to %q", meta.Transactions["B1|response"])
	}
}

func TestBillingTransactionShape(t *testing.T) {
	meta := buildOrFail(t)

	if meta.Types["request.Billing"].Kind != "transaction" {
		t.Errorf("request.Billing kind = %q", meta.Types["request.Billing"].Kind)
	}

	header := findField(t, meta, "request.Billing", "Header")
	if header.Role != "header" || header.Type != "ncpdp.RequestHeader" {
		t.Errorf("Header exported as %+v", header)
	}

	insurance := findField(t, meta, "request.Billing", "Insurance")
	if insurance.Role != "segment" || insurance.Code != "AM04" || insurance.Order != 1 || insurance.Type != "requestsegment.Insurance" {
		t.Errorf("Insurance exported as %+v", insurance)
	}

	claims := findField(t, meta, "request.Billing", "Claims")
	if claims.Role != "group" || claims.Max != 4 || claims.Type != "request.BillingRecord" {
		t.Errorf("Claims exported as %+v", claims)
	}

	dynamicSegments := findField(t, meta, "request.Billing", "DynamicSegments")
	if dynamicSegments.Role != "dynamicSegments" {
		t.Errorf("DynamicSegments exported as %+v", dynamicSegments)
	}

	raw := findField(t, meta, "request.BillingRecord", "Raw")
	if raw.Role != "raw" {
		t.Errorf("BillingRecord.Raw exported as %+v", raw)
	}
}

func TestRequestHeaderLayouts(t *testing.T) {
	meta := buildOrFail(t)

	if meta.Types["ncpdp.RequestHeader"].Kind != "header" {
		t.Errorf("ncpdp.RequestHeader kind = %q", meta.Types["ncpdp.RequestHeader"].Kind)
	}

	bin := findField(t, meta, "ncpdp.RequestHeader", "Bin")
	if bin.Layout == nil || *bin.Layout != (Layout{Start: 0, Length: 6, Order: 1}) {
		t.Errorf("Bin layout exported as %+v", bin.Layout)
	}
	if bin.LayoutF6 == nil || *bin.LayoutF6 != (Layout{Start: 4, Length: 8, Order: 3}) {
		t.Errorf("Bin layoutF6 exported as %+v", bin.LayoutF6)
	}

	recordCount := findField(t, meta, "ncpdp.RequestHeader", "RecordCount")
	if recordCount.GoType != "int" {
		t.Errorf("RecordCount goType = %q", recordCount.GoType)
	}
}

func TestLeafFieldAttributes(t *testing.T) {
	meta := buildOrFail(t)

	quantity := findField(t, meta, "requestsegment.Claim", "QuantityDispensed")
	if quantity.Role != "field" || quantity.Code != "E7" || quantity.DecimalPlaces != 3 || quantity.Order != 10 || quantity.GoType != "float64" {
		t.Errorf("QuantityDispensed exported as %+v", quantity)
	}

	dateWritten := findField(t, meta, "requestsegment.Claim", "DateWritten")
	if dateWritten.Format != "YYYYMMdd" || dateWritten.GoType != "time" {
		t.Errorf("DateWritten exported as %+v", dateWritten)
	}

	ingredientCost := findField(t, meta, "requestsegment.Pricing", "IngredientCostSubmitted")
	if !ingredientCost.Overpunch || ingredientCost.DecimalPlaces != 2 || ingredientCost.Code != "D9" {
		t.Errorf("IngredientCostSubmitted exported as %+v", ingredientCost)
	}

	segmentId := findField(t, meta, "ncpdp.SegmentId", "Id")
	if segmentId.Code != "AM" || segmentId.Order != 1 {
		t.Errorf("SegmentId.Id exported as %+v", segmentId)
	}

	segmentIdRaw := findField(t, meta, "ncpdp.SegmentId", "Raw")
	if segmentIdRaw.Role != "raw" {
		t.Errorf("SegmentId.Raw exported as %+v", segmentIdRaw)
	}
}

func TestNestedAndRepeatingFields(t *testing.T) {
	meta := buildOrFail(t)

	cardholder := findField(t, meta, "requestsegment.Insurance", "Cardholder")
	if cardholder.Role != "struct" || cardholder.Type != "requestsegment.Cardholder" {
		t.Errorf("Cardholder exported as %+v", cardholder)
	}

	otherPayer := findField(t, meta, "requestsegment.CoordinationOfBenefits", "OtherPayer")
	if otherPayer.Role != "repeating" || otherPayer.Type != "requestsegment.OtherPayer" {
		t.Errorf("OtherPayer exported as %+v", otherPayer)
	}

	dynamicFields := findField(t, meta, "requestsegment.Claim", "DynamicFields")
	if dynamicFields.Role != "dynamicFields" {
		t.Errorf("DynamicFields exported as %+v", dynamicFields)
	}

	// A code-less sinceVersion field tag must not demote a repeating group to a
	// leaf field.
	patientIds := findField(t, meta, "requestsegment.Patient", "Ids")
	if patientIds.Role != "repeating" || patientIds.Type != "requestsegment.PatientId" {
		t.Errorf("Patient.Ids exported as %+v", patientIds)
	}

	payers := findField(t, meta, "responsesegment.Insurance", "Payers")
	if payers.Role != "repeating" || payers.Type != "responsesegment.Payer" {
		t.Errorf("Insurance.Payers exported as %+v", payers)
	}
}

// Version scoping must reach non-Go consumers: an F6-only leaf, an F6-only
// derived counter, and an F6 repeating group all export sinceVersion, while
// D0 fields export nothing.
func TestSinceVersionExported(t *testing.T) {
	meta := buildOrFail(t)

	species := findField(t, meta, "requestsegment.Patient", "Species")
	if species.SinceVersion != "F6" {
		t.Errorf("Species exported as %+v", species)
	}

	payerIdCount := findField(t, meta, "responsesegment.Insurance", "PayerIdCount")
	if payerIdCount.SinceVersion != "F6" || payerIdCount.CountFor != "Payers" {
		t.Errorf("PayerIdCount exported as %+v", payerIdCount)
	}

	patientIds := findField(t, meta, "requestsegment.Patient", "Ids")
	if patientIds.SinceVersion != "F6" || patientIds.Role != "repeating" {
		t.Errorf("Patient.Ids exported as %+v", patientIds)
	}

	firstName := findField(t, meta, "requestsegment.Patient", "FirstName")
	if firstName.SinceVersion != "" {
		t.Errorf("D0 field FirstName must not carry sinceVersion: %+v", firstName)
	}
}
