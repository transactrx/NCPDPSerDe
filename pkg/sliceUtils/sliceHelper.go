package sliceutils

import (
	"slices"
	"strings"
)

func IndexOfStartsWith(lookfor string, lookin []string) int {
	// Define a function that returns true when the element starts with the requested value
	startsWithFunc := func(s string) bool {
		return strings.HasPrefix(s, lookfor)
	}

	return slices.IndexFunc(lookin, startsWithFunc)
}
