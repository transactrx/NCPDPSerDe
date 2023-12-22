package ncpdp

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

type FieldSettings struct {
	DecimalPlaces int
	Format        string
	Overpunch     bool
	Repeating     bool
}

const (
	negativeSign = "-"
)

var signedPositive = map[string]string{
	"{": "0",
	"A": "1",
	"B": "2",
	"C": "3",
	"D": "4",
	"E": "5",
	"F": "6",
	"G": "7",
	"H": "8",
	"I": "9",
}

var signedNegative = map[string]string{
	"}": "0",
	"J": "1",
	"K": "2",
	"L": "3",
	"M": "4",
	"N": "5",
	"O": "6",
	"P": "7",
	"Q": "8",
	"R": "9",
}

func (fs FieldSettings) toImpliedDecimalString(value interface{}) string {
	pattern := "%." + fmt.Sprint(fs.DecimalPlaces) + "f"
	return strings.Replace(fmt.Sprintf(pattern, value), ".", Empty, 1)
}

// Parse field to float data type with implied decimals.
// Decimal places may be zero to indicate no decimals.
func (fs FieldSettings) parseFloatWithImpliedDecimals(value string) (*float64, error) {
	if value == Empty {
		return nil, nil
	}

	divisor := math.Pow(10, float64(fs.DecimalPlaces))

	i, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse value to float.  Value: %v  Error: %q", value, err)
	}

	result := float64(i) / divisor

	return &result, nil
}

// Unsign potentially overpunched value as float with implied decimals
func (fs FieldSettings) Unsign(value string) (*float64, error) {
	rawValue := removeOverpunch(value)
	fltValue, err := fs.parseFloatWithImpliedDecimals(rawValue)
	if err != nil {
		return nil, err
	}

	return fltValue, nil
}

// Apply overpunch to value with implied decimal places.
func (fs FieldSettings) sign(value interface{}) string {
	return applyOverpunch(fs.toImpliedDecimalString(value))
}

// Apply overpunch to value.
func applyOverpunch(value string) string {
	if value == Empty {
		return Empty
	}

	lastChar := stringutils.Substring(value, len(value)-1, 1)
	firstChar := stringutils.Substring(value, 0, 1)

	signedValues := signedPositive
	if lastChar == negativeSign || firstChar == negativeSign {
		signedValues = signedNegative
		value = strings.Replace(value, negativeSign, Empty, -1)
	}

	lastChar = stringutils.Substring(value, len(value)-1, 1)

	for key, val := range signedValues {
		if val == lastChar {
			return fmt.Sprintf("%v%v", stringutils.Substring(value, 0, len(value)-1), key)
		}
	}

	return value
}

// Remove overpunch from field.
func removeOverpunch(value string) string {
	if value == Empty {
		return Empty
	}

	lastChar := stringutils.Substring(value, len(value)-1, 1)

	// Look for signed negative value.
	// Replace char found with the index which signifies the actual number.
	negValue := signedNegative[lastChar]
	if negValue != Empty {
		return fmt.Sprintf("-%v%v", stringutils.Substring(value, 0, len(value)-1), negValue)
	}

	// Look for signed positive value.
	// Replace char found with the index which signifies the actual number.
	posValue := signedPositive[lastChar]
	if posValue != Empty {
		return fmt.Sprintf("%v%v", stringutils.Substring(value, 0, len(value)-1), posValue)
	}

	// If none found, field was not overpunched.
	// Return raw value.
	return value
}

// Convert field value to raw string representation
func (fs FieldSettings) convertFieldValueToString(fieldValue interface{}) string {
	switch t := fieldValue.(type) {
	case float64:
		if fs.Overpunch {
			return fs.sign(t)
		}

		return fs.toImpliedDecimalString(t)
	}

	return fmt.Sprint(fieldValue)
}
