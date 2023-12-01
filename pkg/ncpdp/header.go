package ncpdp

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

const layoutTag = "layout"

type NcpdpHeader[V RequestHeader | ResponseHeader] struct {
	RawValue string
	Size     int
	Value    V
}

type RequestHeader struct {
	Bin                           string `layout:"start=0,length=6,order=1"`
	Version                       string `layout:"start=6,length=2,order=2"`
	TransactionCode               string `layout:"start=8,length=2,order=3"`
	Pcn                           string `layout:"start=10,length=10,order=4"`
	RecordCount                   int    `layout:"start=20,length=1,order=5"`
	ServiceProviderIdQualifier    string `layout:"start=21,length=2,order=6"`
	ServiceProviderId             string `layout:"start=23,length=15,order=7"`
	DateOfService                 string `layout:"start=38,length=8,order=8"`
	SoftwareVendorCertificationId string `layout:"start=46,length=10,order=9"`
}

type ResponseHeader struct {
	Version                    string `layout:"start=0,length=2,order=1"`
	TransactionCode            string `layout:"start=2,length=2,order=2"`
	RecordCount                int    `layout:"start=4,length=1,order=3"`
	Status                     string `layout:"start=5,length=1,order=4"`
	ServiceProviderIdQualifier string `layout:"start=6,length=2,order=5"`
	ServiceProviderId          string `layout:"start=8,length=15,order=6"`
	DateOfService              string `layout:"start=23,length=8,order=7"`
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
func (h *NcpdpHeader[V]) buildNcpdpHeader() error {
	headerType := reflect.TypeOf(h.Value)
	objectRef := reflect.ValueOf(h.Value)

	// Compile list of fields by order
	mappedFields := make(map[int]fieldLayout)

	for i := 0; i < headerType.NumField(); i++ {
		field := headerType.Field(i)

		tag := field.Tag.Get(layoutTag)
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

	h.RawValue = sb.String()

	return nil
}

// Parse NCPDP header.
func (h *NcpdpHeader[V]) ParseNcpdpHeader() error {
	if h.RawValue == Empty {
		return fmt.Errorf("NCPDP data is empty")
	}

	item := new(V)
	size := 0

	headerType := reflect.TypeOf(item)
	objectRef := reflect.ValueOf(item)

	//Set fields with layout tags
	for i := 0; i < headerType.Elem().NumField(); i++ {
		field := headerType.Elem().Field(i)
		prop := objectRef.Elem().FieldByName(field.Name)

		if prop.CanSet() {
			tag := field.Tag.Get(layoutTag)

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
