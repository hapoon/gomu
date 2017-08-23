package gomu

import (
	"testing"
	"time"

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
		{testStructReqURL{StringFromPtr(nil)}, true},
		{testStructReqURL{NewString("", false, false)}, true},
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

	var tests2 = []struct {
		param    testStructReqURI
		expected bool
	}{
		{testStructReqURI{StringFrom("")}, false},
		{testStructReqURI{StringFrom("abc")}, false},
		{testStructReqURI{StringFrom("http://sample.com")}, true},
		{testStructReqURI{StringFromPtr(nil)}, true},
		{testStructReqURI{NewString("", false, false)}, true},
	}
	for _, test := range tests2 {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURI(%+v) to be %v, got %v", test.param, test.expected, actual)
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

func TestValidateGomuInt(t *testing.T) {
	t.Parallel()

	type testStructGomuInt struct {
		ID Int
	}

	type testStructRequiredGomuInt struct {
		ID Int `valid:"required"`
	}

	// not reuqired Gomu Int
	gomuIntTests := []struct {
		param    testStructGomuInt
		expected bool
	}{
		{testStructGomuInt{Int{0, false, false}}, true},
		{testStructGomuInt{IntFrom(1)}, true},
		{testStructGomuInt{IntFromPtr(nil)}, true},
	}

	for _, test := range gomuIntTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Int
	requiredGomuIntTests := []struct {
		param    testStructRequiredGomuInt
		expected bool
	}{
		{testStructRequiredGomuInt{Int{0, false, false}}, false},
		{testStructRequiredGomuInt{IntFrom(1)}, true},
		{testStructRequiredGomuInt{IntFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuIntTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateGomuBool(t *testing.T) {
	t.Parallel()

	type testStructGomuBool struct {
		ID Bool
	}

	type testStructRequiredGomuBool struct {
		ID Bool `valid:"required"`
	}

	// not reuqired Gomu Bool
	gomuBoolTests := []struct {
		param    testStructGomuBool
		expected bool
	}{
		{testStructGomuBool{Bool{false, false, false}}, true},
		{testStructGomuBool{BoolFrom(false)}, true},
		{testStructGomuBool{BoolFromPtr(nil)}, true},
	}

	for _, test := range gomuBoolTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected Validate(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Bool
	requiredGomuBoolTests := []struct {
		param    testStructRequiredGomuBool
		expected bool
	}{
		{testStructRequiredGomuBool{Bool{false, false, false}}, false},
		{testStructRequiredGomuBool{BoolFrom(false)}, true},
		{testStructRequiredGomuBool{BoolFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuBoolTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected Validate(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateGomuTime(t *testing.T) {
	t.Parallel()

	type testStructGomuTime struct {
		ID Time
	}

	type testStructRequiredGomuTime struct {
		ID Time `valid:"required"`
	}

	now := time.Now()
	// not reuqired Gomu Time
	gomuTimeTests := []struct {
		param    testStructGomuTime
		expected bool
	}{
		{testStructGomuTime{Time{time.Time{}, false, false}}, true},
		{testStructGomuTime{TimeFrom(now)}, true},
		{testStructGomuTime{TimeFromPtr(nil)}, true},
	}

	for _, test := range gomuTimeTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected Validate(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Time
	requiredGomuTimeTests := []struct {
		param    testStructRequiredGomuTime
		expected bool
	}{
		{testStructRequiredGomuTime{Time{time.Time{}, false, false}}, false},
		{testStructRequiredGomuTime{TimeFrom(now)}, true},
		{testStructRequiredGomuTime{TimeFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuTimeTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected Validate(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}
