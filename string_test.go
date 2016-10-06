package gomu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStructString struct {
	Name String `json:"name"`
}

func TestStringFrom(t *testing.T) {
	// "test"
	target := StringFrom("test")
	expect := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	assert.Equal(t, target, expect, `StringFrom("test") fail`)
}

func TestStringFromPtr(t *testing.T) {
	// a string pointer
	s := "test"
	target := StringFromPtr(&s)
	expect := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	assert.Equal(t, target, expect, "StringFromPtr() fail")
	// a nil pointer
	target = StringFromPtr(nil)
	expect = String{
		String: "",
		Null:   true,
		Valid:  true,
	}
	assert.Equal(t, target, expect, "StringFromPtr(nil) fail")
}

func TestUnmarshalJSONString(t *testing.T) {
	// key and value assigned
	j := []byte(`{"name":"test"}`)
	target := testStructString{}
	expect := testStructString{
		Name: String{
			String: "test",
			Null:   false,
			Valid:  true,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON fail")
	// value is null
	j = []byte(`{"name":null}`)
	target = testStructString{}
	expect = testStructString{
		Name: String{
			String: "",
			Null:   true,
			Valid:  true,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON(null) fail")
	// key is not assigned
	j = []byte(`{}`)
	target = testStructString{}
	expect = testStructString{
		Name: String{
			String: "",
			Null:   false,
			Valid:  false,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON(key is not assigned) fail")
}

func TestUnmarshalTextString(t *testing.T) {
	// key and value assigned
	target := String{}
	expect := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	err := target.UnmarshalText([]byte("test"))
	checkError(err)
	assert.Equal(t, target, expect, "UnmarshalText() fail")
	// value is ""
	target = String{}
	expect = String{
		String: "null",
		Null:   true,
		Valid:  true,
	}
	err = target.UnmarshalText([]byte("null"))
	checkError(err)
	assert.Equal(t, target, expect, `UnmarshalText("null") fail`)
	// value is null
	target = String{}
	expect = String{
		String: "",
		Null:   false,
		Valid:  false,
	}
	err = target.UnmarshalText(nil)
	checkError(err)
	assert.Equal(t, target, expect, "UnmarshalText(null)")
}

func TestMarshalJSONString(t *testing.T) {
	// key=Name value=test
	ts := testStructString{
		Name: StringFrom("test"),
	}
	expect := []byte(`{"name":"test"}`)
	target, err := json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON("test") fail`)
	// key=Name value=null
	ts = testStructString{
		Name: String{
			String: "",
			Null:   true,
			Valid:  true,
		},
	}
	expect = []byte(`{"name":null}`)
	target, err = json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON(null) fail`)
	// key is not assigned
	ts = testStructString{
		Name: String{
			String: "",
			Null:   false,
			Valid:  false,
		},
	}
	expect = []byte(`{"name":null}`)
	target, err = json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON(key is not assigned) fail`)
}

func TestMarshalTextString(t *testing.T) {
	// "test"
	ts := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	expect := []byte("test")
	target, err := ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, `MarshalText("test") fail`)
	// null
	ts = String{
		String: "",
		Null:   true,
		Valid:  true,
	}
	expect = []byte("null")
	target, err = ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, `MarshalText(null) fail`)
	// key is not assigned
	ts = String{
		String: "",
		Null:   false,
		Valid:  false,
	}
	expect = []byte(nil)
	target, err = ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, `MarshalText() fail`)
}

func TestSetValidString(t *testing.T) {
	// "test"
	target := String{}
	target.SetValid("test")
	expect := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	assert.Equal(t, target, expect, `SetValid("test") fail`)
}

func TestPtrString(t *testing.T) {
	// "test"
	ts := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	target := ts.Ptr()
	expect := "test"
	assert.Equal(t, *target, expect, "Ptr() fail")
	// null
	ts = String{
		String: "",
		Null:   true,
		Valid:  true,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() null fail")
	// not Valid
	ts = String{
		String: "",
		Null:   false,
		Valid:  false,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() not valid fail")
}

func TestScanString(t *testing.T) {
	// "test"
	target := String{}
	expect := String{
		String: "test",
		Null:   false,
		Valid:  true,
	}
	err := target.Scan("test")
	checkError(err)
	assert.Equal(t, target, expect, `Scan("test") fail`)
}
