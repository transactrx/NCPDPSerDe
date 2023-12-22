# NCPDP Serializer/Deserializer

## Overview
Parse and build NCPDP claims.

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