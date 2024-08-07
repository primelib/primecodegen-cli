package util

import (
	"testing"
)

func TestCapitalizeAfterChars(t *testing.T) {
	testCases := []struct {
		input           string
		chars           []int32
		capitalizeFirst bool
		expected        string
	}{
		{"hello/world", []int32{'/'}, false, "helloWorld"},
		{"_hello", []int32{'_'}, false, "Hello"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := CapitalizeAfterChars(tc.input, tc.chars, tc.capitalizeFirst)
			if result != tc.expected {
				t.Errorf("CharToCapitalize(%s, %v, %t) = %s; expected %s", tc.input, tc.chars, tc.capitalizeFirst, result, tc.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"HelloWorld", "HelloWorld"},
		{"plain text", "PlainText"},
		{"PLAIN TEXT", "PlainText"},
		{"api_endpoint", "APIEndpoint"}, // acronyms
		{"vcs-release", "VCSRelease"},   // acronyms
		{"", ""},
	}

	for _, test := range tests {
		result := ToPascalCase(test.input)
		if result != test.expected {
			t.Errorf("ToPascalCase(%s) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello_world"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("ToSnakeCase(%s) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello-world"},
		{"helloWorld", "hello-world"},
		{"hello_world", "hello-world"},
		{"hello-world", "hello-world"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToKebabCase(test.input)
		if result != test.expected {
			t.Errorf("ToKebabCase(%s) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "helloWorld"},
		{"helloWorld", "helloWorld"},
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToCamelCase(test.input)
		if result != test.expected {
			t.Errorf("ToCamelCase(%s) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"hello world", "hello-world"},
		{"hello_world", "hello-world"},
		{"hello-world", "hello-world"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToSlug(test.input)
		if result != test.expected {
			t.Errorf("ToSlug(%s) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}
