package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
)

const REQUEST_B1 = "880151D0B1          1011234567893     20231219          AM04C2POLICYNUMBERTHATISLOCCJOHNCDDOEC1D0PAIDDURC61AM01CX99CYVERIC419341231C51CAJOHNCBDOECM9876 TESTING LANECNSPARTANBURGCOSCCP293011234CQ8642538600AM07EM1D26000001E103D700172240780E70000001000D300D530D61D80DE20231219DF00DJ1NX2DK20DK21C8128EAEVPANUMU701AM02EY05E9               AM03EZ01DB1234587693DRERVINGPM86458212342E01DL12345876934EERVING2JJULIUS DUNKI2K15 SLAM DUNK LANE2MPHILADELPHIA2NPA2P123456789AM11D9IDC40{DQ83DDU40IDN01AM13VE2WE02DOM06.9          WE02DOZ79.899        XE1ZE20231219H11358H2ABH3CDH4EFGHIJKLMNOPXE2ZE20231218H10835H2QRH3STH4UVXYZAM06DY20230203CFREDSAIL TECHNOLOGIES INC.CG                              CH                    CI  CJ               DZ0790001976AM054C15C016C037C015581    E820231219HB1HC07DV0000000{5E026E70 6EA5 AM142Q0112V200709152U12S42R62T200709152Z114B1A4KJ29204B1B4J404B1C4J14B44KY4B5A4K14B5B4K34B84KHEART INSTITUTE4B94KHEARTSVILLE4B104KMO4B114G200709114B124KNAM158CVPC283QVISTA PACIFIC CTR3U8888 PACIFIC AVE5JJURUPA VALLEY3VCA6D92509AM16BMTBU WITH BOE E0570 DOP12/2018AM087E1E4TDE5M0E61G8E117E2E4HDE5M0E61G8E117E3E4DDE5M0E61G8E11J9QFH6COAGENT-IDAM09KE01ME321034NE100{AM10EF03EG2EC03RE03TE51927177800        ED0000000003EE0000005AUE01RE03TE38779016308        ED0000000030EE0000036HUE01RE03TE62991277601        ED0000030000EE0000380HUE01AMXX&BDICHD&CB&FEA590&G[D&H00&ID&JEA590RCP&K145023&LN&MN&N1992211627&U410&V1013&W3&Y9GYB9&Z20241218#A3#BE#CFS8133578#E00005#F073416#G00228202950#H595560106973525236#I1934959901#KFS7324558#L1730616749#M466#N -#OY#PN#WALPR5CBC#XPANCIERA#YPENNY#ZPP!F5000!MCT"
const REQUEST_B2 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B3 = "880151D0B3PCN       1011234567893     20231219cert      AM04C2POLICYNUMBERTHATISLOCCJOHNCDDOEC1D0PAIDDURC61AM01CX99CYVERIC419341231C51CAJOHNCBDOECM9876 TESTING LANECNSPARTANBURGCOSCCP293011234CQ8642538600AM07EM1D26000001E103D700172240780E70000001000D300D530D61D80DE20231219DF00DJ1NX2DK20DK21C8128EAEVPANUMU701AM02EY05E9               AM03EZ01DB1234587693DRERVINGPM86458212342E01DL12345876934EERVING2JJULIUS DUNKI2K15 SLAM DUNK LANE2MPHILADELPHIA2NPA2P123456789AM11D9IDC40{DQ83DDU40IDN01AM13VE2WE02DOM06.9          WE02DOZ79.899        XE1ZE20231219H11358H2ABH3CDH4EFGHIJKLMNOPXE2ZE20231218H10835H2QRH3STH4UVXYZAM06DY20230203CFREDSAIL TECHNOLOGIES INC.CG                              CH                    CI  CJ               DZ0790001976AM054C15C016C037C015581    E820231219HB1HC07DV0000000{5E026E70 6EA5 AM142Q0112V200709152U12S42R62T200709152Z114B1A4KJ29204B1B4J404B1C4J14B44KY4B5A4K14B5B4K34B84KHEART INSTITUTE4B94KHEARTSVILLE4B104KMO4B114G200709114B124KNAM158CVPC283QVISTA PACIFIC CTR3U8888 PACIFIC AVE5JJURUPA VALLEY3VCA6D92509AM16BMTBU WITH BOE E0570 DOP12/2018AM087E1E4TDE5M0E61G8E117E2E4HDE5M0E61G8E117E3E4DDE5M0E61G8E11J9QFH6COAGENT-IDAM09KE01ME321034NE100{AM10EF03EG2EC03RE03TE51927177800        ED0000000003EE0000005AUE01RE03TE38779016308        ED0000000030EE0000036HUE01RE03TE62991277601        ED0000030000EE0000380HUE01AMXX&BDICHD&CB&FEA590&G[D&H00&ID&JEA590RCP&K145023&LN&MN&N1992211627&U410&V1013&W3&Y9GYB9&Z20241218#A3#BE#CFS8133578#E00005#F073416#G00228202950#H595560106973525236#I1934959901#KFS7324558#L1730616749#M466#N -#OY#PN#WALPR5CBC#XPANCIERA#YPENNY#ZPP!F5000!MCT"
const REQUEST_E1 = "880151D0E1          1011730433129     20210531          AM04C2D0ELIGCOBCCELIGIBILITYCDCOOLC61AM01C419420501C51CAELIGIBILITYCBCOOL4X00"

const REQUEST_B1_BATCH = "003858D0B1MA        4011083061303     20210310          AM04C255558662900         CCMINYA       CDSIDNEY         FO        C90C1RXINN01        C3   C61AM01HN                                                                                CX99CY                6211C419711111C52CATEST*       CBTEST**         CM1444 N 4TH ST APT Z           CNCOLUMBUS            COOHCP43201          CQ5559005838C7014X00AM07EM1D2000003671354E103D731722070010        E70000030000D308D5030D61D80DE20200729DF11DJ3ET0000030000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90000601IDC0000099IDX0000000{DQ0000701HDU0000701HDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671352E103D770069009101        E70000006000D308D5024D61D80DE20200729DF11DJ3ET0000006000C800DT0EK            28MLDI00EU00EV00000000000U701AM11D90001248GDC0000099IDX0000000{DQ0001348FDU0001348FDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671349E103D759746012110        E70000045000D308D5015D61D80DE20200729DF11DJ3ET0000045000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90000216GDC0000099IDX0000000{DQ0000316FDU0000316FDN01AM03EZ01DB1578020970     DRWALKER         2NOHAM07EM1D2000003671361E103D729300012510        E70000030000D308D5030D61D80DE20200729DF11DJ3ET0000030000C800DT0EK            28EADI00EU00EV00000000000U701AM11D90001744BDC0000099IDX0000000{DQ0001844ADU0001844ADN01AM03EZ01DB1578020970     DRWALKER         2NOH"

const REQUEST_B2_AM96 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM96AA1AB1BAD1DAA2AB2BAC2CAD2DAD2D-2AE33AE34"
const REQUEST_B2_AM97 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM97AA1AB1BAD1DAA2AB2BAC2CAD2DAE33AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B2_AM98 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM98A11A2BBB3c45S41S42S4559932AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY"
const REQUEST_B2_AM99 = "880151D0B2TESTTEST1 1011234567893     20231219          AM04C2TESTC1TESTC61AM014X04AM07EM1D21159262E103D759746017110E70000028000D301D528D61D80DE20231110DF05DJ328EAU705AM99A11A2BBB3c45S41S42S4559932AMXX&BTEST&CB&FEA590&GrC&H00&IV&JEA590SSLG&K083417&LY&MN&N1679198717&U1481&V649&W3&XN&YOVNU0&Z20241218#A3#BE#CFE4702127#E00028#F054101#G59746017110#KFS7324558#L1417228453#M1684#N -#OY#PN#WPRED5C81!F28000!MCT!P201 TEST ROAD!QCITY!ZUDEF"

const RESPONSE_B1 = "D0B11A011234567893     20210118AM20F4QS/1 POWERLINE D.0 TESTING TRANSMISSION LEVEL MESSAGE TEXT GOES HERE.  THE MESSAGE CAN BE UP TO 200 BYTES LONG AND SHOULD CONTAIN INFORMATION ABOUT THE TRANSMISSION OF THE CLAIM, NOT JUST ABOUT THE RXAM21ANPF31234567891234567895F36F0026F0046F012UF2UH01FQRX LEVEL MESSAGE TEXT FIRST FQ FIELDUH02FQRX LEVEL MESSAGE TEXT SECOND FQ FIELD7F038F8008457558AM22EM1D299999999F1AP03AR17236056901AS52EAUPREF PROD DESCRIPTIONAM23F5100{F6557{F7100{AV1J21J301J4150{F9707{FM1FN20{FI80{MW20{EQ20{"
const RESPONSE_B2 = "D0B21A011679198717     20231219AM21ANAAM22EM1D27159262*AD0"
const RESPONSE_B3 = "D0B31R011851480545     20231124AM20F4TRANSACTION CODE/TYPE NOT SUPPORTEDAM21ANRFA1FB1SUF03UH01FQRTR-GATEWAY: TRANSACTION CODE (103-A3) IUG+UH02FQS INVALID. SUBMIT CLAIM AND REVERSAL SEPUG+UH03FQARATELY*AD0**5556615020101872700000100010000"
const RESPONSE_E1 = "D0E11A011952758781     20231212AM20F4LISLVL:3;LISEFF:20230101;LISTERM:20231231;PLAN:PDP ;MBI:5TJ2K60AM25;ED:20170520;QMB:N;QED:        ;QTERM:        ;AM27UR1UQYU1S5921U6376AM29CATESTCBPOWERLINEC419570507AM21ANAAM28NT15C016C037C610097MH9999NU0212345911MJPDPINDUB8778896510UW1UX20211001UY20261231"

type serializerTest struct {
	rawData string
}

var dynamicTests = []serializerTest{
	{rawData: REQUEST_B1},
	{rawData: REQUEST_B2},
	{rawData: REQUEST_B3},
	{rawData: REQUEST_E1},
	{rawData: RESPONSE_B1},
	{rawData: RESPONSE_B2},
	{rawData: RESPONSE_B3},
	{rawData: RESPONSE_E1},
}

var extraSegmentsRequestTests = []serializerTest{
	{rawData: REQUEST_B2_AM96},
	{rawData: REQUEST_B2_AM97},
	{rawData: REQUEST_B2_AM98},
	{rawData: REQUEST_B2_AM99},
}

var billingRequestTests = []serializerTest{
	{rawData: REQUEST_B1},
	{rawData: REQUEST_B1_BATCH},
}

func Test_CanSerialize(t *testing.T) {
	for _, test := range dynamicTests {
		i, err := claimdeserializer.Deserialize(test.rawData)
		if err != nil {
			t.Error(err)
			break
		}

		var rawSerialized string
		var serErr error

		switch i.(type) {
		case request.Billing:
			item, _ := i.(request.Billing)
			rawSerialized, serErr = Serialize(&item)

		case request.Reversal:
			item, _ := i.(request.Reversal)
			rawSerialized, serErr = Serialize(&item)

		case request.Rebill:
			item, _ := i.(request.Rebill)
			rawSerialized, serErr = Serialize(&item)

		case request.Eligibility:
			item, _ := i.(request.Eligibility)
			rawSerialized, serErr = Serialize(&item)

		case response.Billing:
			item, _ := i.(response.Billing)
			rawSerialized, serErr = Serialize(&item)

		case response.Reversal:
			item, _ := i.(response.Reversal)
			rawSerialized, serErr = Serialize(&item)

		case response.Rebill:
			item, _ := i.(response.Rebill)
			rawSerialized, serErr = Serialize(&item)

		case response.Eligibility:
			item, _ := i.(response.Eligibility)
			rawSerialized, serErr = Serialize(&item)

		default:
			t.Errorf("unknown type")
		}

		if serErr != nil {
			t.Error(err)
			break
		}

		if strings.TrimSpace(rawSerialized) == "" {
			t.Errorf("empty serialization response")
		}
	}
}

func Test_CanSerializeDynamicSegments(t *testing.T) {
	for _, test := range extraSegmentsRequestTests {
		i, err := claimdeserializer.Deserialize(test.rawData)
		if err != nil {
			t.Error(err)
			break
		}

		var rawSerialized string
		var serErr error

		switch i.(type) {
		case request.Billing:
			item, _ := i.(request.Billing)
			rawSerialized, serErr = Serialize(&item)

		case request.Reversal:
			item, _ := i.(request.Reversal)
			rawSerialized, serErr = Serialize(&item)

		case request.Rebill:
			item, _ := i.(request.Rebill)
			rawSerialized, serErr = Serialize(&item)

		case request.Eligibility:
			item, _ := i.(request.Eligibility)
			rawSerialized, serErr = Serialize(&item)

		case response.Billing:
			item, _ := i.(response.Billing)
			rawSerialized, serErr = Serialize(&item)

		case response.Reversal:
			item, _ := i.(response.Reversal)
			rawSerialized, serErr = Serialize(&item)

		case response.Rebill:
			item, _ := i.(response.Rebill)
			rawSerialized, serErr = Serialize(&item)

		case response.Eligibility:
			item, _ := i.(response.Eligibility)
			rawSerialized, serErr = Serialize(&item)

		default:
			t.Errorf("unknown type")
		}

		if serErr != nil {
			t.Error(err)
			break
		}

		if strings.TrimSpace(rawSerialized) == "" {
			t.Errorf("empty serialization response")
		}
	}
}

func Test_CanUpdateField(t *testing.T) {
	for _, test := range billingRequestTests {
		i, err := claimdeserializer.Deserialize(test.rawData)
		if err != nil {
			t.Error(err)
			break
		}

		var rawSerialized string
		var serErr error
		var index int

		switch i.(type) {
		case request.Billing:
			item, _ := i.(request.Billing)
			index = len(item.Claims) - 1

			item.Claims[index].Pricing.UsualAndCustmaryCharge = reflectionutils.ToPointer(78.1)
			rawSerialized, serErr = Serialize(&item)

		default:
			t.Errorf("unknown type")
		}

		if serErr != nil {
			t.Error(err)
			break
		}

		if strings.TrimSpace(rawSerialized) == "" {
			t.Errorf("empty serialization response")
		}

		builder := strings.Builder{}
		builder.WriteByte(ncpdp.FIELD)
		builder.WriteString("DQ")

		var dqIndex int

		if index == 0 {
			dqIndex = strings.Index(rawSerialized, builder.String())
			if dqIndex < 0 {
				t.Errorf("field not found: DQ / UsualAndCustmaryCharge")
			}
		} else {
			dqIndex = strings.LastIndex(rawSerialized, builder.String())
		}

		endIndex := strings.IndexByte(rawSerialized[dqIndex+1:], ncpdp.FIELD)
		if endIndex < 0 {
			endIndex = len(rawSerialized)
		} else {
			endIndex = endIndex + dqIndex
		}

		dqRaw := rawSerialized[dqIndex : endIndex+1]
		dqStringVal := dqRaw[3:]

		if dqStringVal != "781{" {
			t.Errorf("field value mismatch. Wanted: 781{  Got: %v", dqStringVal)
		}
	}
}

// F6 request header is 58 bytes: version(2) tranCode(2) IIN(8) PCN(10) recordCount(1) SPIQ(2) SPI(15) DOS(8) vendorCertId(10)
const F6_REQUEST_B1_HEADER = "F6B100880151TEST      1011234567893     20260611SVCID     "

// Raw test constants embed invisible separator bytes, so the F6 sample is built programmatically.
// F6 (vEB+) has no group separators; the single transaction's segments
// directly follow the shared segments.
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

func Test_CanRoundTripF6BillingRequest(t *testing.T) {
	obj := request.Billing{}
	err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	if len(serialized) < 58 || serialized[:58] != F6_REQUEST_B1_HEADER {
		t.Errorf("Serialized F6 header mismatch.\nWanted prefix: %q\nGot:           %q", F6_REQUEST_B1_HEADER, serialized[:min(len(serialized), 58)])
	}

	// Re-deserialize and verify the F6 fields survive the round trip
	obj2 := request.Billing{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	if obj2.Header.Value != obj.Header.Value {
		t.Errorf("Header mismatch after round trip.\nWanted: %+v\nGot:    %+v", obj.Header.Value, obj2.Header.Value)
	}
	if obj2.Patient.MiddleName == nil || *obj2.Patient.MiddleName != "QUINCY" {
		t.Errorf("Patient middle name (F6 field 0C) mismatch after round trip. Wanted: QUINCY   Got: %v", obj2.Patient.MiddleName)
	}

	if got := strings.Count(serialized, string(ncpdp.GROUP)); got != 0 {
		t.Errorf("F6 transmissions must not contain group separators. Found: %v", got)
	}

	// Repeating patient ID: the Ids slice must survive the round trip and the
	// scalar Id/IdQualifier must not cause duplicate CX/CY field emission.
	fs := string(ncpdp.FIELD)
	if got := strings.Count(serialized, fs+"CX"); got != 2 {
		t.Errorf("Serialized CX occurrence mismatch. Wanted: 2   Got: %v", got)
	}
	if got := strings.Count(serialized, fs+"CY"); got != 2 {
		t.Errorf("Serialized CY occurrence mismatch. Wanted: 2   Got: %v", got)
	}
	if len(obj2.Patient.Ids) != 2 {
		t.Errorf("Patient ids count mismatch after round trip. Wanted: 2   Got: %v", len(obj2.Patient.Ids))
	}
	if len(obj2.Claims) != 1 {
		t.Fatalf("Group count mismatch after round trip. Wanted: 1   Got: %v", len(obj2.Claims))
	}
	intermediary := obj2.Claims[0].Intermediary
	if intermediary.IdCount == nil || *intermediary.IdCount != 1 {
		t.Errorf("Intermediary id count (F6 field 8G) mismatch after round trip. Wanted: 1   Got: %v", intermediary.IdCount)
	}
	if len(intermediary.Ids) != 1 || intermediary.Ids[0].Id == nil || *intermediary.Ids[0].Id != "INTERMEDIARY01" {
		t.Errorf("Intermediary ids (F6 segment AM19) mismatch after round trip. Got: %+v", intermediary.Ids)
	}
}

func Test_CanRoundTripF6BillingResponse(t *testing.T) {
	obj := response.Billing{}
	err := claimdeserializer.DeserializeType(buildF6BillingResponse(), &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	if got := strings.Count(serialized, string(ncpdp.GROUP)); got != 0 {
		t.Errorf("F6 transmissions must not contain group separators. Found: %v", got)
	}

	// Repeating payer ID: the Payers slice must survive the round trip and the
	// singular Payer must not cause duplicate J7/J8 field emission.
	fs := string(ncpdp.FIELD)
	if got := strings.Count(serialized, fs+"J7"); got != 2 {
		t.Errorf("Serialized J7 occurrence mismatch. Wanted: 2   Got: %v", got)
	}
	if got := strings.Count(serialized, fs+"J8"); got != 2 {
		t.Errorf("Serialized J8 occurrence mismatch. Wanted: 2   Got: %v", got)
	}

	obj2 := response.Billing{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	if len(obj2.Insurance.Payers) != 2 {
		t.Fatalf("Payers count mismatch after round trip. Wanted: 2   Got: %v", len(obj2.Insurance.Payers))
	}
	for i, want := range []struct{ qualifier, id string }{{"01", "PAYER001"}, {"02", "PAYER002"}} {
		got := obj2.Insurance.Payers[i]
		if got.Qualifier == nil || *got.Qualifier != want.qualifier || got.Id == nil || *got.Id != want.id {
			t.Errorf("Payers[%d] mismatch after round trip. Wanted: %s/%s   Got: %v/%v", i, want.qualifier, want.id, got.Qualifier, got.Id)
		}
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

func Test_CanRoundTripF6FullFieldSample(t *testing.T) {
	obj := request.Billing{}
	err := claimdeserializer.DeserializeType(buildF6FullBillingRequest(), &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	if len(serialized) < 58 || serialized[:58] != F6_REQUEST_FULL_HEADER {
		t.Errorf("Serialized F6 header mismatch.\nWanted prefix: %q\nGot:           %q", F6_REQUEST_FULL_HEADER, serialized[:min(len(serialized), 58)])
	}

	if got := strings.Count(serialized, string(ncpdp.GROUP)); got != 0 {
		t.Errorf("F6 transmissions must not contain group separators. Found: %v", got)
	}

	// Single-occurrence dual-mapped patient ID must serialize exactly once.
	fs := string(ncpdp.FIELD)
	if got := strings.Count(serialized, fs+"CX"); got != 1 {
		t.Errorf("Serialized CX occurrence mismatch. Wanted: 1   Got: %v", got)
	}
	if got := strings.Count(serialized, fs+"CY"); got != 1 {
		t.Errorf("Serialized CY occurrence mismatch. Wanted: 1   Got: %v", got)
	}
	if got := strings.Count(serialized, fs+"RL"); got != 1 {
		t.Errorf("Serialized RL occurrence mismatch. Wanted: 1   Got: %v", got)
	}

	obj2 := request.Billing{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	if obj2.Header.Value != obj.Header.Value {
		t.Errorf("Header mismatch after round trip.\nWanted: %+v\nGot:    %+v", obj.Header.Value, obj2.Header.Value)
	}

	patient := obj2.Patient
	if patient.Address.StreetLine1 == nil || *patient.Address.StreetLine1 != "1234 FIRST ADDRESS LINE" {
		t.Errorf("Patient street line 1 (F6 field 7A) mismatch after round trip. Got: %v", patient.Address.StreetLine1)
	}
	if patient.Address.StreetLine2 == nil || *patient.Address.StreetLine2 != "7890 SECOND ADDRESS LINE" {
		t.Errorf("Patient street line 2 (F6 field 7B) mismatch after round trip. Got: %v", patient.Address.StreetLine2)
	}
	if patient.Address.CountryCode == nil || *patient.Address.CountryCode != "US" {
		t.Errorf("Patient country code (F6 field 1K) mismatch after round trip. Got: %v", patient.Address.CountryCode)
	}
	if patient.Species == nil || *patient.Species != "337915000" {
		t.Errorf("Patient species (F6 field S8) mismatch after round trip. Got: %v", patient.Species)
	}
	if len(patient.Ids) != 1 || patient.Ids[0].Id == nil || *patient.Ids[0].Id != "PATIENTID" {
		t.Errorf("Patient ids mismatch after round trip. Got: %+v", patient.Ids)
	}
	if patient.IdQualifier == nil || *patient.IdQualifier != "99" || patient.Id == nil || *patient.Id != "PATIENTID" {
		t.Errorf("Patient scalar id mismatch after round trip. Got: %v/%v", patient.IdQualifier, patient.Id)
	}

	if len(obj2.Claims) != 1 {
		t.Fatalf("Group count mismatch after round trip. Wanted: 1   Got: %v", len(obj2.Claims))
	}

	pricing := obj2.Claims[0].Pricing
	for _, want := range []struct {
		name string
		got  *float64
		val  float64
	}{
		{"ingredient cost (D9)", pricing.IngredientCostSubmitted, 123456.78},
		{"dispensing fee (DC)", pricing.DispensingFeeSubmitted, 25.00},
		{"professional service fee (BE)", pricing.ProfessionalServiceFeeSubmitted, 20.00},
		{"patient paid amount (DX, empty in source)", pricing.PatientPaidAmountSubmitted, 0},
		{"incentive amount (E3)", pricing.IncentiveAmountSubmitted, 5.00},
		{"flat sales tax (HA)", pricing.FlatSalesTaxAmountSubmitted, 4.32},
		{"percentage sales tax amount (GE)", pricing.PercentageSalesTaxAmountSubmitted, 1234.57},
		{"percentage sales tax rate (HE)", pricing.PercentageSalesTaxRateSubmitted, 1.0000},
		{"usual and customary charge (DQ)", pricing.UsualAndCustmaryCharge, 134567.00},
		{"gross amount due (DU)", pricing.GrossAmountDue, 124753.22},
	} {
		if want.got == nil || *want.got != want.val {
			t.Errorf("Pricing %s mismatch after round trip. Wanted: %v   Got: %v", want.name, want.val, want.got)
		}
	}
	if len(pricing.OtherAmountClaimSubmitted) != 1 || pricing.OtherAmountClaimSubmitted[0].AmountSubmitted == nil || *pricing.OtherAmountClaimSubmitted[0].AmountSubmitted != 7.55 {
		t.Errorf("Pricing other amounts (H8/H9) mismatch after round trip. Got: %+v", pricing.OtherAmountClaimSubmitted)
	}
	if len(pricing.RegulatoryFees) != 1 || pricing.RegulatoryFees[0].TypeCode == nil || *pricing.RegulatoryFees[0].TypeCode != "AB" {
		t.Errorf("Pricing regulatory fees (RK/RL) mismatch after round trip. Got: %+v", pricing.RegulatoryFees)
	}
}

// F6 (vEB+) allows only a single transaction per transmission, and without
// group separators multiple claim groups cannot be represented.
func Test_F6SerializeRejectsMultipleClaimGroups(t *testing.T) {
	obj := request.Billing{}
	err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	obj.Claims = append(obj.Claims, obj.Claims[0])

	_, err = Serialize(&obj)
	if err == nil {
		t.Fatal("Expected error serializing an F6 claim with multiple claim groups, got nil")
	}
}

func Test_CanRoundTripF6PricingRegulatoryFees(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	rawData := F6_REQUEST_B1_HEADER +
		ss + fs + "AM04" + fs + "C2POLICY123" + fs + "C61" +
		ss + fs + "AM07" + fs + "EM1" + fs + "D26000001" + fs + "E103" + fs + "D700172240780" +
		ss + fs + "AM11" + fs + "D9100{" + fs + "RK2" + fs + "RL01" + fs + "RL02"

	obj := request.Billing{}
	err := claimdeserializer.DeserializeType(rawData, &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("Group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}

	pricing := obj.Claims[0].Pricing
	if pricing.RegulatoryFeeCount == nil || *pricing.RegulatoryFeeCount != 2 {
		t.Errorf("Regulatory fee count (F6 field RK) mismatch. Wanted: 2   Got: %v", pricing.RegulatoryFeeCount)
	}
	if len(pricing.RegulatoryFees) != 2 {
		t.Fatalf("Regulatory fees count mismatch. Wanted: 2   Got: %v", len(pricing.RegulatoryFees))
	}
	for i, want := range []string{"01", "02"} {
		got := pricing.RegulatoryFees[i].TypeCode
		if got == nil || *got != want {
			t.Errorf("Regulatory fee type code (F6 field RL) [%d] mismatch. Wanted: %s   Got: %v", i, want, got)
		}
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	if got := strings.Count(serialized, fs+"RL"); got != 2 {
		t.Errorf("Serialized RL occurrence mismatch. Wanted: 2   Got: %v", got)
	}

	obj2 := request.Billing{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	if len(obj2.Claims) != 1 || len(obj2.Claims[0].Pricing.RegulatoryFees) != 2 {
		t.Errorf("Regulatory fees not preserved through round trip. Got: %+v", obj2.Claims)
	}
}

// Test_D0DualMappedFieldsSerializeOnce guards the serializer skip logic for D0
// data: a single CX/CY occurrence populates both the backward-compatible
// scalars and the repeating Ids slice, but must serialize exactly once.
func Test_D0DualMappedFieldsSerializeOnce(t *testing.T) {
	obj := request.Billing{}
	err := claimdeserializer.DeserializeType(REQUEST_B1, &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	if obj.Patient.IdQualifier == nil || *obj.Patient.IdQualifier != "99" {
		t.Errorf("Patient id qualifier mismatch. Wanted: 99   Got: %v", obj.Patient.IdQualifier)
	}
	if obj.Patient.Id == nil || *obj.Patient.Id != "VERI" {
		t.Errorf("Patient id mismatch. Wanted: VERI   Got: %v", obj.Patient.Id)
	}
	if len(obj.Patient.Ids) != 1 {
		t.Fatalf("Patient ids count mismatch. Wanted: 1   Got: %v", len(obj.Patient.Ids))
	}
	if obj.Patient.Ids[0].Qualifier == nil || *obj.Patient.Ids[0].Qualifier != "99" ||
		obj.Patient.Ids[0].Id == nil || *obj.Patient.Ids[0].Id != "VERI" {
		t.Errorf("Patient ids[0] mismatch. Wanted: 99/VERI   Got: %v/%v", obj.Patient.Ids[0].Qualifier, obj.Patient.Ids[0].Id)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	fs := string(ncpdp.FIELD)
	if got := strings.Count(serialized, fs+"CX"); got != 1 {
		t.Errorf("Serialized CX occurrence mismatch. Wanted: 1   Got: %v", got)
	}
	if got := strings.Count(serialized, fs+"CY"); got != 1 {
		t.Errorf("Serialized CY occurrence mismatch. Wanted: 1   Got: %v", got)
	}

	obj2 := request.Billing{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	if obj2.Patient.IdQualifier == nil || *obj2.Patient.IdQualifier != "99" {
		t.Errorf("Patient id qualifier mismatch after round trip. Wanted: 99   Got: %v", obj2.Patient.IdQualifier)
	}
	if obj2.Patient.Id == nil || *obj2.Patient.Id != "VERI" {
		t.Errorf("Patient id mismatch after round trip. Wanted: VERI   Got: %v", obj2.Patient.Id)
	}
}

// Test_PreservesTrailingWhitespace verifies that trailing whitespace in field values
// is preserved through the deserialize -> serialize roundtrip.
func Test_PreservesTrailingWhitespace(t *testing.T) {
	// Deserialize an existing claim first
	obj := request.Reversal{}
	err := claimdeserializer.DeserializeType(REQUEST_B2, &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Modify the cardholder ID to have trailing whitespace
	cardholderWithSpaces := "TEST    " // TEST with 4 trailing spaces
	obj.Insurance.Cardholder.Id = &cardholderWithSpaces

	// Serialize the claim
	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Deserialize again to verify whitespace is preserved
	obj2 := request.Reversal{}
	err = claimdeserializer.DeserializeType(serialized, &obj2)
	if err != nil {
		t.Fatalf("Failed to deserialize serialized claim: %v", err)
	}

	// Verify the cardholder ID has trailing whitespace preserved
	if obj2.Insurance.Cardholder.Id == nil {
		t.Fatal("Cardholder ID is nil after roundtrip")
	}
	cardholderID := *obj2.Insurance.Cardholder.Id
	if cardholderID != cardholderWithSpaces {
		t.Errorf("Cardholder ID whitespace not preserved through roundtrip.\nWanted: %q (len=%d)\nGot:    %q (len=%d)",
			cardholderWithSpaces, len(cardholderWithSpaces), cardholderID, len(cardholderID))
	}

	// Also verify the serialized string contains the field with trailing whitespace
	expectedFieldValue := string(ncpdp.FIELD) + "C2TEST    "
	if !strings.Contains(serialized, expectedFieldValue) {
		// Find the C2 field in the serialized output for debugging
		c2Index := strings.Index(serialized, string(ncpdp.FIELD)+"C2")
		if c2Index >= 0 {
			endIndex := strings.Index(serialized[c2Index+1:], string(ncpdp.FIELD))
			if endIndex < 0 {
				endIndex = len(serialized) - c2Index - 1
			}
			actualField := serialized[c2Index : c2Index+1+endIndex]
			t.Errorf("Cardholder ID whitespace not preserved in serialized output.\nWanted field to contain: %q\nGot field: %q",
				expectedFieldValue, actualField)
		} else {
			t.Error("C2 field not found in serialized output")
		}
	}
}

// Test_PreservesTrailingWhitespaceFromRawInput verifies that trailing whitespace
// in the original raw claim is preserved through deserialization.
func Test_PreservesTrailingWhitespaceFromRawInput(t *testing.T) {
	// Use an existing claim that has a known field, deserialize it,
	// then use reflection to check the raw deserialized value contains whitespace

	// First, create a modified version of REQUEST_B2 with trailing whitespace in the group ID (C1)
	// The C1 field is the GroupId in the Insurance segment
	rawClaimWithWhitespace := strings.Replace(
		REQUEST_B2,
		string(ncpdp.FIELD)+"C1TEST"+string(ncpdp.FIELD),
		string(ncpdp.FIELD)+"C1TEST    "+string(ncpdp.FIELD), // Add trailing spaces to GroupId
		1,
	)

	// Deserialize the modified claim
	obj := request.Reversal{}
	err := claimdeserializer.DeserializeType(rawClaimWithWhitespace, &obj)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Verify the GroupId has trailing whitespace preserved
	groupId := obj.Insurance.GroupId
	if groupId == nil {
		t.Fatal("GroupId is nil")
	}
	expectedGroupId := "TEST    " // with trailing spaces
	if *groupId != expectedGroupId {
		t.Errorf("GroupId whitespace not preserved during deserialization.\nWanted: %q (len=%d)\nGot:    %q (len=%d)",
			expectedGroupId, len(expectedGroupId), *groupId, len(*groupId))
	}

	// Serialize and verify whitespace is preserved in output
	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	expectedFieldValue := string(ncpdp.FIELD) + "C1TEST    "
	if !strings.Contains(serialized, expectedFieldValue) {
		fieldMarker := string(ncpdp.FIELD) + "C1"
		fieldIndex := strings.Index(serialized, fieldMarker)
		if fieldIndex >= 0 {
			endIndex := strings.Index(serialized[fieldIndex+1:], string(ncpdp.FIELD))
			if endIndex < 0 {
				endIndex = strings.Index(serialized[fieldIndex+1:], string(ncpdp.SEGMENT))
			}
			if endIndex < 0 {
				endIndex = len(serialized) - fieldIndex - 1
			}
			actualField := serialized[fieldIndex : fieldIndex+1+endIndex]
			t.Errorf("GroupId whitespace not preserved in serialized output.\nWanted field to contain: %q\nGot field: %q",
				expectedFieldValue, actualField)
		} else {
			t.Error("C1 field not found in serialized output")
		}
	}
}
