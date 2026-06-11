package claimdeserializer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
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

const REQUEST_B2_AM96 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM96AA1AB1BAD1DAA2AB2BAC2CAD2DAE33AE34"
const REQUEST_B2_AM97 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM97AA1AB1BAD1DAA2AB2BAC2CAD2DAE33AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B2_AM98 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM98A11A2BBB3c45S41S42S4559932AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B2_AM99 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY!ZUDEF"

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

var extraSegmentsRequestTests = []parseTest{
	{rawData: REQUEST_B2_AM96},
	{rawData: REQUEST_B2_AM97},
	{rawData: REQUEST_B2_AM98},
	{rawData: REQUEST_B2_AM99},
}

func Test_CanParseDynamic(t *testing.T) {
	for _, test := range dynamicTests {

		i, err := Deserialize(test.rawData)

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

		i, err := DeserializeRequest(test.rawData)

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

		i, err := DeserializeResponse(test.rawData)

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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
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

		err := DeserializeType(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "TEST" {
			t.Errorf("Patient name  mismatch. Wanted: TEST   Got: %v", obj.Patient.FirstName)
		}
	}
}

func Test_CanParseUndefinedReversalSegments(t *testing.T) {
	for _, test := range extraSegmentsRequestTests {

		obj := request.Reversal{}

		err := DeserializeType(test.rawData, &obj)
		if err != nil {
			t.Error(err)
			break
		}

		if len(obj.Claims) != 1 {
			t.Errorf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
		}

		if len(obj.Claims[0].DynamicSegments) == 0 {
			t.Errorf("Group DynamicSegments count mismatch. Wanted: >0   Got: %v", len(obj.Claims[0].DynamicSegments))
		}

		bytes, err := json.Marshal(obj)
		if err != nil {
			t.Errorf("Json error: %q", err)
		}

		fmt.Print(string(bytes))
	}
}

// F6 request header is 58 bytes: version(2) tranCode(2) IIN(8) PCN(10) recordCount(1) SPIQ(2) SPI(15) DOS(8) vendorCertId(10)
const F6_REQUEST_B1_HEADER = "F6B100880151TEST      1011234567893     20260611SVCID     "

// Raw test constants embed invisible separator bytes, so the F6 sample is
// built programmatically. F6 (vEB+) has no group separators; the single
// transaction's segments directly follow the shared segments.
func buildF6BillingRequest() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	return F6_REQUEST_B1_HEADER +
		ss + fs + "AM04" + fs + "C2POLICY123" + fs + "C1TESTGROUP" + fs + "C61" +
		ss + fs + "AM01" + fs + "RR2" + fs + "CX01" + fs + "CY111111111" + fs + "CX02" + fs + "CY222222222" + fs + "C419341231" + fs + "C51" + fs + "CAJOHN" + fs + "0CQUINCY" + fs + "CBDOE" +
		ss + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" + fs + "E70000001000" + fs + "D300" + fs + "D530" + fs + "D61" + fs + "D80" + fs + "DE20260611" + fs + "U701" +
		ss + fs + "AM19" + fs + "8G1" + fs + "8H01" + fs + "8KQQ" + fs + "8MINTERMEDIARY01"
}

// F6 response header layout is identical to D0 (31 bytes).
const F6_RESPONSE_B1_HEADER = "F6B11A011234567893     20260611"

func buildF6BillingResponse() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	return F6_RESPONSE_B1_HEADER +
		ss + fs + "AM25" + fs + "C1TESTGROUP" + fs + "KR2" + fs + "J701" + fs + "J8PAYER001" + fs + "J702" + fs + "J8PAYER002" + fs + "C2CARD123" +
		ss + fs + "AM21" + fs + "ANA" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001"
}

func Test_CanParseF6BillingRequest(t *testing.T) {
	i, err := Deserialize(buildF6BillingRequest())
	if err != nil {
		t.Fatal(err)
	}

	obj, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing, got: %T", i)
	}

	header := obj.Header.Value
	if header.Version != ncpdp.F6 {
		t.Errorf("Version mismatch. Wanted: F6   Got: %q", header.Version)
	}
	if header.Bin != "00880151" {
		t.Errorf("Bin/IIN mismatch. Wanted: 00880151   Got: %q", header.Bin)
	}
	if header.TransactionCode != "B1" {
		t.Errorf("TransactionCode mismatch. Wanted: B1   Got: %q", header.TransactionCode)
	}
	if header.Pcn != "TEST" {
		t.Errorf("Pcn mismatch. Wanted: TEST   Got: %q", header.Pcn)
	}
	if header.RecordCount != 1 {
		t.Errorf("RecordCount mismatch. Wanted: 1   Got: %v", header.RecordCount)
	}
	if obj.Header.Size != 58 {
		t.Errorf("Header size mismatch. Wanted: 58   Got: %v", obj.Header.Size)
	}

	if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "JOHN" {
		t.Errorf("Patient first name mismatch. Wanted: JOHN   Got: %v", obj.Patient.FirstName)
	}
	if obj.Patient.MiddleName == nil || *obj.Patient.MiddleName != "QUINCY" {
		t.Errorf("Patient middle name (F6 field 0C) mismatch. Wanted: QUINCY   Got: %v", obj.Patient.MiddleName)
	}

	// Repeating patient ID (F6): scalars keep one occurrence for backward
	// compatibility, the Ids slice captures every occurrence.
	if obj.Patient.IdCount == nil || *obj.Patient.IdCount != 2 {
		t.Errorf("Patient id count (F6 field RR) mismatch. Wanted: 2   Got: %v", obj.Patient.IdCount)
	}
	if obj.Patient.Id == nil || obj.Patient.IdQualifier == nil {
		t.Error("Patient scalar Id/IdQualifier not populated")
	}
	if len(obj.Patient.Ids) != 2 {
		t.Fatalf("Patient ids count mismatch. Wanted: 2   Got: %v", len(obj.Patient.Ids))
	}
	for i, want := range []struct{ qualifier, id string }{{"01", "111111111"}, {"02", "222222222"}} {
		got := obj.Patient.Ids[i]
		if got.Qualifier == nil || *got.Qualifier != want.qualifier || got.Id == nil || *got.Id != want.id {
			t.Errorf("Patient ids[%d] mismatch. Wanted: %s/%s   Got: %v/%v", i, want.qualifier, want.id, got.Qualifier, got.Id)
		}
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}

	claim := obj.Claims[0]
	intermediary := claim.Intermediary
	if intermediary.IdCount == nil || *intermediary.IdCount != 1 {
		t.Errorf("Intermediary id count (F6 field 8G) mismatch. Wanted: 1   Got: %v", intermediary.IdCount)
	}
	if len(intermediary.Ids) != 1 {
		t.Fatalf("Intermediary ids (F6 segment AM19) count mismatch. Wanted: 1   Got: %v", len(intermediary.Ids))
	}
	if intermediary.Ids[0].Qualifier == nil || *intermediary.Ids[0].Qualifier != "QQ" {
		t.Errorf("Intermediary id qualifier (F6 field 8K) mismatch. Wanted: QQ   Got: %v", intermediary.Ids[0].Qualifier)
	}
	if intermediary.Ids[0].Id == nil || *intermediary.Ids[0].Id != "INTERMEDIARY01" {
		t.Errorf("Intermediary id (F6 field 8M) mismatch. Wanted: INTERMEDIARY01   Got: %v", intermediary.Ids[0].Id)
	}
}

func Test_CanParseF6BillingResponse(t *testing.T) {
	i, err := Deserialize(buildF6BillingResponse())
	if err != nil {
		t.Fatal(err)
	}

	obj, ok := i.(response.Billing)
	if !ok {
		t.Fatalf("expected response.Billing, got: %T", i)
	}

	if obj.Header.Value.Version != ncpdp.F6 {
		t.Errorf("Version mismatch. Wanted: F6   Got: %q", obj.Header.Value.Version)
	}

	// Repeating payer ID (F6): the singular Payer keeps one occurrence for
	// backward compatibility, the Payers slice captures every occurrence.
	insurance := obj.Insurance
	if insurance.PayerIdCount == nil || *insurance.PayerIdCount != 2 {
		t.Errorf("Payer id count (F6 field KR) mismatch. Wanted: 2   Got: %v", insurance.PayerIdCount)
	}
	if insurance.Payer.Id == nil || insurance.Payer.Qualifier == nil {
		t.Error("Singular Payer not populated")
	}
	if len(insurance.Payers) != 2 {
		t.Fatalf("Payers count mismatch. Wanted: 2   Got: %v", len(insurance.Payers))
	}
	for i, want := range []struct{ qualifier, id string }{{"01", "PAYER001"}, {"02", "PAYER002"}} {
		got := insurance.Payers[i]
		if got.Qualifier == nil || *got.Qualifier != want.qualifier || got.Id == nil || *got.Id != want.id {
			t.Errorf("Payers[%d] mismatch. Wanted: %s/%s   Got: %v/%v", i, want.qualifier, want.id, got.Qualifier, got.Id)
		}
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}
}

const F6_REQUEST_FULL_HEADER = "F6B188015600PCN12345671011234567893     20260611VENDORCERT"

// Built from the raw F6 sample (RawSample_F6.txt) with CRLF line breaks
// removed. F6 has no group separator: AM11 follows the shared segments
// directly and is assigned to the single claim group by segment ID.
func buildF6FullBillingRequest() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	return F6_REQUEST_FULL_HEADER +
		ss + fs + "AM04" + fs + "C2CARDID" + fs + "CCCARDFIRST" + fs + "CDCARDLAST" + fs + "FOPLANID" + fs + "C90" + fs + "C1D0PAID" + fs + "C3001" + fs + "C61" + fs + "2AMEDIGAPID" + fs + "2BSC" + fs + "2DN" + fs + "N5MEDICAIDIDNUMBER" +
		ss + fs + "AM01" + fs + "RR1" + fs + "CX99" + fs + "CYPATIENTID" + fs + "C419850601" + fs + "C51" + fs + "CAFIRSTNAME" + fs + "CBLASTNAME" + fs + "7A1234 FIRST ADDRESS LINE" + fs + "7B7890 SECOND ADDRESS LINE" + fs + "CNROEBUCK" + fs + "COSC" + fs + "CP293764444" + fs + "1KUS" + fs + "CQ8645555555" + fs + "C701" + fs + "CZEMPLOYERID" + fs + "2C1" + fs + "HNEMAIL@DOMAIN.COM" + fs + "4X1" + fs + "S8337915000" +
		ss + fs + "AM11" + fs + "D91234567H" + fs + "DC250{" + fs + "BE200{" + fs + "DX" + fs + "E350{" + fs + "H71" + fs + "H801" + fs + "H975E" + fs + "RK1" + fs + "RLAB" + fs + "HA43B" + fs + "GE12345G" + fs + "HE10000" + fs + "JE02" + fs + "DQ1345670{" + fs + "DU1247532B" + fs + "DN01"
}

func assertString(t *testing.T, name string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Errorf("%s mismatch. Wanted: %q   Got: nil", name, want)
	} else if *got != want {
		t.Errorf("%s mismatch. Wanted: %q   Got: %q", name, want, *got)
	}
}

func assertFloat(t *testing.T, name string, got *float64, want float64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s mismatch. Wanted: %v   Got: nil", name, want)
	} else if *got != want {
		t.Errorf("%s mismatch. Wanted: %v   Got: %v", name, want, *got)
	}
}

func Test_CanParseF6FullFieldSample(t *testing.T) {
	i, err := Deserialize(buildF6FullBillingRequest())
	if err != nil {
		t.Fatal(err)
	}

	obj, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing, got: %T", i)
	}

	header := obj.Header.Value
	if header.Version != ncpdp.F6 {
		t.Errorf("Version mismatch. Wanted: F6   Got: %q", header.Version)
	}
	if header.Bin != "88015600" {
		t.Errorf("Bin/IIN mismatch. Wanted: 88015600   Got: %q", header.Bin)
	}
	if header.TransactionCode != "B1" {
		t.Errorf("TransactionCode mismatch. Wanted: B1   Got: %q", header.TransactionCode)
	}
	if header.Pcn != "PCN1234567" {
		t.Errorf("Pcn mismatch. Wanted: PCN1234567   Got: %q", header.Pcn)
	}
	if header.RecordCount != 1 {
		t.Errorf("RecordCount mismatch. Wanted: 1   Got: %v", header.RecordCount)
	}
	if header.ServiceProviderIdQualifier != "01" {
		t.Errorf("ServiceProviderIdQualifier mismatch. Wanted: 01   Got: %q", header.ServiceProviderIdQualifier)
	}
	if header.ServiceProviderId != "1234567893" {
		t.Errorf("ServiceProviderId mismatch. Wanted: 1234567893   Got: %q", header.ServiceProviderId)
	}
	if header.DateOfService != "20260611" {
		t.Errorf("DateOfService mismatch. Wanted: 20260611   Got: %q", header.DateOfService)
	}
	if header.SoftwareVendorCertificationId != "VENDORCERT" {
		t.Errorf("SoftwareVendorCertificationId mismatch. Wanted: VENDORCERT   Got: %q", header.SoftwareVendorCertificationId)
	}
	if obj.Header.Size != 58 {
		t.Errorf("Header size mismatch. Wanted: 58   Got: %v", obj.Header.Size)
	}

	insurance := obj.Insurance
	assertString(t, "Insurance cardholder id (C2)", insurance.Cardholder.Id, "CARDID")
	assertString(t, "Insurance cardholder first name (CC)", insurance.Cardholder.FirstName, "CARDFIRST")
	assertString(t, "Insurance cardholder last name (CD)", insurance.Cardholder.LastName, "CARDLAST")
	assertString(t, "Insurance plan id (FO)", insurance.PlanId, "PLANID")
	assertString(t, "Insurance eligibility clarification code (C9)", insurance.EligbilityClarificationCode, "0")
	assertString(t, "Insurance group id (C1)", insurance.GroupId, "D0PAID")
	assertString(t, "Insurance person code (C3)", insurance.PersonCode, "001")
	assertString(t, "Insurance patient relationship code (C6)", insurance.PatientRelationshipCode, "1")
	assertString(t, "Insurance medigap id (2A)", insurance.MedigapId, "MEDIGAPID")
	assertString(t, "Insurance medicaid indicator (2B)", insurance.Medicaid.Indicator, "SC")
	assertString(t, "Insurance provider accept assignment (2D)", insurance.ProviderAcceptAssignment, "N")
	assertString(t, "Insurance medicaid id (N5)", insurance.Medicaid.Id, "MEDICAIDIDNUMBER")
	if len(insurance.DynamicFields) != 0 {
		t.Errorf("Insurance dynamic field spillover. Wanted: 0   Got: %v", len(insurance.DynamicFields))
	}

	patient := obj.Patient
	if patient.IdCount == nil || *patient.IdCount != 1 {
		t.Errorf("Patient id count (RR) mismatch. Wanted: 1   Got: %v", patient.IdCount)
	}
	assertString(t, "Patient id qualifier (CX)", patient.IdQualifier, "99")
	assertString(t, "Patient id (CY)", patient.Id, "PATIENTID")
	if len(patient.Ids) != 1 {
		t.Errorf("Patient Ids count mismatch. Wanted: 1   Got: %v", len(patient.Ids))
	} else {
		assertString(t, "Patient Ids[0] qualifier", patient.Ids[0].Qualifier, "99")
		assertString(t, "Patient Ids[0] id", patient.Ids[0].Id, "PATIENTID")
	}
	if patient.BirthDate == nil || patient.BirthDate.Format("20060102") != "19850601" {
		t.Errorf("Patient birth date (C4) mismatch. Wanted: 19850601   Got: %v", patient.BirthDate)
	}
	assertString(t, "Patient gender code (C5)", patient.GenderCode, "1")
	assertString(t, "Patient first name (CA)", patient.FirstName, "FIRSTNAME")
	assertString(t, "Patient last name (CB)", patient.LastName, "LASTNAME")
	assertString(t, "Patient address street line 1 (F6 field 7A)", patient.Address.StreetLine1, "1234 FIRST ADDRESS LINE")
	assertString(t, "Patient address street line 2 (F6 field 7B)", patient.Address.StreetLine2, "7890 SECOND ADDRESS LINE")
	assertString(t, "Patient address city (CN)", patient.Address.City, "ROEBUCK")
	assertString(t, "Patient address state (CO)", patient.Address.State, "SC")
	assertString(t, "Patient address zip (CP)", patient.Address.Zip, "293764444")
	assertString(t, "Patient address country code (F6 field 1K)", patient.Address.CountryCode, "US")
	assertString(t, "Patient phone (CQ)", patient.Phone, "8645555555")
	assertString(t, "Patient place of service (C7)", patient.PlaceOfService, "01")
	assertString(t, "Patient employer id (CZ)", patient.EmployerId, "EMPLOYERID")
	assertString(t, "Patient pregnant (2C)", patient.Pregnant, "1")
	assertString(t, "Patient email (HN)", patient.Email, "EMAIL@DOMAIN.COM")
	assertString(t, "Patient residence (4X)", patient.Residence, "1")
	assertString(t, "Patient species (F6 field S8)", patient.Species, "337915000")
	if len(patient.DynamicFields) != 0 {
		t.Errorf("Patient dynamic field spillover. Wanted: 0   Got: %v", len(patient.DynamicFields))
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}

	pricing := obj.Claims[0].Pricing
	assertFloat(t, "Pricing ingredient cost (D9)", pricing.IngredientCostSubmitted, 123456.78)
	assertFloat(t, "Pricing dispensing fee (DC)", pricing.DispensingFeeSubmitted, 25.00)
	assertFloat(t, "Pricing professional service fee (BE)", pricing.ProfessionalServiceFeeSubmitted, 20.00)
	// DX is present but empty in the raw sample; the deserializer allocates the
	// pointer and leaves the zero value.
	assertFloat(t, "Pricing patient paid amount (DX, empty)", pricing.PatientPaidAmountSubmitted, 0)
	assertFloat(t, "Pricing incentive amount (E3)", pricing.IncentiveAmountSubmitted, 5.00)
	if pricing.OtherAmountClaimSubmittedCount == nil || *pricing.OtherAmountClaimSubmittedCount != 1 {
		t.Errorf("Pricing other amount count (H7) mismatch. Wanted: 1   Got: %v", pricing.OtherAmountClaimSubmittedCount)
	}
	if len(pricing.OtherAmountClaimSubmitted) != 1 {
		t.Errorf("Pricing other amounts count mismatch. Wanted: 1   Got: %v", len(pricing.OtherAmountClaimSubmitted))
	} else {
		assertString(t, "Pricing other amount qualifier (H8)", pricing.OtherAmountClaimSubmitted[0].Qualifier, "01")
		assertFloat(t, "Pricing other amount (H9)", pricing.OtherAmountClaimSubmitted[0].AmountSubmitted, 7.55)
	}
	if pricing.RegulatoryFeeCount == nil || *pricing.RegulatoryFeeCount != 1 {
		t.Errorf("Pricing regulatory fee count (F6 field RK) mismatch. Wanted: 1   Got: %v", pricing.RegulatoryFeeCount)
	}
	if len(pricing.RegulatoryFees) != 1 {
		t.Errorf("Pricing regulatory fees count mismatch. Wanted: 1   Got: %v", len(pricing.RegulatoryFees))
	} else {
		assertString(t, "Pricing regulatory fee type (F6 field RL)", pricing.RegulatoryFees[0].TypeCode, "AB")
	}
	assertFloat(t, "Pricing flat sales tax (HA)", pricing.FlatSalesTaxAmountSubmitted, 4.32)
	assertFloat(t, "Pricing percentage sales tax amount (GE)", pricing.PercentageSalesTaxAmountSubmitted, 1234.57)
	assertFloat(t, "Pricing percentage sales tax rate (HE)", pricing.PercentageSalesTaxRateSubmitted, 1.0000)
	assertString(t, "Pricing percentage sales tax basis (JE)", pricing.PercentageSalesTaxBasisSubmitted, "02")
	assertFloat(t, "Pricing usual and customary charge (DQ)", pricing.UsualAndCustmaryCharge, 134567.00)
	assertFloat(t, "Pricing gross amount due (DU)", pricing.GrossAmountDue, 124753.22)
	assertString(t, "Pricing basis of cost determination (DN)", pricing.BasisOfCostDetermination, "01")
	if len(pricing.DynamicFields) != 0 {
		t.Errorf("Pricing dynamic field spillover. Wanted: 0   Got: %v", len(pricing.DynamicFields))
	}
}

// An F6 transmission with only shared segments must parse without error and
// without inventing a claim group.
func Test_CanParseF6RequestWithoutTransactionSegments(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	rawData := F6_REQUEST_B1_HEADER +
		ss + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		ss + fs + "AM01" + fs + "CX99" + fs + "CYPATIENTID" + fs + "CAJOHN" + fs + "CBDOE"

	obj := request.Billing{}
	err := DeserializeType(rawData, &obj)
	if err != nil {
		t.Fatal(err)
	}

	if obj.Insurance.Cardholder.Id == nil || *obj.Insurance.Cardholder.Id != "POLICY123" {
		t.Errorf("Insurance cardholder id mismatch. Wanted: POLICY123   Got: %v", obj.Insurance.Cardholder.Id)
	}
	if obj.Patient.LastName == nil || *obj.Patient.LastName != "DOE" {
		t.Errorf("Patient last name mismatch. Wanted: DOE   Got: %v", obj.Patient.LastName)
	}
	if len(obj.Claims) != 0 {
		t.Errorf("Group count mismatch. Wanted: 0   Got: %v", len(obj.Claims))
	}
}

// Without group separators the claim group starts at the first segment whose
// ID belongs to the group struct. Unknown segments before that boundary stay
// in the shared DynamicSegments; unknown segments after it belong to the claim.
func Test_F6UnknownSegmentsRespectGroupBoundary(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	rawData := F6_REQUEST_B1_HEADER +
		ss + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		ss + fs + "AM98" + fs + "A11" + fs + "A2BB" +
		ss + fs + "AM07" + fs + "EM1" + fs + "D26000001" +
		ss + fs + "AM99" + fs + "B3CC" + fs + "B41"

	obj := request.Billing{}
	err := DeserializeType(rawData, &obj)
	if err != nil {
		t.Fatal(err)
	}

	if len(obj.DynamicSegments) != 1 {
		t.Errorf("Shared dynamic segment count mismatch. Wanted: 1 (AM98)   Got: %v", len(obj.DynamicSegments))
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}

	claim := obj.Claims[0]
	rxNumber := claim.Claim.PrescriptionServiceReference.Number
	if rxNumber == nil || *rxNumber != "6000001" {
		t.Errorf("Prescription number (D2) mismatch. Wanted: 6000001   Got: %v", rxNumber)
	}
	if len(claim.DynamicSegments) != 1 {
		t.Errorf("Claim dynamic segment count mismatch. Wanted: 1 (AM99)   Got: %v", len(claim.DynamicSegments))
	}
}
