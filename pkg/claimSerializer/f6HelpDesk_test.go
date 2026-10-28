package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

// buildF6HelpDeskResponse builds an F6 B1 response whose AM21 carries two
// help desk support types, each with its own business unit count and contact,
// and a contact extension (BD) on the first contact only — the structure the
// vF6 guide defines ("The Help Desk Business Unit Type loop will repeat
// within an occurrence of the Help Desk Support Type loop").
func buildF6HelpDeskResponse() string {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	return "F6B11A011234567893     20261026" +
		ss + fs + "AM21" + fs + "ANA" + fs + "F3AUTH123456" +
		fs + "BH02" +
		fs + "BG1" + fs + "BB01" + fs + "BA1" + fs + "BF1" + fs + "BC5551234567" + fs + "BD45678901" +
		fs + "BG2" + fs + "BB01" + fs + "BA3" + fs + "BF4" + fs + "BChelp@example.com" +
		fs + "ZR3" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001"
}

func assertHelpDeskSupportTypes(t *testing.T, status responsesegment.Status, context string) {
	t.Helper()

	if len(status.HelpDeskSupportTypes) != 2 {
		t.Fatalf("%v: support type count mismatch. Wanted: 2   Got: %v", context, len(status.HelpDeskSupportTypes))
	}

	for i, want := range []struct {
		supportType string
		unitType    string
		qualifier   string
		contact     string
		extension   *string
	}{
		{"1", "1", "1", "5551234567", strPtr("45678901")},
		{"2", "3", "4", "help@example.com", nil},
	} {
		st := status.HelpDeskSupportTypes[i]

		if st.Type == nil || *st.Type != want.supportType {
			t.Errorf("%v: support type %v mismatch. Wanted: %q   Got: %v", context, i+1, want.supportType, st.Type)
		}
		if st.BusinessUnitCount == nil || *st.BusinessUnitCount != 1 {
			t.Errorf("%v: support type %v BB count mismatch. Wanted: 1   Got: %v", context, i+1, st.BusinessUnitCount)
		}
		if len(st.BusinessUnits) != 1 {
			t.Fatalf("%v: support type %v business unit count mismatch. Wanted: 1   Got: %v", context, i+1, len(st.BusinessUnits))
		}

		bu := st.BusinessUnits[0]
		if bu.Type == nil || *bu.Type != want.unitType {
			t.Errorf("%v: support type %v unit type mismatch. Wanted: %q   Got: %v", context, i+1, want.unitType, bu.Type)
		}
		if bu.ContactInformationQualifier == nil || *bu.ContactInformationQualifier != want.qualifier {
			t.Errorf("%v: support type %v qualifier mismatch. Wanted: %q   Got: %v", context, i+1, want.qualifier, bu.ContactInformationQualifier)
		}
		if bu.ContactInformation == nil || *bu.ContactInformation != want.contact {
			t.Errorf("%v: support type %v contact mismatch. Wanted: %q   Got: %v", context, i+1, want.contact, bu.ContactInformation)
		}
		if want.extension == nil {
			if bu.ContactInformationExtension != nil {
				t.Errorf("%v: support type %v extension mismatch. Wanted: nil   Got: %q", context, i+1, *bu.ContactInformationExtension)
			}
		} else if bu.ContactInformationExtension == nil || *bu.ContactInformationExtension != *want.extension {
			t.Errorf("%v: support type %v extension mismatch. Wanted: %q   Got: %v", context, i+1, *want.extension, bu.ContactInformationExtension)
		}
	}
}

func strPtr(s string) *string { return &s }

// The F6 AM21 help desk loop nests business units within each support type.
// Both the parsed shape and the rebuilt message must keep the second support
// type's business unit count and the contact-to-support-type relationship,
// with the contact extension on its own contact.
func Test_RoundTripF6HelpDeskContacts(t *testing.T) {
	raw := buildF6HelpDeskResponse()

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj.Claims) != 1 {
		t.Fatalf("claim group count mismatch. Wanted: 1   Got: %v", len(obj.Claims))
	}
	assertHelpDeskSupportTypes(t, obj.Claims[0].Status, "parsed")

	for _, serialize := range []struct {
		name string
		run  func() (string, error)
	}{
		{"direct", func() (string, error) { return Serialize(&obj) }},
		{"via JSON", func() (string, error) { return jsonRoundTripSerialize(obj) }},
	} {
		serialized, err := serialize.run()
		if err != nil {
			t.Fatalf("%v serialization failed: %v", serialize.name, err)
		}

		assertFieldValues(t, serialized, "BH", []string{"2"})
		assertFieldValues(t, serialized, "BG", []string{"1", "2"})
		assertFieldValues(t, serialized, "BB", []string{"1", "1"})
		assertFieldValues(t, serialized, "BA", []string{"1", "3"})
		assertFieldValues(t, serialized, "BF", []string{"1", "4"})
		assertFieldValues(t, serialized, "BC", []string{"5551234567", "help@example.com"})
		assertFieldValues(t, serialized, "BD", []string{"45678901"})

		// The fields of one support type stay contiguous: BD follows its own
		// contact, before the second BG.
		fs := string(ncpdp.FIELD)
		wantOrder := fs + "BG1" + fs + "BB1" + fs + "BA1" + fs + "BF1" + fs + "BC5551234567" + fs + "BD45678901" +
			fs + "BG2" + fs + "BB1" + fs + "BA3" + fs + "BF4" + fs + "BChelp@example.com"
		if !strings.Contains(serialized, wantOrder) {
			t.Errorf("%v: help desk loop order mismatch.\nWanted substring: %q\nGot: %q", serialize.name, wantOrder, serialized)
		}

		// A rebuilt message must parse back to the same nested shape.
		reparsed := response.Billing{}
		if err := claimdeserializer.DeserializeType(serialized, &reparsed); err != nil {
			t.Fatalf("%v: failed to reparse: %v", serialize.name, err)
		}
		assertHelpDeskSupportTypes(t, reparsed.Claims[0].Status, serialize.name+" reparsed")
	}
}

// The intermediary help desk loop (AM36) nests the same way: each support
// type (KB) carries its own business unit count (G9) and contacts.
func Test_RoundTripF6IntermediaryHelpDesk(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)

	raw := "F6B11A011234567893     20261026" +
		ss + fs + "AM21" + fs + "ANA" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D26000001" +
		ss + fs + "AM36" + fs + "KC02" +
		fs + "KB1" + fs + "G901" + fs + "G81" + fs + "KA1" + fs + "JP5559876543" + fs + "JR00001234" +
		fs + "KB2" + fs + "G901" + fs + "G82" + fs + "KA4" + fs + "JPintermediary@example.com"

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	intermediary := obj.Claims[0].Intermediary
	if len(intermediary.HelpDeskSupportTypes) != 2 {
		t.Fatalf("support type count mismatch. Wanted: 2   Got: %v", len(intermediary.HelpDeskSupportTypes))
	}
	for i, wantContact := range []string{"5559876543", "intermediary@example.com"} {
		st := intermediary.HelpDeskSupportTypes[i]
		if len(st.BusinessUnits) != 1 {
			t.Fatalf("support type %v business unit count mismatch. Wanted: 1   Got: %v", i+1, len(st.BusinessUnits))
		}
		if got := st.BusinessUnits[0].ContactInformation; got == nil || *got != wantContact {
			t.Errorf("support type %v contact mismatch. Wanted: %q   Got: %v", i+1, wantContact, got)
		}
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertFieldValues(t, serialized, "KC", []string{"2"})
	assertFieldValues(t, serialized, "KB", []string{"1", "2"})
	assertFieldValues(t, serialized, "G9", []string{"1", "1"})
	assertFieldValues(t, serialized, "JP", []string{"5559876543", "intermediary@example.com"})
	assertFieldValues(t, serialized, "JR", []string{"00001234"})
}
