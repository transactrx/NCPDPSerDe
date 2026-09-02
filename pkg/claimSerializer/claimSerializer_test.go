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

// A blank-padded field at the very end of the message (e.g. an all-space DUR
// co-agent ID) is field data sent by the pharmacy: it must survive a
// deserialize/serialize round trip with its padding intact instead of being
// emitted as a bare field ID with no value.
func Test_TrailingBlankPaddedFieldRoundTrips(t *testing.T) {
	fs := string(ncpdp.FIELD)
	paddedTail := fs + "J9  " + fs + "H6                   "

	cut := strings.Index(REQUEST_B1, fs+"J9")
	if cut == -1 {
		t.Fatal("sample claim is missing the J9 field")
	}

	i, err := claimdeserializer.Deserialize(REQUEST_B1[:cut] + paddedTail)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing.  Got: %T", i)
	}

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(serialized, paddedTail) {
		t.Errorf("trailing blank-padded fields not preserved.\nWanted substring: %q\nGot: %q", paddedTail, serialized)
	}
}

// Same guarantee on the response path: a blank-padded value on the last field
// of a response must round-trip unchanged.
func Test_ResponseTrailingBlankPaddedFieldRoundTrips(t *testing.T) {
	fs := string(ncpdp.FIELD)
	paddedTail := fs + "8F8008457558   "

	cut := strings.Index(RESPONSE_B1, fs+"7F")
	if cut == -1 {
		t.Fatal("sample response is missing the 7F field")
	}

	i, err := claimdeserializer.Deserialize(RESPONSE_B1[:cut] + fs + "7F03" + paddedTail)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := i.(response.Billing)
	if !ok {
		t.Fatalf("expected response.Billing.  Got: %T", i)
	}

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(serialized, paddedTail) {
		t.Errorf("trailing blank-padded field not preserved.\nWanted substring: %q\nGot: %q", paddedTail, serialized)
	}
}

// ETX framing bytes after a blank-padded field are stripped without taking the
// padding with them (padded value + ETX + ETX is how these claims arrive off
// the wire).
func Test_TrailingPaddingBeforeEtxRoundTrips(t *testing.T) {
	fs := string(ncpdp.FIELD)
	etx := string(ncpdp.ETX)
	paddedTail := fs + "J9  " + fs + "H6                   "

	cut := strings.Index(REQUEST_B1, fs+"J9")
	if cut == -1 {
		t.Fatal("sample claim is missing the J9 field")
	}

	i, err := claimdeserializer.Deserialize(REQUEST_B1[:cut] + paddedTail + etx + etx)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing.  Got: %T", i)
	}

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(serialized, paddedTail) {
		t.Errorf("padding before ETX not preserved.\nWanted substring: %q\nGot: %q", paddedTail, serialized)
	}

	if strings.Contains(serialized, etx) {
		t.Errorf("ETX framing bytes leaked into serialized output: %q", serialized)
	}
}
