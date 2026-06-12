package ncpdp

import (
	"strings"
	"testing"
	"time"
)

type ncpdpBodyValueTest struct {
	claim        string
	recordCount  int
	segmentCount int
	fields       []fieldValueTest
}

type fieldValueTest struct {
	recordIndex int //-1 = shared segment
	segmentId   string
	fieldId     string
	fieldValue  string
}

var requestBodyTests = []ncpdpBodyValueTest{
	{
		claim:        REQUEST_B1,
		recordCount:  1,
		segmentCount: 2,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM04", fieldId: "C2", fieldValue: "POLICYNUMBERTHATISLO"},
			{recordIndex: -1, segmentId: "AM01", fieldId: "CA", fieldValue: "JOHN"},
			{recordIndex: 0, segmentId: "AM07", fieldId: "D7", fieldValue: "00172240780"},
			{recordIndex: 0, segmentId: "AM03", fieldId: "DR", fieldValue: "ERVING"},
			{recordIndex: 0, segmentId: "AM11", fieldId: "DC", fieldValue: "40{"},
		},
	},
	{
		claim:        REQUEST_B1_BATCH,
		recordCount:  4,
		segmentCount: 2,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM04", fieldId: "C2", fieldValue: "55558662900"},
			{recordIndex: -1, segmentId: "AM01", fieldId: "CA", fieldValue: "TEST*"},
			{recordIndex: 0, segmentId: "AM07", fieldId: "D7", fieldValue: "31722070010"},
			{recordIndex: 1, segmentId: "AM07", fieldId: "D7", fieldValue: "70069009101"},
			{recordIndex: 2, segmentId: "AM07", fieldId: "D7", fieldValue: "59746012110"},
			{recordIndex: 3, segmentId: "AM07", fieldId: "D7", fieldValue: "29300012510"},
			{recordIndex: 0, segmentId: "AM11", fieldId: "D9", fieldValue: "0000601I"},
			{recordIndex: 1, segmentId: "AM11", fieldId: "D9", fieldValue: "0001248G"},
			{recordIndex: 2, segmentId: "AM11", fieldId: "D9", fieldValue: "0000216G"},
			{recordIndex: 3, segmentId: "AM11", fieldId: "D9", fieldValue: "0001744B"},
		},
	},
	{
		claim:        REQUEST_E1,
		recordCount:  0,
		segmentCount: 2,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM04", fieldId: "C2", fieldValue: "D0ELIGCOB"},
			{recordIndex: -1, segmentId: "AM01", fieldId: "CA", fieldValue: "ELIGIBILITY"},
		},
	},
}

var responseBodyTests = []ncpdpBodyValueTest{
	{
		claim:        RESPONSE_PAID,
		recordCount:  1,
		segmentCount: 1,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM20", fieldId: "F4", fieldValue: "TRANSMISSION MESSAGE TEXT QS/1 TEST RESPONSE"},
			{recordIndex: 0, segmentId: "AM21", fieldId: "AN", fieldValue: "P"},
			{recordIndex: 0, segmentId: "AM22", fieldId: "D2", fieldValue: "000000010120"},
			{recordIndex: 0, segmentId: "AM23", fieldId: "F6", fieldValue: "654C"},
		},
	},
	{
		claim:        RESPONSE_REJECTED,
		recordCount:  1,
		segmentCount: 1,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM20", fieldId: "F4", fieldValue: "QS/1 POWERLINE D.0 REJECTED CLAIM TESTING TRANSMISSION LEVEL MESSAGE TEXT GOES HERE.  THE MESSAGE CAN BE UP TO 200 BYTES LONG AND SHOULD CONTAIN INFORMATION ABOUT THE TRANSMISSION OF THE CLAIM."},
			{recordIndex: 0, segmentId: "AM21", fieldId: "AN", fieldValue: "R"},
			{recordIndex: 0, segmentId: "AM22", fieldId: "D2", fieldValue: "000006508020"},
		},
	},
	{
		claim:        RESPONSE_REJECTED_BATCH,
		recordCount:  4,
		segmentCount: 1,
		fields: []fieldValueTest{
			{recordIndex: -1, segmentId: "AM20", fieldId: "F4", fieldValue: "TRANSMISSION MESSAGE TEXT QS/1 TEST RESPONSE"},
			{recordIndex: 0, segmentId: "AM21", fieldId: "AN", fieldValue: "R"},
			{recordIndex: 1, segmentId: "AM21", fieldId: "AN", fieldValue: "R"},
			{recordIndex: 2, segmentId: "AM21", fieldId: "AN", fieldValue: "R"},
			{recordIndex: 3, segmentId: "AM21", fieldId: "AN", fieldValue: "R"},
		},
	},
}

func TestParsingRequestBody(t *testing.T) {
	for _, test := range requestBodyTests {
		validateBody[RequestHeader](t, test)
	}
}

func TestParsingResponseBody(t *testing.T) {
	for _, test := range responseBodyTests {
		validateBody[ResponseHeader](t, test)
	}
}

func TestBuildingNewRequest(t *testing.T) {
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
	if err != nil {
		t.Error(err)
		return
	}

	if request.RawValue == "" {
		t.Errorf("empty RawValue")
	}

	if len(request.RawValue) <= 56 {
		t.Errorf("RawValue length mismatch. Wanted at least: %q   Got: %q", 56, len(request.RawValue))
	}
}

func validateBody[V RequestHeader | ResponseHeader](t *testing.T, test ncpdpBodyValueTest) {
	tran := NcpdpTransaction[V]{
		RawValue: test.claim,
		Created:  time.Now().UTC(),
	}
	err := tran.ParseNcpdp()

	if err != nil {
		t.Error(err)
		return
	}

	if tran.Header.RawValue == "" {
		t.Errorf("empty NcpdpTransaction.Header")
	}

	if len(tran.Records) != test.recordCount {
		t.Errorf("Record count mismatch. Wanted: %q   Got: %q", test.recordCount, len(tran.Records))
	}

	if len(tran.Segments) != test.segmentCount {
		t.Errorf("Shared segment count mismatch. Wanted: %q   Got: %q", test.segmentCount, len(tran.Segments))
	}

	validateRecords(t, tran.Records, test.fields)

	for _, seg := range tran.Segments {
		validateSegment(t, seg, -1, test.fields)
	}
}

func validateRecords(t *testing.T, records []NcpdpRecord, fields []fieldValueTest) {
	for iRecords := 0; iRecords < len(records); iRecords++ {
		record := records[iRecords]

		if record.RawValue == Empty {
			t.Errorf("empty NcpdpRecord.RawValue")
		}

		if len(record.Segments) == 0 {
			t.Errorf("empty NcpdpRecord.Segments")
		}

		for iSegments := 0; iSegments < len(record.Segments); iSegments++ {
			validateSegment(t, record.Segments[iSegments], iRecords, fields)
		}
	}
}

func validateSegment(t *testing.T, seg NcpdpSegment, recordIndex int, fields []fieldValueTest) {
	if seg.RawValue == Empty {
		t.Errorf("empty shared NcpdpSegment.RawValue")
	}

	if len(seg.Fields) == 0 {
		t.Errorf("empty NcpdpSegment.Fields")
	}

	for _, field := range seg.Fields {
		if field.RawValue == Empty {
			t.Errorf("empty shared NcpdpSegment.NcpdpField.RawValue")
		}
	}

	for _, fldTest := range fields {
		if fldTest.segmentId == seg.Id && recordIndex == fldTest.recordIndex {
			field := seg.FindFirstField(fldTest.fieldId)

			if strings.TrimSpace(field.Value) != fldTest.fieldValue {
				t.Errorf("Field value mismatch. SegmentId: %q  FieldId: %q   Wanted: %q  Got: %q", fldTest.segmentId, fldTest.fieldId, fldTest.fieldValue, field.Value)
			}
		}
	}
}

const REQUEST_B1_F6 = "F6B188015600PCN1234567101 1234567893    20260611VENDORCERT" +
	"\x1e\x1cAM04\x1cC2CARDID\x1cC1D0PAID" +
	"\x1e\x1cAM01\x1cCAFIRST\x1cCBLAST" +
	"\x1e\x1cAM07\x1cEM1\x1cD26000001\x1cE103\x1cD700172240780" +
	"\x1e\x1cAM11\x1cD9123A"

func TestParsingF6RequestBody(t *testing.T) {
	tran := NewTransactionRequest(REQUEST_B1_F6)

	err := tran.ParseNcpdp()
	if err != nil {
		t.Fatal(err)
	}

	if tran.Header.Value.Version != F6 {
		t.Errorf("Version mismatch. Wanted: %q   Got: %q", F6, tran.Header.Value.Version)
	}

	// F6 has no group separator: all segments are transaction-level
	if len(tran.Records) != 0 {
		t.Errorf("Record count mismatch. Wanted: 0   Got: %v", len(tran.Records))
	}

	if len(tran.Segments) != 4 {
		t.Errorf("Segment count mismatch. Wanted: 4   Got: %v", len(tran.Segments))
	}

	if tran.RecordCount() != 1 {
		t.Errorf("RecordCount mismatch. Wanted: 1   Got: %v", tran.RecordCount())
	}

	claimSegment := tran.FindSegmentInRecord(0, CLAIM_SEGMENT_ID)
	if claimSegment == nil {
		t.Fatal("FindSegmentInRecord(0, AM07) returned nil for F6 request")
	}

	rxField := claimSegment.FindFirstField(PRESCRIPTION_SERVICE_REFERENCE_NO_FIELD_ID)
	if rxField == nil || rxField.Value != "6000001" {
		t.Errorf("Rx number mismatch. Wanted: %q   Got: %v", "6000001", rxField)
	}

	if tran.FindSegmentInRecord(1, CLAIM_SEGMENT_ID) != nil {
		t.Error("FindSegmentInRecord(1, AM07) should be nil when transaction has no records")
	}
}

func TestRecordHelpersWithD0Batch(t *testing.T) {
	tran := NewTransactionRequest(REQUEST_B1_BATCH)

	err := tran.ParseNcpdp()
	if err != nil {
		t.Fatal(err)
	}

	if tran.RecordCount() != 4 {
		t.Errorf("RecordCount mismatch. Wanted: 4   Got: %v", tran.RecordCount())
	}

	for i := 0; i < 4; i++ {
		if tran.FindSegmentInRecord(i, CLAIM_SEGMENT_ID) == nil {
			t.Errorf("FindSegmentInRecord(%v, AM07) returned nil for D0 batch request", i)
		}
	}

	if tran.FindSegmentInRecord(4, CLAIM_SEGMENT_ID) != nil {
		t.Error("FindSegmentInRecord(4, AM07) should be nil beyond record range")
	}
}

func buildVersionedResponse(version string) NcpdpTransaction[ResponseHeader] {
	response := NewTransactionResponse("")

	response.Header.Value.Version = version
	response.Header.Value.TransactionCode = BILLING
	response.Header.Value.RecordCount = 1
	response.Header.Value.Status = ACCEPTED_STATUS
	response.Header.Value.ServiceProviderIdQualifier = "01"
	response.Header.Value.ServiceProviderId = "1234567893"
	response.Header.Value.DateOfService = "20260611"

	messageSegment := NcpdpSegment{Id: RESPONSE_MESSAGE_SEGMENT_ID}
	messageSegment.AppendField(MESSAGE_FIELD_ID, "TEST MESSAGE")
	response.Segments = append(response.Segments, messageSegment)

	statusSegment := NcpdpSegment{Id: RESPONSE_STATUS_SEGMENT_ID}
	statusSegment.AppendField(STATUS_FIELD_ID, REJECTED_STATUS)
	statusSegment.AppendField(REJECT_CODE_FIELD_ID, "85")

	record := NcpdpRecord{}
	record.Segments = append(record.Segments, statusSegment)
	response.Records = append(response.Records, record)

	return response
}

func TestBuildingF6ResponseOmitsGroupSeparator(t *testing.T) {
	response := buildVersionedResponse(F6)

	err := response.BuildNcpdp()
	if err != nil {
		t.Fatal(err)
	}

	if strings.ContainsRune(response.RawValue, rune(GROUP)) {
		t.Error("F6 response must not contain a group separator")
	}

	if !strings.Contains(response.RawValue, "AM21") {
		t.Error("F6 response is missing the record's status segment")
	}
}

func TestBuildingD0ResponseIncludesGroupSeparator(t *testing.T) {
	response := buildVersionedResponse(D0)

	err := response.BuildNcpdp()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.ContainsRune(response.RawValue, rune(GROUP)) {
		t.Error("D0 response must contain a group separator")
	}
}

func TestScanRawDataForFieldValues_Request(t *testing.T) {
	tran := NcpdpTransaction[RequestHeader]{
		RawValue: REQUEST_B1,
		Created:  time.Now().UTC(),
	}
	err := tran.ParseNcpdp()

	if err != nil {
		t.Error(err)
		return
	}

	valuesE7 := tran.ScanRawDataForFieldValues("E7")

	if len(valuesE7) != 1 {
		t.Errorf("Field count mismatch. FieldId: E7   Wanted: 1   Got: %v", len(valuesE7))
	} else {
		if valuesE7[0] != "0000001000" {
			t.Errorf("Field value mismatch. FieldId: E7   Wanted: 0000001000   Got: %v", valuesE7[0])
		}
	}
}

func TestScanRawDataForFieldValues_Response(t *testing.T) {
	tran := NcpdpTransaction[ResponseHeader]{
		RawValue: RESPONSE_REJECTED,
		Created:  time.Now().UTC(),
	}
	err := tran.ParseNcpdp()

	if err != nil {
		t.Error(err)
		return
	}

	valuesFB := tran.ScanRawDataForFieldValues("FB")

	if len(valuesFB) != 5 {
		t.Errorf("Field count mismatch. FieldId: FB   Wanted: 5   Got: %v", len(valuesFB))
	} else {
		if valuesFB[0] != "7X" {
			t.Errorf("Field value mismatch. FieldId: FB   Wanted: 7X   Got: %v", valuesFB[0])
		}

		if valuesFB[1] != "06" {
			t.Errorf("Field value mismatch. FieldId: FB   Wanted: 06   Got: %v", valuesFB[1])
		}

		if valuesFB[2] != "07" {
			t.Errorf("Field value mismatch. FieldId: FB   Wanted: 07   Got: %v", valuesFB[2])
		}

		if valuesFB[3] != "ZZ" {
			t.Errorf("Field value mismatch. FieldId: FB   Wanted: ZZ   Got: %v", valuesFB[3])
		}

		if valuesFB[4] != "372" {
			t.Errorf("Field value mismatch. FieldId: FB   Wanted: 372   Got: %v", valuesFB[4])
		}
	}
}
