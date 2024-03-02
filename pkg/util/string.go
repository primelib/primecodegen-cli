package util

import (
	"strings"
)

// FirstNonEmptyString returns the first non-empty string from the input strings
func FirstNonEmptyString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// UpperCaseFirstLetter capitalizes the first letter of the input string
func UpperCaseFirstLetter(input string) string {
	if input == "" {
		return input
	}

	return strings.ToUpper(input[0:1]) + input[1:]
}

// LowerCaseFirstLetter lowercases the first letter of the input string
func LowerCaseFirstLetter(input string) string {
	if input == "" {
		return input
	}

	return strings.ToLower(input[0:1]) + input[1:]
}

// TrimNonASCII removes non-ASCII characters from the input string
func TrimNonASCII(input string) string {
	return strings.Map(func(r rune) rune {
		if r > 127 {
			return -1
		}
		return r
	}, input)
}

// FindCommonStrPrefix returns the common prefix of all provided strings if any
func FindCommonStrPrefix(values []string) string {
	if len(values) <= 1 {
		return ""
	}

	prefix := values[0]
	for _, str := range values[1:] {
		for !strings.HasPrefix(str, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return "" // If the prefix becomes empty, there's no common prefix
			}
		}
	}

	return prefix
}
