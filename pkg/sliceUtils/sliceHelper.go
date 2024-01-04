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

func CountOccurencesWithPrefix(data []string, prefix string) int {
	count := 0
	for _, item := range data {
		if strings.HasPrefix(item, prefix) {
			count++
		}
	}
	return count
}
