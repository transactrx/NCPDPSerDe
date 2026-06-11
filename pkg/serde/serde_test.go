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
	attr, err := GetFieldAttribute("code=D7,order=5,decimalPlaces=2,overpunch=true,format=YYYYMMdd")
	if err != nil {
		t.Fatal(err)
	}

	if attr.Code != "D7" {
		t.Errorf("Code mismatch. Wanted: D7   Got: %q", attr.Code)
	}
	if attr.Order != 5 {
		t.Errorf("Order mismatch. Wanted: 5   Got: %v", attr.Order)
	}
	if attr.DecimalPlaces != 2 {
		t.Errorf("DecimalPlaces mismatch. Wanted: 2   Got: %v", attr.DecimalPlaces)
	}
	if !attr.Overpunch {
		t.Error("Overpunch mismatch. Wanted: true   Got: false")
	}
	if attr.Format != "20060102" {
		t.Errorf("Format mismatch. Wanted: 20060102   Got: %q", attr.Format)
	}
	if attr.Dynamic {
		t.Error("Dynamic mismatch. Wanted: false   Got: true")
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
