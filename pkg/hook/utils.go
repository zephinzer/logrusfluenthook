package hook

import (
	"strings"
)

// errorContainsString returns true if the error message of :haystack contains
// the string :needle
func errorContainsString(haystack error, needle string) bool {
	return strings.Contains(haystack.Error(), needle)
}

// sliceContainsString returns true if the provided
func sliceContainsString(haystack []string, needle string) bool {
	for _, hay := range haystack { // no need for order
		if hay == needle {
			return true
		}
	}
	return false
}
