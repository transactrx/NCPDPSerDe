package dynamic

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_ReplaceIllegalCharactersRoundTrip(t *testing.T) {
	codes := []string{"&B", "#A", "!F", "AD", "S4", "{X", "=Y", "[D", "c4", "__"}

	for _, code := range codes {
		mangled := ReplaceIllegalCharacters(code)
		if got := restoreIllegalCharacters(mangled); got != code {
			t.Errorf("restoreIllegalCharacters(%q) mismatch. Wanted: %q   Got: %q", mangled, code, got)
		}
	}
}

func Test_FieldNameParseRoundTrip(t *testing.T) {
	for _, test := range []struct {
		code  string
		order int
	}{
		{code: "&B", order: 25},
		{code: "#A", order: 1},
		{code: "AD", order: 100},
	} {
		name := FieldName(test.code, test.order)
		code, order, ok := parseFieldName(name)
		if !ok || code != test.code || order != test.order {
			t.Errorf("parseFieldName(%q) mismatch. Wanted: (%q, %v)   Got: (%q, %v, %v)", name, test.code, test.order, code, order, ok)
		}
	}
}

// UnmarshalJSON must rebuild DynamicType and a typed Value so a DynamicStruct
// that went through a JSON round-trip can still be serialized.
func Test_DynamicStructJsonRoundTrip(t *testing.T) {
	raw := "AMXX&BTEST#A3"
	segId := "AMXX"
	code1 := "TEST"
	code2 := "3"

	original := struct {
		Raw                *string `field:"code=rawsegment"`
		Field_AM_1         *string `field:"code=AM,order=1"`
		Field_ampersandB_2 *string `field:"code=&B,order=2"`
		Field_hashA_3      *string `field:"code=#A,order=3"`
	}{
		Raw:                &raw,
		Field_AM_1:         &segId,
		Field_ampersandB_2: &code1,
		Field_hashA_3:      &code2,
	}

	var iface any = original
	ds := DynamicStruct{DynamicType: reflect.TypeOf(original), Value: &iface}

	jsonBytes, err := json.Marshal(ds)
	if err != nil {
		t.Fatal(err)
	}

	rebuilt := DynamicStruct{}
	if err := json.Unmarshal(jsonBytes, &rebuilt); err != nil {
		t.Fatal(err)
	}

	if rebuilt.DynamicType == nil {
		t.Fatal("DynamicType not reconstructed")
	}

	if rebuilt.Value == nil {
		t.Fatal("Value not reconstructed")
	}

	rebuiltVal := reflect.ValueOf(*rebuilt.Value)

	for _, want := range []struct {
		name  string
		code  string
		value string
	}{
		{name: "Raw", code: "rawsegment", value: raw},
		{name: "Field_AM_1", code: "AM", value: segId},
		{name: "Field_ampersandB_2", code: "&B", value: code1},
		{name: "Field_hashA_3", code: "#A", value: code2},
	} {
		field, ok := rebuilt.DynamicType.FieldByName(want.name)
		if !ok {
			t.Errorf("field %v missing from reconstructed type", want.name)
			continue
		}

		tag := field.Tag.Get("field")
		if tag != "code="+want.code && !hasCodeAndOrder(tag, want.code) {
			t.Errorf("field %v tag mismatch. Got: %q", want.name, tag)
		}

		fieldVal := rebuiltVal.FieldByName(want.name)
		if fieldVal.IsNil() {
			t.Errorf("field %v value is nil", want.name)
			continue
		}

		if got := fieldVal.Elem().String(); got != want.value {
			t.Errorf("field %v value mismatch. Wanted: %q   Got: %q", want.name, want.value, got)
		}
	}
}

// hasCodeAndOrder reports whether a field tag starts with the expected code
// followed by an order component (e.g. "code=&B,order=2").
func hasCodeAndOrder(tag, code string) bool {
	prefix := "code=" + code + ",order="
	return len(tag) > len(prefix) && tag[:len(prefix)] == prefix
}
