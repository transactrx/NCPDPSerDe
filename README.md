# NCPDP Serializer/Deserializer

## Overview
Parse and build NCPDP telecommunication claims. Supports **D.0** and **F6** (vF6) formats — the version is auto-detected from the header on parse, and selected from the `Version` field on build. No format flags needed.

> **Migrating an existing app to F6?** See [F6_MIGRATION.md](F6_MIGRATION.md) — especially the warning about `tran.Records[0]` being empty for F6.

### Parse Examples (Generic Parser)
Request:
```
requestTran := ncpdp.NewTransactionRequest(rawClaimString)
err := requestTran.ParseNcpdp()
```

Response:
```
responseTran := ncpdp.NewTransactionResponse(rawClaimString)
err := responseTran.ParseNcpdp()
```

### Build Examples (Generic Builder)
Update a Request Header value:
```
requestTran.Header.Value.Bin = "123456"
requestTran.Header.Value.Pcn = "TESTPCN"

//Rebuild raw claim
err := requestTran.BuildNcpdp()
```
Update an existing field value:
```
groupField := requestTran.FindFirstField(ncpdp.INSURANCE_SEGMENT_ID, ncpdp.GROUP_CODE_FIELD_ID, -1)
if groupField != nil {    
    groupField.Value = "NEWVALUE"
}

//Rebuild raw claim 
err := requestTran.BuildNcpdp()
```

Create a new Request:
```
request := NewTransactionRequest("")

// Populate header
request.Header.Value.Bin = "880151"
request.Header.Value.Version = D0
request.Header.Value.TransactionCode = REVERSAL
request.Header.Value.Pcn = "TEST"
request.Header.Value.RecordCount = 1
request.Header.Value.ServiceProviderIdQualifier = "01"
request.Header.Value.ServiceProviderId = "1234567893"
request.Header.Value.DateOfService = "20231201"
request.Header.Value.SoftwareVendorCertificationId = "CERTID"

// Add shared segments
insuranceSegment := NcpdpSegment{Id: INSURANCE_SEGMENT_ID}
insuranceSegment.AppendField(CARDHOLDER_ID_FIELD_ID, "card_id")
insuranceSegment.AppendField(GROUP_CODE_FIELD_ID, "group_code")

request.Segments = append(request.Segments, insuranceSegment)

// Add groups/records
claimRecord := NcpdpRecord{}
claimSegment := NcpdpSegment{Id: CLAIM_SEGMENT_ID}
claimSegment.AppendField(PRESCRIPTION_SERVICE_REFERENCE_NO_QUALIFIER_FIELD_ID, "01")
claimSegment.AppendField(PRESCRIPTION_SERVICE_REFERENCE_NO_FIELD_ID, "rx_number")
claimSegment.AppendField(PRODUCT_SERVICE_ID_QUALIFIER_FIELD_ID, "03")
claimSegment.AppendField(PRODUCT_SERVICE_ID_FIELD_ID, "drug_ndc")

claimRecord.Segments = append(claimRecord.Segments, claimSegment)

request.Records = append(request.Records, claimRecord)

err := request.BuildNcpdp()
```

Build an **F6** Request (single transaction, no group separator):
```go
request := ncpdp.NewTransactionRequest("")

// Populate header — version "F6" drives the F6 layout and omits 0x1D on build
request.Header.Value.Version = ncpdp.F6
request.Header.Value.TransactionCode = ncpdp.BILLING
request.Header.Value.Bin = "88015600" // 101-A1 IIN in F6 is 8 chars
request.Header.Value.Pcn = "PCN1234567"
request.Header.Value.RecordCount = 1
request.Header.Value.ServiceProviderIdQualifier = "01"
request.Header.Value.ServiceProviderId = "1234567893"
request.Header.Value.DateOfService = "20260611"
request.Header.Value.SoftwareVendorCertificationId = "VENDORCERT"

// F6 allows exactly one transaction per transmission; BuildNcpdp errors on >1.
claimRecord := ncpdp.NcpdpRecord{}
// ... populate claimRecord.Segments ...
request.Records = append(request.Records, claimRecord)

err := request.BuildNcpdp()
```

### ⚠️ Records vs. Segments on F6
For the generic parser, F6 transmissions have **no group separator (`0x1D`)**, so all segments land in `tran.Segments` and `tran.Records` stays **empty**. Direct indexing like `tran.Records[0]` will misbehave or panic on F6.

Use the F6-safe helpers instead — they fall back to shared segments when `Records` is empty, so the same code works for both D0 and F6:
```go
// ✅ Works for both D0 and F6
seg := tran.FindSegmentInRecord(0, ncpdp.RESPONSE_STATUS_SEGMENT_ID)
for i := 0; i < tran.RecordCount(); i++ {
    seg := tran.FindSegmentInRecord(i, ncpdp.CLAIM_SEGMENT_ID)
    ...
}

// FindFirstField and InsertField apply the same fallback for record index >= 0
rxField := tran.FindFirstField(ncpdp.CLAIM_SEGMENT_ID, ncpdp.PRESCRIPTION_SERVICE_REFERENCE_NO_FIELD_ID, 0)
```

### Response Helpers
Convenience methods on `NcpdpTransaction[ResponseHeader]` — all F6-aware:
```go
responseTran := ncpdp.NewTransactionResponse(rawClaimString)
_ = responseTran.ParseNcpdp()

status := responseTran.Status()             // "P", "R", "D", etc.
isPaid := responseTran.IsPaid()             // P or D
isRejected := responseTran.IsRejected()     // R
rejectCodes := responseTran.GetRejectCodes() // []string from FB fields
messages := responseTran.GetAdditionalMessages() // map[qualifier]message from UH/FQ pairs
```

### Strongly-Typed Parser/Builder
For type-safe access to NCPDP fields (vs. iterating raw segments), use the `claimDeserializer` / `claimSerializer` packages. They return concrete request/response structs (`request.BillingRequest`, `response.BillingResponse`, etc.) with typed fields, support both D0 and F6, and round-trip cleanly.

Deserialize:
```go
import "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"

obj, err := claimDeserializer.DeserializeRequest(rawClaimString)
billing, ok := obj.(*request.BillingRequest)
if ok {
    rxNumber := billing.Claims[0].Claim.PrescriptionServiceReferenceNo
    // F6: Claims[0] is still populated — the strongly-typed parser locates the
    // single transaction by segment ID when no group separator is present.
}
```

Serialize:
```go
import "github.com/transactrx/NCPDPSerDe/pkg/claimSerializer"

raw, err := claimSerializer.Serialize(billing)
// Header version drives format: "D0" emits 0x1D group separators,
// "F6" omits them and errors if more than one claim group is present.
```

### Repeating Fields (CX/CY, J7/J8) on F6
F6 permits patient IDs (CX/CY) and insurance payers (J7/J8) to repeat. The strongly-typed structs expose both shapes for backward compatibility:

| Segment | Scalar (D0, preserved) | Slice (F6 repeats) |
|---|---|---|
| `request.Patient` | `IdQualifier`, `Id` | `Ids []PatientId` |
| `response.Insurance` | `Payer` | `Payers []Payer` |

- **Reading:** scalar holds the last occurrence; slice holds every occurrence. Read the slice when full F6 data matters.
- **Writing:** when the slice is non-empty, the serializer skips the scalar's codes to prevent double-emit. Populate one or the other.