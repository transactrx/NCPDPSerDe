package ncpdp

import "testing"

const REQUEST_B1 = "880151D0B1          1011234567893     20180207          AM04C2POLICYNUMBERTHATISLOCCJOHNCDDOEC1D0PAIDDURC61AM01CX99CYVERIC419340102C51CAJOHNCBDOECM9876 TESTING LANECNSPARTANBURGCOSCCP293011234CQ8642538600AM07EM1D26000001E103D700172240780E70000001000D300D530D61D80DE20180207DF00DJ1C8128EAEVPANUMU701AM03EZ01DB1234587693DRERVINGPM86458212342E01DL12345876934EERVING2JJULIUS DUNKI2K15 SLAM DUNK LANE2MPHILADELPHIA2NPA2P123456789AM11D9IDC40{DQ83DDU40IDN01"
const REQUEST_B1_BATCH = "003858D0B1MA        4011083061303     20210310CERT      AM04C255558662900         CCMINYA       CDSIDNEY         FO        C90C1RXINN01        C3   C61AM01HN                                                                                CX99CY                6211C419711111C52CATEST*       CBTEST**         CM1444 N 4TH ST APT Z           CNCOLUMBUS            COOHCP43201          CQ5559005838C7014X00AM07EM1D2000003671354E103D731722070010        E70000030000D308D5030D61D80DE20200729DF11DJ3ET0000030000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90000601IDC0000099IDX0000000{DQ0000701HDU0000701HDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671352E103D770069009101        E70000006000D308D5024D61D80DE20200729DF11DJ3ET0000006000C800DT0EK            28MLDI00EU00EV00000000000U701AM11D90001248GDC0000099IDX0000000{DQ0001348FDU0001348FDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671349E103D759746012110        E70000045000D308D5015D61D80DE20200729DF11DJ3ET0000045000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90000216GDC0000099IDX0000000{DQ0000316FDU0000316FDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671361E103D729300012510        E70000030000D308D5030D61D80DE20200729DF11DJ3ET0000030000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90001744BDC0000099IDX0000000{DQ0001844ADU0001844ADN01AM03EZ01DB1578020970     DRWALKER         2NOH"
const REQUEST_E1 = "880151D0E1          1011730433129     20210531          AM04C2D0ELIGCOBCCELIGIBILITYCDCOOLC61AM01C419420501C51CAELIGIBILITYCBCOOL4X00"

const RESPONSE_PAID = "D0B11A011234567893     20180207AM20F4TRANSMISSION MESSAGE TEXT QS/1 TEST RESPONSE AM21ANPF3123456789123456789UF1UH01FQTRANSACTION MESSAGE TEXT 7F038F6023570862AM22EM1D20000000101209F1AP03AR17236056901AM23F5100{F6654CF750{AV1J21J301J4150{F9754CFM1FN20{FI80{AW20{EQ20{*1N*2N*3N*ADX"
const RESPONSE_REJECTED = "D0B11A011619071677     20210118AM20F4QS/1 POWERLINE D.0 REJECTED CLAIM TESTING TRANSMISSION LEVEL MESSAGE TEXT GOES HERE.  THE MESSAGE CAN BE UP TO 200 BYTES LONG AND SHOULD CONTAIN INFORMATION ABOUT THE TRANSMISSION OF THE CLAIM.AM21ANRF3REJECT AUTH NUMBERFA5FB7XFB06FB07FBZZFB372UF2UH01FQFIRST QS/1 TEST REJECT MESSAGEUH02FQSECOND QS/1 TEST REJECT MESSAGE7F038F8008457558MAWWW.QS1.COMAM22EM1D20000065080209F1AP03AR50474091001AS25{AT12EAULORTAB 10-500 TABLET"
const RESPONSE_REJECTED_BATCH = "D0B14A011083061303     20210310AM20F4TRANSMISSION MESSAGE TEXT QS/1 TEST RESPONSE AM21ANRF3123456789123456789UF1UH01FQTRANSACTION MESSAGE TEXT 7F038F6023570862AM22EM1D20000036713549F1AP03AR17236056901AM21ANRF3123456789123456789UF1UH01FQTRANSACTION MESSAGE TEXT 7F038F6023570862AM22EM1D20000036713529F1AP03AR17236056901AM21ANRF3123456789123456789UF1UH01FQTRANSACTION MESSAGE TEXT 7F038F6023570862AM22EM1D20000036713499F1AP03AR17236056901AM21ANRF3123456789123456789UF1UH01FQTRANSACTION MESSAGE TEXT 7F038F6023570862AM22EM1D20000036713619F1AP03AR17236056901"

type requestHeaderValueTest struct {
	claim  string
	header RequestHeader
}

type responseHeaderValueTest struct {
	claim  string
	header ResponseHeader
}

var requestHeaderTests = []requestHeaderValueTest{
	{claim: REQUEST_B1, header: RequestHeader{Bin: "880151", Version: "D0", TransactionCode: "B1", Pcn: "", RecordCount: 1, ServiceProviderIdQualifier: "01", ServiceProviderId: "1234567893", DateOfService: "20180207", SoftwareVendorCertificationId: ""}},
	{claim: REQUEST_B1_BATCH, header: RequestHeader{Bin: "003858", Version: "D0", TransactionCode: "B1", Pcn: "MA", RecordCount: 4, ServiceProviderIdQualifier: "01", ServiceProviderId: "1083061303", DateOfService: "20210310", SoftwareVendorCertificationId: "CERT"}},
	{claim: REQUEST_E1, header: RequestHeader{Bin: "880151", Version: "D0", TransactionCode: "E1", Pcn: "", RecordCount: 1, ServiceProviderIdQualifier: "01", ServiceProviderId: "1730433129", DateOfService: "20210531", SoftwareVendorCertificationId: ""}},
}

var responseHeaderTests = []responseHeaderValueTest{
	{claim: RESPONSE_PAID, header: ResponseHeader{Version: "D0", TransactionCode: "B1", RecordCount: 1, Status: "A", ServiceProviderIdQualifier: "01", ServiceProviderId: "1234567893", DateOfService: "20180207"}},
	{claim: RESPONSE_REJECTED, header: ResponseHeader{Version: "D0", TransactionCode: "B1", RecordCount: 1, Status: "A", ServiceProviderIdQualifier: "01", ServiceProviderId: "1619071677", DateOfService: "20210118"}},
	{claim: RESPONSE_REJECTED_BATCH, header: ResponseHeader{Version: "D0", TransactionCode: "B1", RecordCount: 4, Status: "A", ServiceProviderIdQualifier: "01", ServiceProviderId: "1083061303", DateOfService: "20210310"}},
}

var buildRequestHeaderTests = []requestHeaderValueTest{
	{
		claim:  "880151D0B1TEST      1011234567893     20231127SVCID     ",
		header: RequestHeader{Bin: "880151", Version: "D0", TransactionCode: "B1", Pcn: "TEST", RecordCount: 1, ServiceProviderIdQualifier: "01", ServiceProviderId: "1234567893", DateOfService: "20231127", SoftwareVendorCertificationId: "SVCID"},
	},
	{
		claim:  "99073451E1          4070000240        20231127          ",
		header: RequestHeader{Bin: "990734", Version: "51", TransactionCode: "E1", Pcn: "", RecordCount: 4, ServiceProviderIdQualifier: "07", ServiceProviderId: "0000240", DateOfService: "20231127", SoftwareVendorCertificationId: ""},
	},
}

var buildResponseHeaderTests = []responseHeaderValueTest{
	{
		claim:  "D0B11A011234567893     20231127",
		header: ResponseHeader{Version: "D0", TransactionCode: "B1", RecordCount: 1, Status: "A", ServiceProviderIdQualifier: "01", ServiceProviderId: "1234567893", DateOfService: "20231127"},
	},
	{
		claim:  "D0E14C070000240        20231127",
		header: ResponseHeader{Version: "D0", TransactionCode: "E1", RecordCount: 4, Status: "C", ServiceProviderIdQualifier: "07", ServiceProviderId: "0000240", DateOfService: "20231127"},
	},
}

func TestParsingRequestHeader(t *testing.T) {
	for _, test := range requestHeaderTests {
		header := NcpdpHeader[RequestHeader]{
			RawValue: test.claim,
		}

		err := header.ParseNcpdpHeader()

		if err != nil {
			t.Error(err)
			break
		}

		if header.Value.Bin != test.header.Bin {
			t.Errorf("Bin mismatch. Wanted: %q   Got: %q", test.header.Bin, header.Value.Bin)
		}

		if header.Value.Version != test.header.Version {
			t.Errorf("Version mismatch. Wanted: %q   Got: %q", test.header.Version, header.Value.Version)
		}

		if header.Value.TransactionCode != test.header.TransactionCode {
			t.Errorf("TransactionCode mismatch. Wanted: %q   Got: %q", test.header.TransactionCode, header.Value.TransactionCode)
		}

		if header.Value.Pcn != test.header.Pcn {
			t.Errorf("TransactionCode mismatch. Wanted: %q   Got: %q", test.header.Pcn, header.Value.Pcn)
		}

		if header.Value.RecordCount != test.header.RecordCount {
			t.Errorf("RecordCount mismatch. Wanted: %q   Got: %q", test.header.RecordCount, header.Value.RecordCount)
		}

		if header.Value.ServiceProviderIdQualifier != test.header.ServiceProviderIdQualifier {
			t.Errorf("ServiceProviderIdQualifier mismatch. Wanted: %q   Got: %q", test.header.ServiceProviderIdQualifier, header.Value.ServiceProviderIdQualifier)
		}

		if header.Value.ServiceProviderId != test.header.ServiceProviderId {
			t.Errorf("ServiceProviderId mismatch. Wanted: %q   Got: %q", test.header.ServiceProviderId, header.Value.ServiceProviderId)
		}

		if header.Value.DateOfService != test.header.DateOfService {
			t.Errorf("DateOfService mismatch. Wanted: %q   Got: %q", test.header.DateOfService, header.Value.DateOfService)
		}

		if header.Value.SoftwareVendorCertificationId != test.header.SoftwareVendorCertificationId {
			t.Errorf("SoftwareVendorCertificationId mismatch. Wanted: %q   Got: %q", test.header.SoftwareVendorCertificationId, header.Value.SoftwareVendorCertificationId)
		}

		if header.Size != 56 {
			t.Errorf("Header size mismatch. Wanted: '56'   Got: %q", header.Size)
		}
	}
}

func TestParsingResponseHeader(t *testing.T) {
	for _, test := range responseHeaderTests {
		header := NcpdpHeader[ResponseHeader]{
			RawValue: test.claim,
		}

		err := header.ParseNcpdpHeader()

		if err != nil {
			t.Error(err)
			break
		}

		if header.Value.Version != test.header.Version {
			t.Errorf("Version mismatch. Wanted: %q   Got: %q", test.header.Version, header.Value.Version)
		}

		if header.Value.TransactionCode != test.header.TransactionCode {
			t.Errorf("TransactionCode mismatch. Wanted: %q   Got: %q", test.header.TransactionCode, header.Value.TransactionCode)
		}

		if header.Value.RecordCount != test.header.RecordCount {
			t.Errorf("RecordCount mismatch. Wanted: %q   Got: %q", test.header.RecordCount, header.Value.RecordCount)
		}

		if header.Value.Status != test.header.Status {
			t.Errorf("Status mismatch. Wanted: %q   Got: %q", test.header.Status, header.Value.Status)
		}

		if header.Value.ServiceProviderIdQualifier != test.header.ServiceProviderIdQualifier {
			t.Errorf("ServiceProviderIdQualifier mismatch. Wanted: %q   Got: %q", test.header.ServiceProviderIdQualifier, header.Value.ServiceProviderIdQualifier)
		}

		if header.Value.ServiceProviderId != test.header.ServiceProviderId {
			t.Errorf("ServiceProviderId mismatch. Wanted: %q   Got: %q", test.header.ServiceProviderId, header.Value.ServiceProviderId)
		}

		if header.Value.DateOfService != test.header.DateOfService {
			t.Errorf("DateOfService mismatch. Wanted: %q   Got: %q", test.header.DateOfService, header.Value.DateOfService)
		}

		if header.Size != 31 {
			t.Errorf("Header size mismatch. Wanted: '31'   Got: %q", header.Size)
		}
	}
}

func TestBuildingRequestHeader(t *testing.T) {
	for _, test := range buildRequestHeaderTests {
		header := NcpdpHeader[RequestHeader]{
			Value: test.header,
		}

		err := header.buildNcpdpHeader()
		if err != nil {
			t.Error(err)
			break
		}

		if header.RawValue != test.claim {
			t.Errorf("RawValue mismatch. Wanted: %q   Got: %q", test.claim, header.RawValue)
		}
	}
}

func TestBuildingResponseHeader(t *testing.T) {
	for _, test := range buildResponseHeaderTests {
		header := NcpdpHeader[ResponseHeader]{
			Value: test.header,
		}

		err := header.buildNcpdpHeader()
		if err != nil {
			t.Error(err)
			break
		}

		if header.RawValue != test.claim {
			t.Errorf("RawValue mismatch. Wanted: %q   Got: %q", test.claim, header.RawValue)
		}
	}
}
