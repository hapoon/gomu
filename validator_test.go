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

type testStrcutGomuIntRequired struct {
	ID Int `valid:"required"`
}

type testStructAttachedTagsGomu struct {
	ID   Int    `valid:"required"`
	Name String `valid:"required, stringlength(1|10)"`
	URL  String `valid:"required, requri"`
	Date Time   `valid:"required"`
	Flag Bool   `valid:"required"`
}

type testGomu struct {
	ID   Int    `valid:"required"`
	Name String `valid:"stringlength(1|10)"`
	URL  String `valid:"requri"`
	Date Time   `valid:"required"`
	Flag Bool   `valid:"required"`
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

	type testStrcutGomuInt struct {
		ID Int
	}

	type testStrcutRequiredGomuInt struct {
		ID Int `valid:"required"`
	}

	// not reuqired Gomu Int
	gomuIntTests := []struct {
		param    testStrcutGomuInt
		expected bool
	}{
		{testStrcutGomuInt{Int{0, false, false}}, true},
		{testStrcutGomuInt{IntFrom(1)}, true},
		{testStrcutGomuInt{IntFromPtr(nil)}, true},
	}

	for _, test := range gomuIntTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Int
	requiredGomuIntTests := []struct {
		param    testStrcutRequiredGomuInt
		expected bool
	}{
		{testStrcutRequiredGomuInt{Int{0, false, false}}, false},
		{testStrcutRequiredGomuInt{IntFrom(1)}, true},
		{testStrcutRequiredGomuInt{IntFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuIntTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateGomuBool(t *testing.T) {
	t.Parallel()

	type testStrcutGomuBool struct {
		ID Bool
	}

	type testStrcutRequiredGomuBool struct {
		ID Bool `valid:"required"`
	}

	// not reuqired Gomu Bool
	gomuBoolTests := []struct {
		param    testStrcutGomuBool
		expected bool
	}{
		{testStrcutGomuBool{Bool{false, false, false}}, true},
		{testStrcutGomuBool{BoolFrom(false)}, true},
		{testStrcutGomuBool{BoolFromPtr(nil)}, true},
	}

	for _, test := range gomuBoolTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Bool
	requiredGomuBoolTests := []struct {
		param    testStrcutRequiredGomuBool
		expected bool
	}{
		{testStrcutRequiredGomuBool{Bool{false, false, false}}, false},
		{testStrcutRequiredGomuBool{BoolFrom(false)}, true},
		{testStrcutRequiredGomuBool{BoolFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuBoolTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}

func TestValidateGomuTime(t *testing.T) {
	t.Parallel()

	type testStrcutGomuTime struct {
		ID Time
	}

	type testStrcutRequiredGomuTime struct {
		ID Time `valid:"required"`
	}

	now := time.Now()
	// not reuqired Gomu Time
	gomuTimeTests := []struct {
		param    testStrcutGomuTime
		expected bool
	}{
		{testStrcutGomuTime{Time{time.Time{}, false, false}}, true},
		{testStrcutGomuTime{TimeFrom(now)}, true},
		{testStrcutGomuTime{TimeFromPtr(nil)}, true},
	}

	for _, test := range gomuTimeTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}

	// reuqired Gomu Time
	requiredGomuTimeTests := []struct {
		param    testStrcutRequiredGomuTime
		expected bool
	}{
		{testStrcutRequiredGomuTime{Time{time.Time{}, false, false}}, false},
		{testStrcutRequiredGomuTime{TimeFrom(now)}, true},
		{testStrcutRequiredGomuTime{TimeFromPtr(nil)}, false},
	}

	for _, test := range requiredGomuTimeTests {
		actual, err := Validate(test.param)
		ignoreError(err)
		assert.Equal(t, test.expected, actual, "Expected IsRequestURL(%+v) to be %v, got %v", test.param, test.expected, actual)
	}
}
