package gomu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStructStringLength struct {
	Name String `valid:"stringlength(1|10)"`
}

type testStructLength struct {
	Name String `valid:"length(1|10)"`
}

type testStructReqURL struct {
	URL String `valid:"requrl"`
}

type testStructReqURI struct {
	URI String `valid:"requri"`
}

func TestValidateStringLength(t *testing.T) {
	// stringlength true
	test := testStructStringLength{
		Name: StringFrom("1234567890"),
	}
	result, err := Validate(test)
	ignoreError(err)
	assert.True(t, result, "Validate(stringlength) fail")
	// stringlength(multibyte) true
	test = testStructStringLength{
		Name: StringFrom("あいうえおかきくけこ"),
	}
	result, err = Validate(test)
	ignoreError(err)
	assert.True(t, result, "Validate(stringlength) fail")
	// stringlength false
	test = testStructStringLength{
		Name: StringFrom("12345678901"),
	}
	result, err = Validate(test)
	ignoreError(err)
	assert.False(t, result, "Validate(stringlength) fail")
}

func TestValidateLength(t *testing.T) {
	// length true
	test := testStructLength{
		Name: StringFrom("1234567890"),
	}
	result, err := Validate(test)
	ignoreError(err)
	assert.True(t, result, "Validate(length) fail")
	// length false
	test = testStructLength{
		Name: StringFrom("12345678901"),
	}
	result, err = Validate(test)
	ignoreError(err)
	assert.False(t, result, "Validate(length) fail")
	// length(multibyte) false
	test = testStructLength{
		Name: StringFrom("あいうえ"),
	}
	result, err = Validate(test)
	ignoreError(err)
	assert.False(t, result, "Validate(length) fail")
}

func TestValidateIsRequestURL(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://sample.com", true},
		{"https://sam.ple.com/", true},
		{"ftp://sample.abc/", true},
		{"/abc/def/ghi", false},
		{"http://sample.com:1234", true},
		{"http://sample.com:123456", true},
		{"abc", false},
	}
	for _, test := range tests {
		actual := IsRequestURL(test.param)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%q) to be %v, got %v", test.param, test.expected, actual)
	}

	var tests2 = []struct {
		param    testStructReqURL
		expected bool
	}{
		{testStructReqURL{StringFrom("")}, false},
		{testStructReqURL{StringFrom("http://sample.com")}, true},
	}
	for _, test := range tests2 {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateIsRequestURI(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://sample.com", true},
		{"https://sam.ple.com/", true},
		{"ftp://sample.abc/", true},
		{"/abc/def/ghi", true},
		{"http://sample.com:1234", true},
		{"http://sample.com:123456", true},
		{"abc", false},
	}
	for _, test := range tests {
		actual := IsRequestURI(test.param)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURI(%q) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateIsURL(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://sample.com", true},
		{"https://sam.ple.com/", true},
		{"ftp://sample.abc/", true},
		{"/abc/def/ghi", false},
		{"http://sample.com:1234", true},
		{"http://sample.com:123456", false},
		{"abc", false},
	}

	for _, test := range tests {
		actual := IsURL(test.param)
		assert.Equal(t, test.expected, actual, "Expected IsURL(%q) to be %v, got %v", test.param, test.expected, actual)
	}
}
