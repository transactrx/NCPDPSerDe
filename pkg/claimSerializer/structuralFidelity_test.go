package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// An empty or unparseable field that belongs to a repeating group's element
// type (seen in production: DUR items whose FU/FV arrive empty or as the
// invalid date 00000000) must stay attached to ITS occurrence. Segment-level
// preservation re-emitted it after the last item, silently moving it to the
// wrong DUR item.
func Test_RoundTripKeepsEmptyDurFieldsInTheirItem(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "D0B11A011234567893     20261001" +
		ss + fs + "AM20" + fs + "F4TEXT" +
		gs + ss + fs + "AM21" + fs + "ANP" + fs + "F3AUTH123" +
		ss + fs + "AM22" + fs + "EM1" + fs + "D21234567" +
		ss + fs + "AM24" +
		fs + "J61" + fs + "E4DC" + fs + "FT0" + fs + "FU00000000" + fs + "FYfirst alert" +
		fs + "J62" + fs + "E4DD" + fs + "FT1" + fs + "FU20261008" + fs + "FYsecond alert"

	obj := response.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	items := obj.Claims[0].Dur.Items
	if len(items) != 2 {
		t.Fatalf("DUR item count mismatch. Wanted: 2   Got: %v", len(items))
	}
	if len(items[0].DynamicFields) != 1 {
		t.Errorf("First item dynamic field count mismatch. Wanted: 1 (invalid FU)   Got: %v", len(items[0].DynamicFields))
	}
	if len(items[1].DynamicFields) != 0 {
		t.Errorf("Second item dynamic field spillover. Wanted: 0   Got: %v", len(items[1].DynamicFields))
	}
	if len(obj.Claims[0].Dur.DynamicFields) != 0 {
		t.Errorf("Segment-level dynamic field spillover. Wanted: 0   Got: %v", len(obj.Claims[0].Dur.DynamicFields))
	}

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

		assertFieldValues(t, serialized, "FU", []string{"00000000", "20261008"})

		// The preserved FU must re-emit inside the first item: after that
		// item's parsed fields, before the second item's counter.
		fuIndex := strings.Index(serialized, fs+"FU00000000")
		secondItemIndex := strings.Index(serialized, fs+"J62")
		if fuIndex == -1 || secondItemIndex == -1 || fuIndex > secondItemIndex {
			t.Errorf("%v: preserved FU emitted outside its item. FU index: %v   J62 index: %v", serialize.name, fuIndex, secondItemIndex)
		}
	}
}

// An empty segment ("<RS><RS>", seen in production responses) must not abort
// the rest of the part: every following segment still parses, and the stray
// separator normalizes away on rebuild.
func Test_RoundTripSkipsEmptySegmentAndKeepsRest(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" +
		ss + fs + "AM04" + fs + "C2ID12345" +
		ss + ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "E103" + fs + "D700172240780" + fs + "DJ1"

	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "JOHN" {
		t.Errorf("Patient first name mismatch. Wanted: JOHN   Got: %v", obj.Patient.FirstName)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertFieldValues(t, serialized, "CA", []string{"JOHN"})
	if strings.Contains(serialized, ss+ss) {
		t.Errorf("Rebuilt message still contains an empty segment: %q", serialized)
	}
}

// A segment whose id arrives without its leading field separator (seen in
// production: a response carrying "<RS>AM29<FS>..." with the patient's name
// and birth date) must parse instead of silently ending the part. The missing
// separator normalizes in on rebuild.
func Test_RoundTripRecoversSegmentWithoutLeadingFieldSeparator(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "D0E11A073350321        20260820" +
		ss + fs + "AM25" + fs + "FODU332263" +
		ss + "AM29" + fs + "C419590930" + fs + "CAMARTIN" + fs + "CBLAW" +
		gs + ss + fs + "AM21" + fs + "ANA" + fs + "F300000000000000"

	obj := response.Eligibility{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assertFieldValues(t, serialized, "C4", []string{"19590930"})
	assertFieldValues(t, serialized, "CA", []string{"MARTIN"})
	assertFieldValues(t, serialized, "CB", []string{"LAW"})
	if !strings.Contains(serialized, ss+fs+"AM29") {
		t.Errorf("Rebuilt message does not normalize the missing field separator: %q", serialized)
	}
}

// A non-repeating segment sent twice (seen in production: a duplicated AM22)
// must not silently collapse: the first occurrence wins the typed field and
// the duplicate is preserved as a dynamic segment, so both round-trip.
func Test_RoundTripPreservesDuplicateNonRepeatingSegment(t *testing.T) {
	fs := string(ncpdp.FIELD)
	ss := string(ncpdp.SEGMENT)
	gs := string(ncpdp.GROUP)

	raw := "999999D0B1TESTPCN   1011234567893     20261001VENDOR0001" +
		ss + fs + "AM04" + fs + "C2ID12345" +
		ss + fs + "AM01" + fs + "C419800101" + fs + "CAJOHN" +
		ss + fs + "AM01" + fs + "C419990101" + fs + "CAJANE" +
		gs + ss + fs + "AM07" + fs + "EM1" + fs + "D21234567" + fs + "E103" + fs + "D700172240780" + fs + "DJ1"

	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(raw, &obj); err != nil {
		t.Fatal(err)
	}

	// First occurrence wins the typed segment, mirroring duplicate fields.
	if obj.Patient.FirstName == nil || *obj.Patient.FirstName != "JOHN" {
		t.Errorf("Patient first name mismatch. Wanted: JOHN (first occurrence)   Got: %v", obj.Patient.FirstName)
	}

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

		assertFieldValues(t, serialized, "C4", []string{"19800101", "19990101"})
		assertFieldValues(t, serialized, "CA", []string{"JOHN", "JANE"})

		// Reparsing the rebuilt message must stay stable: the duplicate keeps
		// round-tripping, not multiplying.
		reparsed := request.Billing{}
		if err := claimdeserializer.DeserializeType(serialized, &reparsed); err != nil {
			t.Fatalf("%v: failed to reparse: %v", serialize.name, err)
		}
		again, err := Serialize(&reparsed)
		if err != nil {
			t.Fatalf("%v: reserialization failed: %v", serialize.name, err)
		}
		assertFieldValues(t, again, "C4", []string{"19800101", "19990101"})
	}
}
