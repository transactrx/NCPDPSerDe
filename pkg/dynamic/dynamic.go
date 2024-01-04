package dynamic

import (
	"fmt"
	"reflect"
	"strings"
)

type DynamicStruct struct {
	DynamicType reflect.Type `json:"-"`
	Value       *any
}

type DynamicFieldData struct {
	Name    string
	Type    reflect.Type
	TagName string
	Tags    map[string]string
}

type DynamicFieldList []DynamicFieldData

// Convert to struct field.
func (d DynamicFieldData) ToStructField() reflect.StructField {
	tags := []string{}

	for k, v := range d.Tags {
		tags = append(tags, fmt.Sprintf("%v=%v", k, v))
	}

	sTag := ""
	if len(tags) > 0 {
		sTag = fmt.Sprintf(`%v:"%v"`, d.TagName, strings.Join(tags[:], ","))
	}

	sf := reflect.StructField{
		Name: d.Name,
		Type: d.Type,
		Tag:  reflect.StructTag(sTag),
	}

	return sf
}

func (dfl DynamicFieldList) ToStructFieldList() []reflect.StructField {
	structFields := []reflect.StructField{}

	for _, v := range dfl {
		structFields = append(structFields, v.ToStructField())
	}

	return structFields
}
