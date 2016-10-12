package gomu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Parallel()

	customErr := &Error{
		Name: "Custom error name",
		Err:  fmt.Errorf("Custom error"),
	}
	customErrWithCustomErrorMessage := &Error{
		Name: "Custom error with custom error message name",
		Err:  fmt.Errorf("Custom error2"),
		CustomErrorMessageExists: true,
	}

	var tests = []struct {
		param    Errors
		expected string
	}{
		{Errors{}, ""},
		{Errors{fmt.Errorf("Error 1")}, "Error 1;"},
		{Errors{fmt.Errorf("Error 1"), fmt.Errorf("Error 2")}, "Error 1;Error 2;"},
		{Errors{customErr, fmt.Errorf("Error 2")}, "Custom error name: Custom error;Error 2;"},
		{Errors{fmt.Errorf("Error 123"), customErrWithCustomErrorMessage}, "Error 123;Custom error2;"},
	}
	for _, test := range tests {
		actual := test.param.Error()
		assert.Equal(t, test.expected, actual, "Expected Error() to return '%v', got '%v'", test.expected, actual)
	}
}
