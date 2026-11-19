package serde

import "testing"

func Test_ParseTag(t *testing.T) {
	tests := []struct {
		name     string
		tagField string
		tagKey   string
		want     string
	}{
		{"simple value", "code=AM07", "code", "AM07"},
		{"empty value", "code=", "code", ""},
		{"key mismatch", "order=5", "code", ""},
		{"key that is a prefix of another key does not match", "decimalPlaces=2", "decimal", ""},
		{"value containing whitespace is preserved", "format=YYYY MMdd", "format", "YYYY MMdd"},
		{"no equals sign", "code", "code", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseTag(test.tagField, test.tagKey)
			if got != test.want {
				t.Errorf("parseTag(%q, %q) mismatch. Wanted: %q   Got: %q", test.tagField, test.tagKey, test.want, got)
			}
		})
	}
}

func Test_GetFieldAttribute_ParsesAllKeys(t *testing.T) {
	attr, err := GetFieldAttribute("code=KR,order=5,countfor=Payers,sinceVersion=F6")
	if err != nil {
		t.Fatal(err)
	}

	if attr.Code != "KR" {
		t.Errorf("Code mismatch. Wanted: KR   Got: %q", attr.Code)
	}
	if attr.Order != 5 {
		t.Errorf("Order mismatch. Wanted: 5   Got: %v", attr.Order)
	}
	if attr.CountFor != "Payers" {
		t.Errorf("CountFor mismatch. Wanted: Payers   Got: %q", attr.CountFor)
	}
	if attr.SinceVersion != "F6" {
		t.Errorf("SinceVersion mismatch. Wanted: F6   Got: %q", attr.SinceVersion)
	}

	// Numeric/format keys still parse alongside the above.
	attr2, err := GetFieldAttribute("code=D7,order=5,decimalPlaces=2,overpunch=true,format=YYYYMMdd")
	if err != nil {
		t.Fatal(err)
	}
	if attr2.DecimalPlaces != 2 {
		t.Errorf("DecimalPlaces mismatch. Wanted: 2   Got: %v", attr2.DecimalPlaces)
	}
	if !attr2.Overpunch {
		t.Error("Overpunch mismatch. Wanted: true   Got: false")
	}
	if attr2.Format != "20060102" {
		t.Errorf("Format mismatch. Wanted: 20060102   Got: %q", attr2.Format)
	}
	if attr2.Dynamic {
		t.Error("Dynamic mismatch. Wanted: false   Got: true")
	}

	// CountFor and SinceVersion default to empty when the tag omits them.
	if attr2.CountFor != "" {
		t.Errorf("CountFor should default empty. Got: %q", attr2.CountFor)
	}
	if attr2.SinceVersion != "" {
		t.Errorf("SinceVersion should default empty. Got: %q", attr2.SinceVersion)
	}
}

func Test_GetFieldAttribute_VersionedOverrides(t *testing.T) {
	// 448-ED: 3 implied decimals in D0, 7 since F6.
	attr, err := GetFieldAttribute("code=ED,decimalPlaces=3,decimalPlacesSince=F6:7,order=7")
	if err != nil {
		t.Fatal(err)
	}
	if attr.DecimalPlaces != 3 || attr.DecimalPlacesSinceVersion != "F6" || attr.DecimalPlacesSince != 7 {
		t.Errorf("decimalPlacesSince parse mismatch. Got: %+v", attr)
	}
	for version, want := range map[string]int{"": 3, "D0": 3, "F6": 7} {
		if got := attr.DecimalPlacesFor(version); got != want {
			t.Errorf("DecimalPlacesFor(%q) mismatch. Wanted: %v   Got: %v", version, want, got)
		}
	}

	// 483-HE: signed (overpunch) in D0, unsigned since F6.
	attr2, err := GetFieldAttribute("code=HE,decimalPlaces=4,overpunch=true,unsignedSince=F6,order=12")
	if err != nil {
		t.Fatal(err)
	}
	if attr2.UnsignedSinceVersion != "F6" {
		t.Errorf("unsignedSince parse mismatch. Got: %+v", attr2)
	}
	for version, want := range map[string]bool{"": true, "D0": true, "F6": false} {
		if got := attr2.OverpunchFor(version); got != want {
			t.Errorf("OverpunchFor(%q) mismatch. Wanted: %v   Got: %v", version, want, got)
		}
	}

	// A malformed decimalPlacesSince value (no colon) records no override.
	attr3, err := GetFieldAttribute("code=ZZ,decimalPlaces=2,decimalPlacesSince=F6,order=1")
	if err != nil {
		t.Fatal(err)
	}
	if attr3.DecimalPlacesSinceVersion != "" || attr3.DecimalPlacesFor("F6") != 2 {
		t.Errorf("malformed decimalPlacesSince should leave no override. Got: %+v", attr3)
	}

	// No override tags: resolvers return the base values for every version.
	attr4, err := GetFieldAttribute("code=D9,decimalPlaces=2,overpunch=true,order=2")
	if err != nil {
		t.Fatal(err)
	}
	if attr4.DecimalPlacesFor("F6") != 2 || !attr4.OverpunchFor("F6") {
		t.Errorf("base attributes must apply to every version. Got: %+v", attr4)
	}
}

func Test_NewVersionedFieldSettings(t *testing.T) {
	attr, err := GetFieldAttribute("code=ED,decimalPlaces=3,decimalPlacesSince=F6:7,unsignedSince=F6,overpunch=true,order=7")
	if err != nil {
		t.Fatal(err)
	}

	d0 := NewVersionedFieldSettings(attr, "D0")
	if d0.DecimalPlaces != 3 || !d0.Overpunch {
		t.Errorf("D0 settings mismatch. Got: %+v", d0)
	}

	f6 := NewVersionedFieldSettings(attr, "F6")
	if f6.DecimalPlaces != 7 || f6.Overpunch {
		t.Errorf("F6 settings mismatch. Got: %+v", f6)
	}

	// NewFieldSettings keeps its historical (oldest-version) behavior.
	base := NewFieldSettings(attr)
	if base.DecimalPlaces != 3 || !base.Overpunch {
		t.Errorf("base settings mismatch. Got: %+v", base)
	}
}

func Test_GetFieldAttribute_Dynamic(t *testing.T) {
	attr, err := GetFieldAttribute("code=dynamic")
	if err != nil {
		t.Fatal(err)
	}

	if !attr.Dynamic {
		t.Error("Dynamic mismatch. Wanted: true   Got: false")
	}
}

func Test_GetFieldAttribute_CacheReturnsSameResult(t *testing.T) {
	const tag = "code=E7,order=9,decimalPlaces=3"

	first, err := GetFieldAttribute(tag)
	if err != nil {
		t.Fatal(err)
	}

	second, err := GetFieldAttribute(tag)
	if err != nil {
		t.Fatal(err)
	}

	if first != second {
		t.Errorf("Cached attribute mismatch.\nFirst:  %+v\nSecond: %+v", first, second)
	}
}

func Test_GetSegmentAttribute_ParsesAllKeys(t *testing.T) {
	attr, err := GetSegmentAttribute("code=AM07,order=1")
	if err != nil {
		t.Fatal(err)
	}

	if attr.Code != "AM07" {
		t.Errorf("Code mismatch. Wanted: AM07   Got: %q", attr.Code)
	}
	if attr.Order != 1 {
		t.Errorf("Order mismatch. Wanted: 1   Got: %v", attr.Order)
	}
	if attr.Dynamic {
		t.Error("Dynamic mismatch. Wanted: false   Got: true")
	}
}

func Test_GetRequestType(t *testing.T) {
	RegisterTypes()

	tranType, err := GetRequestType("B1")
	if err != nil {
		t.Fatal(err)
	}
	if tranType.Name() != "Billing" {
		t.Errorf("Request type mismatch. Wanted: Billing   Got: %v", tranType.Name())
	}

	if _, err := GetRequestType("ZZ"); err == nil {
		t.Error("Expected error for unknown request transaction code, got nil")
	}
}

func Test_GetResponseType(t *testing.T) {
	RegisterTypes()

	tranType, err := GetResponseType("B2")
	if err != nil {
		t.Fatal(err)
	}
	if tranType.Name() != "Reversal" {
		t.Errorf("Response type mismatch. Wanted: Reversal   Got: %v", tranType.Name())
	}

	if _, err := GetResponseType("ZZ"); err == nil {
		t.Error("Expected error for unknown response transaction code, got nil")
	}
}

func Test_DynamicFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldCode string
		order     int
		want      string
	}{
		{"plain code", "D7", 3, "Field_D7_3"},
		{"ampersand replaced", "&B", 1, "Field_ampersandB_1"},
		{"hash replaced", "#A", 2, "Field_hashA_2"},
		{"exclamation replaced", "!F", 7, "Field_exclamationF_7"},
		{"multiple illegal characters", "{}", 4, "Field_leftbracerightbrace_4"},
		{"at sign replaced", "@A", 5, "Field_atsignA_5"},
		{"control byte hex escaped", "\x04Z", 6, "Field_u0004Z_6"},
		{"high byte hex escaped", "�B", 7, "Field_ufffdB_7"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DynamicFieldName(test.fieldCode, test.order)
			if got != test.want {
				t.Errorf("DynamicFieldName(%q, %v) mismatch. Wanted: %q   Got: %q", test.fieldCode, test.order, test.want, got)
			}
		})
	}
}

func Test_GetSegmentAttribute_CacheReturnsSameResult(t *testing.T) {
	const tag = "code=AM11,order=7"

	first, err := GetSegmentAttribute(tag)
	if err != nil {
		t.Fatal(err)
	}

	second, err := GetSegmentAttribute(tag)
	if err != nil {
		t.Fatal(err)
	}

	if first != second {
		t.Errorf("Cached attribute mismatch.\nFirst:  %+v\nSecond: %+v", first, second)
	}
}
