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
