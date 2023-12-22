package claimparser

import (
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

const REQUEST_B1 = "880151D0B1          1011234567893     20231219          AM04C2POLICYNUMBERTHATISLOCCJOHNCDDOEC1D0PAIDDURC61AM01CX99CYVERIC419341231C51CAJOHNCBDOECM9876 TESTING LANECNSPARTANBURGCOSCCP293011234CQ8642538600AM07EM1D26000001E103D700172240780E70000001000D300D530D61D80DE20231219DF00DJ1NX2DK20DK21C8128EAEVPANUMU701AM02EY05E9               AM03EZ01DB1234587693DRERVINGPM86458212342E01DL12345876934EERVING2JJULIUS DUNKI2K15 SLAM DUNK LANE2MPHILADELPHIA2NPA2P123456789AM11D9IDC40{DQ83DDU40IDN01AM13VE2WE02DOM06.9          WE02DOZ79.899        XE1ZE20231219H11358H2ABH3CDH4EFGHIJKLMNOPXE2ZE20231218H10835H2QRH3STH4UVXYZAM06DY20230203CFREDSAIL TECHNOLOGIES INC.CG                              CH                    CI  CJ               DZ0790001976AM054C15C016C037C015581    E820231219HB1HC07DV0000000{5E026E70 6EA5 AM142Q0112V200709152U12S42R62T200709152Z114B1A4KJ29204B1B4J404B1C4J14B44KY4B5A4K14B5B4K34B84KHEART INSTITUTE4B94KHEARTSVILLE4B104KMO4B114G200709114B124KNAM158CVPC283QVISTA PACIFIC CTR3U8888 PACIFIC AVE5JJURUPA VALLEY3VCA6D92509AM16BMTBU WITH BOE E0570 DOP12/2018AM087E1E4TDE5M0E61G8E117E2E4HDE5M0E61G8E117E3E4DDE5M0E61G8E11J9QFH6COAGENT-IDAM09KE01ME321034NE100{AM10EF03EG2EC03RE03TE51927177800        ED0000000003EE0000005AUE01RE03TE38779016308        ED0000000030EE0000036HUE01RE03TE62991277601        ED0000030000EE0000380HUE01AMXX&BDICHD&CB&FEA590&G[D&H00&ID&JEA590RCP&K145023&LN&MN&N1992211627&U410&V1013&W3&Y9GYB9&Z20241218#A3#BE#CFS8133578#E00005#F073416#G00228202950#H595560106973525236#I1934959901#KFS7324558#L1730616749#M466#N -#OY#PN#WALPR5CBC#XPANCIERA#YPENNY#ZPP!F5000!MCT"
const REQUEST_B2 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B3 = "880151D0B3PCN       1011234567893     20231219cert      AM04C2POLICYNUMBERTHATISLOCCJOHNCDDOEC1D0PAIDDURC61AM01CX99CYVERIC419341231C51CAJOHNCBDOECM9876 TESTING LANECNSPARTANBURGCOSCCP293011234CQ8642538600AM07EM1D26000001E103D700172240780E70000001000D300D530D61D80DE20231219DF00DJ1NX2DK20DK21C8128EAEVPANUMU701AM02EY05E9               AM03EZ01DB1234587693DRERVINGPM86458212342E01DL12345876934EERVING2JJULIUS DUNKI2K15 SLAM DUNK LANE2MPHILADELPHIA2NPA2P123456789AM11D9IDC40{DQ83DDU40IDN01AM13VE2WE02DOM06.9          WE02DOZ79.899        XE1ZE20231219H11358H2ABH3CDH4EFGHIJKLMNOPXE2ZE20231218H10835H2QRH3STH4UVXYZAM06DY20230203CFREDSAIL TECHNOLOGIES INC.CG                              CH                    CI  CJ               DZ0790001976AM054C15C016C037C015581    E820231219HB1HC07DV0000000{5E026E70 6EA5 AM142Q0112V200709152U12S42R62T200709152Z114B1A4KJ29204B1B4J404B1C4J14B44KY4B5A4K14B5B4K34B84KHEART INSTITUTE4B94KHEARTSVILLE4B104KMO4B114G200709114B124KNAM158CVPC283QVISTA PACIFIC CTR3U8888 PACIFIC AVE5JJURUPA VALLEY3VCA6D92509AM16BMTBU WITH BOE E0570 DOP12/2018AM087E1E4TDE5M0E61G8E117E2E4HDE5M0E61G8E117E3E4DDE5M0E61G8E11J9QFH6COAGENT-IDAM09KE01ME321034NE100{AM10EF03EG2EC03RE03TE51927177800        ED0000000003EE0000005AUE01RE03TE38779016308        ED0000000030EE0000036HUE01RE03TE62991277601        ED0000030000EE0000380HUE01AMXX&BDICHD&CB&FEA590&G[D&H00&ID&JEA590RCP&K145023&LN&MN&N1992211627&U410&V1013&W3&Y9GYB9&Z20241218#A3#BE#CFS8133578#E00005#F073416#G00228202950#H595560106973525236#I1934959901#KFS7324558#L1730616749#M466#N -#OY#PN#WALPR5CBC#XPANCIERA#YPENNY#ZPP!F5000!MCT"
const REQUEST_E1 = "880151D0E1          1011730433129     20210531          AM04C2D0ELIGCOBCCELIGIBILITYCDCOOLC61AM01C419420501C51CAELIGIBILITYCBCOOL4X00"
const REQUEST_S1 = "880151D0S1TEST      1011234567893     20231214          AM04C2TESTC1TESTHEALTHC61AM01CX99CYTESTC420231109C51CATESTCBPWLCM41 GREEN STREETCNTESTCOMACP01844CQ34787077374X01AM07EM2D26007766D300D80DF00DJ0DT0DI12E200U701AM11D90000000{DQ0000000{DU0000000{DN00AMXX&GWX&ID&LN&N1538201264&U        &V        &WM#AM#E     #KBC3234666#M        #N02!F          !H  "
const REQUEST_S2 = "880151D0S29999      1011234567893     20231212          AM04C2TESTCDTESTC1COSC61G2YAM014X03AM07EM2D29152724D300D80DF00DJ0DT0E200U705AMXX&G+C&IV&LY&N1558764159&U        &V        &W3#A3#E     #KFP4868153#M        #N04!F          !H  "
const REQUEST_P1 = "880151D0P1          1010424919        20231212          AM04C2TESTCCTESTCDTESTC1RX8909C301C60G2YAM01CX01CY431042504C419530421C52CATESTCBTESTCM524 TEST RDCNHOT SPRINGSCOSCCP819018213CQ8007012333C7012C14X03AM07EM1D29127715E103D765862017001E70000007000D304D57D61D80DE20231115DF99DJ3DT328EAU705AM03EZ12DBBB6746804DRBHARANYPM50162533342E12DLBB67468044EBHARANY2JNEERAJ2K180 TEST PARK PL2MGREER2NSC2P719018067AM11D96DDC55ADQ0000000{DU61EDN01AMXX&BTEST&CA&DG&FEC880&GMG&H00&JEC880TPNR&K004440&LN&MN&N1558764159&U31&V615&WC&XN&YN5XXXO&Z20241211#AC#BE#CBB6746804#E00007#FE-3199#G65862017001#H3077fdb2e0344f78b2b05e94ac9b7021#I5430707#KFP4868153#L1497751390#M64#N01#OY#PN#Q431042504#WATEN8C2F#XWEISS#YKAREN#ZKW!F7000!MAR!P610 TEST ROAD!QHOT SPRINGS!RSC!S819018213"

const RESPONSE_B1 = "D0B11A011234567893     20210118AM20F4QS/1 POWERLINE D.0 TESTING TRANSMISSION LEVEL MESSAGE TEXT GOES HERE.  THE MESSAGE CAN BE UP TO 200 BYTES LONG AND SHOULD CONTAIN INFORMATION ABOUT THE TRANSMISSION OF THE CLAIM, NOT JUST ABOUT THE RXAM21ANPF31234567891234567895F36F0026F0046F012UF2UH01FQRX LEVEL MESSAGE TEXT FIRST FQ FIELDUH02FQRX LEVEL MESSAGE TEXT SECOND FQ FIELD7F038F8008457558AM22EM1D299999999F1AP03AR17236056901AS52EAUPREF PROD DESCRIPTIONAM23F5100{F6557{F7100{AV1J21J301J4150{F9707{FM1FN20{FI80{MW20{EQ20{"
const RESPONSE_B2 = "D0B21A011679198717     20231219AM21ANAAM22EM1D27159262*AD0"
const RESPONSE_B3 = "D0B31R011851480545     20231124AM20F4TRANSACTION CODE/TYPE NOT SUPPORTEDAM21ANRFA1FB1SUF03UH01FQRTR-GATEWAY: TRANSACTION CODE (103-A3) IUG+UH02FQS INVALID. SUBMIT CLAIM AND REVERSAL SEPUG+UH03FQARATELY*AD0**5556615020101872700000100010000"
const RESPONSE_E1 = "D0E11A011952758781     20231212AM20F4LISLVL:3;LISEFF:20230101;LISTERM:20231231;PLAN:PDP ;MBI:5TJ2K60AM25;ED:20170520;QMB:N;QED:        ;QTERM:        ;AM27UR1UQYU1S5921U6376AM29CATESTCBPOWERLINEC419570507AM21ANAAM28NT15C016C037C610097MH9999NU0212345911MJPDPINDUB8778896510UW1UX20211001UY20261231"
const RESPONSE_S2 = "D0S21R011851005763     20231120AM20F4TRANSACTION CODE/TYPE NOT SUPPORTEDAM21ANRFA1FB1SUF02UH01FQRTR-GATEWAY: TRANSACTION CODE (103-A3) IUG+UH02FQS INVALID.*AD0"

type parseTest struct {
	rawData string
}

var dynamicTests = []parseTest{
	{rawData: REQUEST_B1},
	{rawData: REQUEST_B2},
	{rawData: REQUEST_B3},
	{rawData: RESPONSE_B1},
	{rawData: RESPONSE_B2},
	{rawData: RESPONSE_B3},
}

var requestTests = []parseTest{
	{rawData: REQUEST_B1},
	{rawData: REQUEST_B2},
	{rawData: REQUEST_B3},
}

var responseTests = []parseTest{
	{rawData: RESPONSE_B1},
	{rawData: RESPONSE_B2},
	{rawData: RESPONSE_B3},
}

var billingRequestTests = []parseTest{
	{rawData: REQUEST_B1},
}

var reversalRequestTests = []parseTest{
	{rawData: REQUEST_B2},
}

var rebillRequestTests = []parseTest{
	{rawData: REQUEST_B3},
}

var eligibilityRequestTests = []parseTest{
	{rawData: REQUEST_E1},
}

var serviceBillingRequestTests = []parseTest{
	{rawData: REQUEST_S1},
}

var serviceReversalRequestTests = []parseTest{
	{rawData: REQUEST_S2},
}

var priorAuthRequestTests = []parseTest{
	{rawData: REQUEST_P1},
}

var billingResponseTests = []parseTest{
	{rawData: RESPONSE_B1},
}

var reversalResponseTests = []parseTest{
	{rawData: RESPONSE_B2},
}

var rebillResponseTests = []parseTest{
	{rawData: RESPONSE_B3},
}

var serviceReversalResponseTests = []parseTest{
	{rawData: RESPONSE_S2},
}

var eligibilityResponseTests = []parseTest{
	{rawData: RESPONSE_E1},
}

func Test_CanParseDynamic(t *testing.T) {
	for _, test := range dynamicTests {

		i, err := ParseDynamic(test.rawData)

		if err != nil {
			t.Error(err)
			break
		}

		if i == nil {
			t.Error("Result object is null")
		}

		switch i.(type) {
		case request.Billing:
			item, ok := i.(request.Billing)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case request.Reversal:
			item, ok := i.(request.Reversal)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case request.Rebill:
			item, ok := i.(request.Rebill)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case response.Billing:
			item, ok := i.(response.Billing)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case response.Reversal:
			item, ok := i.(response.Reversal)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case response.Rebill:
			item, ok := i.(response.Rebill)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}

		default:
			t.Errorf("unknown type")
		}
	}
}

func Test_CanParseDynamicRequest(t *testing.T) {
	for _, test := range requestTests {

		i, err := ParseRequestDynamic(test.rawData)

		if err != nil {
			t.Error(err)
			break
		}

		if i == nil {
			t.Error("Result object is null")
		}

		switch i.(type) {
		case request.Billing:
			item, ok := i.(request.Billing)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case request.Reversal:
			item, ok := i.(request.Reversal)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case request.Rebill:
			item, ok := i.(request.Rebill)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}

		default:
			t.Errorf("unknown type")
		}
	}
}

func Test_CanParseDynamicResponse(t *testing.T) {
	for _, test := range responseTests {

		i, err := ParseResponseDynamic(test.rawData)

		if err != nil {
			t.Error(err)
			break
		}

		if i == nil {
			t.Error("Result object is null")
		}

		switch i.(type) {
		case response.Billing:
			item, ok := i.(response.Billing)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case response.Reversal:
			item, ok := i.(response.Reversal)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}
		case response.Rebill:
			item, ok := i.(response.Rebill)
			if !ok {
				t.Errorf("unable to cast to type: %T", i)
			}
			if len(item.Claims) <= 0 {
				t.Errorf("Group count mismatch. Wanted: >=1   Got: %v", len(item.Claims))
			}

		default:
			t.Errorf("unknown type")
		}
	}
}

func Test_CanParseBillingRequest(t *testing.T) {
	for _, test := range billingRequestTests {

		obj := request.Billing{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseReversalRequest(t *testing.T) {
	for _, test := range reversalRequestTests {

		obj := request.Reversal{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseRebillRequest(t *testing.T) {
	for _, test := range rebillRequestTests {

		obj := request.Rebill{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseServiceBillingRequest(t *testing.T) {
	for _, test := range serviceBillingRequestTests {

		obj := request.ServiceBilling{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) == 0 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: 0")
		}
	}
}

func Test_CanParseServiceReversalRequest(t *testing.T) {
	for _, test := range serviceReversalRequestTests {

		obj := request.ServiceReversal{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParsePriorAuthRequest(t *testing.T) {
	for _, test := range priorAuthRequestTests {

		obj := request.PriorAuthorization{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Authorizations) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Authorizations))
		}
	}
}

func Test_CanParseEligibilityRequest(t *testing.T) {
	for _, test := range eligibilityRequestTests {

		obj := request.Eligibility{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "ELIGIBILITY" {
			t.Errorf("Patient name  mismatch. Wanted: ELIGIBILITY   Got: %v", obj.Patient.FirstName)
		}
	}
}

func Test_CanParseBillingResponse(t *testing.T) {
	for _, test := range billingResponseTests {

		obj := response.Billing{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseRebillResponse(t *testing.T) {
	for _, test := range rebillResponseTests {

		obj := response.Rebill{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseReversalResponse(t *testing.T) {
	for _, test := range reversalResponseTests {

		obj := response.Reversal{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseServiceReversalResponse(t *testing.T) {
	for _, test := range serviceReversalResponseTests {

		obj := response.ServiceReversal{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}
	}
}

func Test_CanParseEligibilityResponse(t *testing.T) {
	for _, test := range eligibilityResponseTests {

		obj := response.Eligibility{}

		err := ParseRawClaim(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "TEST" {
			t.Errorf("Patient name  mismatch. Wanted: TEST   Got: %v", obj.Patient.FirstName)
		}
	}
}
