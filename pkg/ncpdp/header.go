package ncpdp

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

const (
	layoutTag   = "layout"
	layoutF6Tag = "layoutF6"
)

// RawValue carries the header bytes exactly as received so they can be
// re-emitted unchanged on serialization (it survives the JSON round trip for
// that reason); Size is internal bookkeeping and is excluded from JSON.
type NcpdpHeader[V RequestHeader | ResponseHeader | FinancialRequestHeader | FinancialResponseHeader] struct {
	RawValue string `json:",omitempty"`
	Size     int    `json:"-"`
	Value    V
}

// RequestHeader holds the parsed request header fields for both D0 and F6.
// In F6 the field 101-A1 is called IIN; this code reuses Bin to hold either value.
// The code tag carries each field's NCPDP field code (e.g. 101-A1) for
// metadata export; the fixed-width wire format never contains the codes.
type RequestHeader struct {
	Bin                           string `layout:"start=0,length=6,order=1"   layoutF6:"start=4,length=8,order=3"   code:"A1"`
	Version                       string `layout:"start=6,length=2,order=2"   layoutF6:"start=0,length=2,order=1"   code:"A2"`
	TransactionCode               string `layout:"start=8,length=2,order=3"   layoutF6:"start=2,length=2,order=2"   code:"A3"`
	Pcn                           string `layout:"start=10,length=10,order=4" layoutF6:"start=12,length=10,order=4" code:"A4"`
	RecordCount                   int    `layout:"start=20,length=1,order=5"  layoutF6:"start=22,length=1,order=5"  code:"A9"`
	ServiceProviderIdQualifier    string `layout:"start=21,length=2,order=6"  layoutF6:"start=23,length=2,order=6"  code:"B2"`
	ServiceProviderId             string `layout:"start=23,length=15,order=7" layoutF6:"start=25,length=15,order=7" code:"B1"`
	DateOfService                 string `layout:"start=38,length=8,order=8"  layoutF6:"start=40,length=8,order=8"  code:"D1"`
	SoftwareVendorCertificationId string `layout:"start=46,length=10,order=9" layoutF6:"start=48,length=10,order=9" code:"AK"`
}

type ResponseHeader struct {
	Version                    string `layout:"start=0,length=2,order=1"  code:"A2"`
	TransactionCode            string `layout:"start=2,length=2,order=2"  code:"A3"`
	RecordCount                int    `layout:"start=4,length=1,order=3"  code:"A9"`
	Status                     string `layout:"start=5,length=1,order=4"  code:"F1"`
	ServiceProviderIdQualifier string `layout:"start=6,length=2,order=5"  code:"B2"`
	ServiceProviderId          string `layout:"start=8,length=15,order=6" code:"B1"`
	DateOfService              string `layout:"start=23,length=8,order=7" code:"D1"`
}

// The accumulator fields carry the FIR standard's codes: Accumulator Year is
// 650-S1 and Transaction Identifier is 651-S2 (FIR Implementation Guide v1.3,
// "Financial Information Reporting Transaction Header Segment").
type FinancialRequestHeader struct {
	Bin                           string `layout:"start=0,length=6,order=1"   code:"A1"`
	Version                       string `layout:"start=6,length=2,order=2"   code:"A2"`
	TransactionCode               string `layout:"start=8,length=2,order=3"   code:"A3"`
	Pcn                           string `layout:"start=10,length=10,order=4" code:"A4"`
	RecordCount                   int    `layout:"start=20,length=1,order=5"  code:"A9"`
	AccumulatorYear               string `layout:"start=21,length=4,order=6"  code:"S1"`
	TransactionId                 string `layout:"start=25,length=21,order=7" code:"S2"`
	SoftwareVendorCertificationId string `layout:"start=46,length=10,order=8" code:"AK"`
}

type FinancialResponseHeader struct {
	Version         string `layout:"start=0,length=2,order=1"   code:"A2"`
	TransactionCode string `layout:"start=2,length=2,order=2"   code:"A3"`
	RecordCount     int    `layout:"start=4,length=1,order=3"   code:"A9"`
	Status          string `layout:"start=5,length=1,order=4"   code:"F1"`
	AccumulatorYear string `layout:"start=6,length=4,order=5"   code:"S1"`
	TransactionId   string `layout:"start=10,length=21,order=6" code:"S2"`
}

type Layout struct {
	Start  int
	Length int
	Order  int
}

type fieldLayout struct {
	field  reflect.StructField
	layout Layout
}

// Build header
func (h *NcpdpHeader[V]) BuildNcpdpHeader() error {
	if h == nil {
		return fmt.Errorf("header is null")
	}

	headerType := reflect.TypeOf(h.Value)
	objectRef := reflect.ValueOf(h.Value)

	tagName := resolveLayoutTagName(getStructVersion(objectRef))

	// Compile list of fields by order
	mappedFields := make(map[int]fieldLayout)

	for i := 0; i < headerType.NumField(); i++ {
		field := headerType.Field(i)

		tag := resolveFieldTag(field, tagName)
		if tag != Empty {
			layout, err := getLayoutFromTag(tag)
			if err != nil {
				return err
			}
			mappedFields[layout.Order] = fieldLayout{layout: layout, field: field}
		}
	}

	// Sort by key (order)
	keys := make([]int, 0, len(mappedFields))
	for k := range mappedFields {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	// Iterate keys and concatenate field data
	sb := strings.Builder{}

	for _, key := range keys {
		fl := mappedFields[key]

		prop := objectRef.FieldByName(fl.field.Name)
		val := fmt.Sprintf("%v", prop)
		sb.WriteString(stringutils.RightPadExact(val, ' ', fl.layout.Length))
	}

	built := sb.String()

	// Keep the header bytes exactly as received while they still parse to the
	// same field values: fixed-width fields a sender padded unusually (e.g. a
	// leading-space PCN) must round-trip unchanged. Rebuild normalized only
	// when a consumer actually changed a field.
	if h.RawValue != Empty && h.RawValue != built {
		reparsed := NcpdpHeader[V]{RawValue: h.RawValue}
		if err := reparsed.ParseNcpdpHeader(); err == nil && reparsed.Value == h.Value && reparsed.Size == len(built) {
			h.RawValue = reparsed.RawValue
			return nil
		}
	}

	h.RawValue = built

	return nil
}

// Parse NCPDP header.
func (h *NcpdpHeader[V]) ParseNcpdpHeader() error {
	if h == nil {
		return fmt.Errorf("header is null")
	}

	if strings.TrimSpace(h.RawValue) == Empty {
		return fmt.Errorf("NCPDP data is empty")
	}

	item := new(V)
	size := 0

	headerType := reflect.TypeOf(item)
	objectRef := reflect.ValueOf(item)

	tagName := resolveLayoutTagName(detectVersionFromRaw(h.RawValue))

	//Set fields with layout tags
	for i := 0; i < headerType.Elem().NumField(); i++ {
		field := headerType.Elem().Field(i)
		prop := objectRef.Elem().FieldByName(field.Name)

		if prop.CanSet() {
			tag := resolveFieldTag(field, tagName)

			if tag != Empty {
				layout, err := getLayoutFromTag(tag)
				if err != nil {
					return err
				}

				endIndex := layout.Start + layout.Length

				if endIndex > size {
					size = endIndex
				}

				dataLen := len(h.RawValue)

				if dataLen >= layout.Start && dataLen >= endIndex {
					value := strings.TrimSpace(h.RawValue[layout.Start:endIndex])

					switch field.Type.String() {
					case "string":
						prop.Set(reflect.ValueOf(value))
					case "int":
						iValue, err := strconv.Atoi(value)
						if err == nil {
							prop.Set(reflect.ValueOf(iValue))
						}
					}
				}
			}
		}
	}

	if size > len(h.RawValue) {
		return fmt.Errorf("NCPDP header data length %d is shorter than expected header size %d", len(h.RawValue), size)
	}

	// Separator bytes can never occur inside the fixed-width header: finding
	// one means the header is short or garbled and every field after it is
	// misaligned. Without this check the separator is swallowed into the
	// preserved header bytes and re-emitted, producing output the parser
	// itself cannot route.
	if idx := stringutils.IndexOfAny(h.RawValue[:size], 0, []byte{FIELD, SEGMENT, GROUP}); idx != -1 {
		return fmt.Errorf("malformed NCPDP header: separator byte at offset %d inside the %d-byte fixed-width header", idx, size)
	}

	// Fixed-width parsing cannot detect a malformed header by content — a
	// short, long, or multi-byte-polluted header (real traffic: 46/55/154-byte
	// headers, non-breaking spaces in the PCN) simply shifts every following
	// field — but any such header breaks the alignment of the separator that
	// must follow it. Refuse loudly instead of silently garbling the fields.
	if len(h.RawValue) > size {
		if next := h.RawValue[size]; next != SEGMENT && next != GROUP {
			return fmt.Errorf("malformed NCPDP header: expected a segment or group separator after the %d header bytes, found %q", size, next)
		}
	}

	h.Value = *item
	h.Size = size
	h.RawValue = h.RawValue[:size]

	return nil
}

func getLayoutFromTag(tag string) (Layout, error) {
	layout := Layout{}

	_, err := fmt.Sscanf(tag, "start=%d,length=%d,order=%d", &layout.Start, &layout.Length, &layout.Order)
	if err != nil {
		return layout, err
	}

	return layout, nil
}

// resolveLayoutTagName returns the struct tag name to use for the given NCPDP
// version. Each version with its own header layout needs a case here (and a
// matching struct tag); everything else uses the default D0-era layout.
func resolveLayoutTagName(version string) string {
	if version == F6 {
		return layoutF6Tag
	}
	return layoutTag
}

// resolveFieldTag returns the tag value for the requested tag name, falling back to the default layout tag if absent.
func resolveFieldTag(field reflect.StructField, tagName string) string {
	if tagName != layoutTag {
		if tag := field.Tag.Get(tagName); tag != Empty {
			return tag
		}
	}
	return field.Tag.Get(layoutTag)
}

// detectVersionFromRaw inspects the first bytes of the raw header to determine the NCPDP version.
// Returns the known version code found there (e.g., "F6") or empty string if it cannot be determined.
func detectVersionFromRaw(raw string) string {
	if len(raw) >= 2 && KnownVersion(raw[:2]) {
		return raw[:2]
	}
	return Empty
}

// getStructVersion reads the Version field from a header struct value if present.
func getStructVersion(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return Empty
	}
	vf := v.FieldByName("Version")
	if !vf.IsValid() || vf.Kind() != reflect.String {
		return Empty
	}
	return vf.String()
}
