package ncpdp

import (
	"fmt"
	"strconv"
)

type NcpdpField struct {
	Id         string
	Value      string
	RawValue   string
	StartIndex int
}

// Create new field
func NewField(id, value string) *NcpdpField {
	field := NcpdpField{Id: id, Value: value}
	field.RawValue = fmt.Sprintf("%v%v", id, value)

	return &field
}

// Get value as string
func (field *NcpdpField) GetString() string {
	if field == nil || field.Value == Empty {
		return Empty
	}

	return field.Value
}

// Get value as integer
func (field *NcpdpField) GetInt() *int {
	if field == nil || field.Value == Empty {
		return nil
	}

	i, err := strconv.Atoi(field.Value)
	if err != nil {
		return nil
	}

	return &i
}

// Get value as float
func (field *NcpdpField) GetFloat(fs *FieldSettings) (*float64, error) {
	if field == nil {
		return nil, nil
	}

	if fs == nil {
		fs = &FieldSettings{}
	}

	return fs.unsign(field.Value)
}
