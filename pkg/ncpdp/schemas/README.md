# NCPDP D.0/F6 JSON Schemas

This directory contains JSON Schema definitions for NCPDP Telecommunication Standard Version D.0 and F6 transactions (the version is auto-detected from the header). These schemas can be used by external users to:

1. Understand the data structure of NCPDP transactions
2. Validate JSON payloads before serialization
3. Generate client code in various languages
4. Document API contracts

## Directory Structure

```
schemas/
├── ncpdp-schemas.json          # Shared type definitions ($defs)
├── request/                     # Request transaction schemas
│   ├── billing.json            # B1 - Billing Request
│   ├── reversal.json           # B2 - Reversal Request
│   ├── rebill.json             # B3 - Rebill Request
│   ├── eligibility.json        # E1 - Eligibility Request
│   ├── information.json        # N1 - Information Reporting
│   ├── priorAuthorization.json # P1 - Prior Authorization
│   ├── predeterminationOfBenefits.json # D1 - Predetermination
│   ├── serviceBilling.json     # S1 - Service Billing
│   └── controlledSubstanceReporting.json # C1 - Controlled Substance
├── response/                    # Response transaction schemas
│   ├── billing.json            # B1 - Billing Response
│   ├── reversal.json           # B2 - Reversal Response
│   ├── rebill.json             # B3 - Rebill Response
│   ├── eligibility.json        # E1 - Eligibility Response
│   ├── information.json        # N1 - Information Reporting Response
│   ├── priorAuthorization.json # P1 - Prior Authorization Response
│   ├── predeterminationOfBenefits.json # D1 - Predetermination Response
│   ├── serviceBilling.json     # S1 - Service Billing Response
│   └── controlledSubstanceReporting.json # C1 - Controlled Substance Response
└── examples/                    # Example JSON payloads
    ├── billing-request-example.json
    └── billing-response-example.json
```

## Transaction Codes

| Code | Transaction Type |
|------|------------------|
| B1   | Billing (New Rx) |
| B2   | Reversal |
| B3   | Rebill |
| E1   | Eligibility Verification |
| N1   | Information Reporting |
| N2   | Information Reporting Rebill |
| N3   | Information Reporting Reversal |
| P1   | Prior Authorization Request & Billing |
| P2   | Prior Authorization Reversal |
| P3   | Prior Authorization Inquiry |
| P4   | Prior Authorization Request Only |
| D1   | Predetermination of Benefits |
| S1   | Service Billing |
| S2   | Service Rebill |
| S3   | Service Reversal |
| C1   | Controlled Substance Reporting |
| C2   | Controlled Substance Reporting Rebill |
| C3   | Controlled Substance Reporting Reversal |

## Key Concepts

### Header
Every transaction has a Header containing:
- **Bin**: Bank Identification Number (D.0, 6 chars) or IIN (F6, 8 chars)
- **Version**: "D0" or "F6"
- **TransactionCode**: Identifies the transaction type (B1, B2, etc.)
- **RecordCount**: Number of claims in the transaction (D.0 allows 1-4; F6 allows exactly 1)
- **ServiceProviderId**: Pharmacy NPI or other identifier
- **DateOfService**: Date in CCYYMMDD format

### D.0 / F6 Versions
Every D.0 field is unchanged in F6 (full backward compatibility). Fields and segments that only exist in F6 are suffixed `- F6` in their schema descriptions and are optional — D.0 transactions simply omit them. Schema constraints use the largest value across versions (e.g. `Bin` is `maxLength: 8` for the F6 IIN even though D.0 BINs are 6 characters). F6 permits exactly one claim per transmission; the serializer rejects F6 transactions with more than one claim group.

### Segments
Transactions consist of segments identified by codes (AM01, AM04, etc.):
- **AM01**: Patient Segment
- **AM02**: Pharmacy Provider Segment
- **AM03**: Prescriber Segment
- **AM04**: Insurance Segment
- **AM05**: Coordination of Benefits Segment
- **AM06**: Workers Compensation Segment
- **AM07**: Claim Segment
- **AM08**: DUR/PPS Segment
- **AM09**: Coupon Segment
- **AM10**: Compound Segment
- **AM11**: Pricing Segment
- **AM12**: Prior Authorization Segment
- **AM13**: Clinical Segment
- **AM14**: Additional Documentation Segment
- **AM15**: Facility Segment
- **AM16**: Narrative Segment
- **AM20**: Message Segment (Response)
- **AM21**: Response Status Segment
- **AM22**: Response Claim Segment
- **AM23**: Response Pricing Segment
- **AM24**: Response DUR/PPS Segment
- **AM25**: Response Insurance Segment
- **AM26**: Response Prior Authorization Segment
- **AM27**: Response Insurance Additional Information Segment
- **AM28**: Response Coordination of Benefits Segment
- **AM29**: Response Patient Segment

### Field Codes
Each field within a segment has a 2-character code (e.g., C2=Cardholder ID, D2=Rx Number). These codes correspond to NCPDP field identifiers.

### Internal Fields
The Go structs contain internal bookkeeping fields (`SegmentId`, header `RawValue`/`Size`, and `Raw` group captures) that are excluded from JSON in both directions: they are not accepted as input (segment identifiers are populated automatically during serialization) and are not emitted as output. They are intentionally absent from these schemas.

### Counter Fields
Repeating-group count fields (e.g. `Count` 4C in Coordination of Benefits, `RejectCodeCount` FA) and per-item counters (e.g. DUR `Counter` 7E) are derived automatically during serialization: counts are set to the array length (correcting any supplied value) and per-item counters are numbered 1..N. Supply a count yourself only when using the D0 scalar form instead of the repeating array (e.g. `Patient.IdCount` alongside the single `IdQualifier`/`Id` fields).

### Nullable Fields
Most fields are nullable (can be omitted or set to null). Only required fields in the schema must be provided.

### DynamicFields
Each segment includes a `DynamicFields` array for non-standard or vendor-specific fields; transactions and claim groups likewise include a `DynamicSegments` array for segments not defined in the standard. Both marshal as objects of the form `{"Value": {...}}`, where each key of the inner object encodes the NCPDP field code and its position within the segment as `Field_<code>_<order>` (characters not valid in identifiers are replaced with words, e.g. `Field_ampersandB_25` for code `&B`). Dynamic segments also carry a `Raw` key holding the raw segment capture. These names are how dynamic data survives a JSON round-trip: unmarshaling reconstructs the underlying dynamic type from them, so a JSON payload containing dynamic fields or segments can be serialized back to NCPDP format. Keep the generated key names intact when editing such payloads.

## Usage Examples

### Go
```go
import (
    "encoding/json"
    request "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// Create a billing request
billing := request.Billing{
    Header: ncpdp.NcpdpHeader[ncpdp.RequestHeader]{
        Value: ncpdp.RequestHeader{
            Bin:             "610014",
            Version:         "D0",
            TransactionCode: "B1",
            // ...
        },
    },
    // ...
}

// Serialize to JSON
jsonBytes, _ := json.Marshal(billing)
```

### JavaScript/TypeScript
```typescript
// Use JSON Schema to generate types
// npm install json-schema-to-typescript

interface BillingRequest {
  Header: {
    Value: {
      Bin: string;
      Version: string;
      TransactionCode: string;
      // ...
    };
  };
  Insurance: InsuranceSegment;
  Patient: PatientSegment;
  Claims: BillingRecord[];
}
```

### Python
```python
import json

# Load and validate against schema
billing_request = {
    "Header": {
        "Value": {
            "Bin": "610014",
            "Version": "D0",
            "TransactionCode": "B1",
            # ...
        }
    },
    # ...
}

json_str = json.dumps(billing_request)
```

## Response Status Codes

| Code | Meaning |
|------|---------|
| A    | Accepted/Approved |
| C    | Captured (claim accepted, payment to follow) |
| D    | DUR Reject (Drug Utilization Review) |
| P    | Paid |
| R    | Rejected |

## Common Reject Codes

See NCPDP External Code List (ECL) for complete reject code definitions. Common codes include:
- 70: Product/Service Not Covered
- 75: Prior Authorization Required
- 79: Refill Too Soon
- 88: DUR Reject

## Date/Time Formats

- Dates: CCYYMMDD (e.g., "20240115")
- Times: HHmm (e.g., "1430") or HHmmss (e.g., "143025")

## Monetary Amounts

All monetary amounts are represented as numbers with 2 decimal places. Values may be positive or negative (overpunch representation in raw NCPDP format is handled internally).

## Documentation (Confluence)

A Word reference document (`NCPDP-JSON-Schema-Reference.docx` in the repo root) is generated from these schemas for upload to Confluence. **After changing anything in this directory, regenerate it** from the repo root:

```
python scripts/gen_schema_docx.py
```

See [scripts/README.md](../../../scripts/README.md) for prerequisites, options, and Confluence upload steps. When adding fields, keep the NCPDP field code as a trailing parenthesized suffix in the `description` (e.g. `"Cardholder ID (C2)"`) — the generator extracts it into the documentation's field-code column.

## Questions?

For questions about NCPDP standards, refer to the official NCPDP Telecommunication Standard Implementation Guide.
