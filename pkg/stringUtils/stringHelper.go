package stringutils

import (
	"fmt"
	"slices"
	"strings"
)

// Split string at any of the specified search characters.
func SplitFirst(str string, fromIndex int, searchFor []byte) string {
	index := IndexOfAny(str, fromIndex+1, searchFor)

	if index > 0 {
		return Substring(str, fromIndex, index-fromIndex)
	}

	return Substring(str, fromIndex, -1)
}

// Find first index of any of the specified characters.
// Returns -1 when not found.
func IndexOfAny(str string, fromIndex int, searchFor []byte) int {
	if str == "" || len(searchFor) == 0 {
		return -1
	}

	for i := fromIndex; i < len(str); i++ {
		ch := str[i]

		if slices.Contains(searchFor, ch) {
			return i
		}
	}

	return -1
}

// Get a section of a string.
// Length of <=0 will return the remaining string starting at the specified index.
// Length >0 will return the length specified when in range.
func Substring(str string, fromIndex, length int) string {
	if str == "" || len(str) < fromIndex {
		return ""
	}

	if length <= 0 {
		return str[fromIndex:]
	}

	endIndex := fromIndex + length

	if len(str) < endIndex {
		return str[fromIndex:]
	}

	return str[fromIndex:endIndex]
}

// Right pad to specified length, or truncate if length is too long.
func RightPadExact(str string, padChar byte, length int) string {
	if len(str) >= length {
		return str[:length]
	}

	padder := string(padChar)
	numNeeded := length - len(str)

	return fmt.Sprintf("%v%v", str, strings.Repeat(padder, numNeeded))
}

// Trim all specified characters from the front and end of the string.
func TrimAll(str string, trimChars []byte) string {
	trimFunc := func(r rune) bool {
		b := byte(r)

		return slices.Contains(trimChars, b)
	}

	str = strings.TrimLeftFunc(str, trimFunc)
	return strings.TrimRightFunc(str, trimFunc)
}

func SplitBySeparator(data string, sep byte) []string {
	sepIndex := IndexOfAny(data, 0, []byte{sep})

	// None present
	if sepIndex < 0 {
		return []string{}
	}

	return strings.Split(data[sepIndex+1:], string(sep))
}
