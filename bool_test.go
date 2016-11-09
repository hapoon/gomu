package gomu

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStructBool struct {
	IsClear Bool `json:"isClear"`
}

func TestBoolFrom(t *testing.T) {
	// true
	target := BoolFrom(true)
	expect := NewBool(true, false, true)
	assert.Equal(t, target, expect, "not equal")
	// false
	target = BoolFrom(false)
	expect = NewBool(false, false, true)
	assert.Equal(t, target, expect, "not equal")
}

func TestBoolFromPtr(t *testing.T) {
	// a bool pointer
	b := true
	target := BoolFromPtr(&b)
	expect := NewBool(true, false, true)
	assert.Equal(t, target, expect, "not equal")
	// a nil pointer
	target = BoolFromPtr(nil)
	expect = NewBool(false, true, true)
	assert.Equal(t, target, expect, "not equal")
}

func TestUnmarshalJSONBool(t *testing.T) {
	// normal pattern(key and value assigned))
	j := []byte(`{"isClear":true}`)
	ts := testStructBool{}
	expect := testStructBool{
		IsClear: NewBool(true, false, true),
	}
	json.Unmarshal(j, &ts)
	assert.Equal(t, ts, expect, "not equal")
	// value is Null
	j = []byte(`{"isClear":null}`)
	ts = testStructBool{}
	expect = testStructBool{
		IsClear: NewBool(false, true, true),
	}
	json.Unmarshal(j, &ts)
	assert.Equal(t, ts, expect, "not equal")
	// key is not assigned
	j = []byte(`{}`)
	ts = testStructBool{}
	expect = testStructBool{
		IsClear: NewBool(false, false, false),
	}
	json.Unmarshal(j, &ts)
	assert.Equal(t, ts, expect, "not equal")
}

func TestUnmarshalTextBool(t *testing.T) {
	// normal "true"
	var b Bool
	expect := NewBool(true, false, true)
	err := b.UnmarshalText([]byte("true"))
	checkError(err)
	assert.Equal(t, b, expect, `UnmarshalText([]byte("true"))`)
	// normal "false"
	b = Bool{}
	expect = NewBool(false, false, true)
	err = b.UnmarshalText([]byte("false"))
	checkError(err)
	assert.Equal(t, b, expect, `UnmarshalText([]byte("false"))`)
	// normal "null"
	b = Bool{}
	expect = NewBool(false, true, true)
	err = b.UnmarshalText([]byte("null"))
	checkError(err)
	assert.Equal(t, b, expect, `UnmarshalText([]byte("null"))`)
	// normal ""
	b = Bool{}
	expect = NewBool(false, true, true)
	err = b.UnmarshalText([]byte(""))
	checkError(err)
	assert.Equal(t, b, expect, `UnmarshalText([]byte(""))`)
	// normal null
	b = Bool{}
	expect = NewBool(false, false, false)
	err = b.UnmarshalText(nil)
	checkError(err)
	assert.Equal(t, b, expect, `UnmarshalText(nil)`)
}

func TestMarshalJSONBool(t *testing.T) {
	// normal key=isClear value=true
	j := testStructBool{
		IsClear: BoolFrom(true),
	}
	expect := []byte(`{"isClear":true}`)
	data, err := json.Marshal(j)
	checkError(err)
	assert.Equal(t, data, expect, `MarshalJSON(true) fail`)
	// normal key=isClear value=null
	j = testStructBool{
		IsClear: NewBool(false, true, true),
	}
	expect = []byte(`{"isClear":null}`)
	data, err = json.Marshal(j)
	checkError(err)
	assert.Equal(t, data, expect, `MarshalJSON(null) fail`)
	// normal key not assigned
	j = testStructBool{
		IsClear: NewBool(false, false, false),
	}
	expect = []byte(`{"isClear":null}`)
	data, err = json.Marshal(j)
	checkError(err)
	assert.Equal(t, data, expect, `MarshalJSON() fail`)
}

func TestMarshalTextBool(t *testing.T) {
	// normal true
	b := NewBool(true, false, true)
	expect := []byte("true")
	data, err := b.MarshalText()
	checkError(err)
	assert.Equal(t, data, expect, "MarshalText(true) fail")
	// normal null
	b = NewBool(false, true, true)
	expect = []byte("null")
	data, err = b.MarshalText()
	checkError(err)
	assert.Equal(t, data, expect, "MarshalText(null) fail")
	// normal key is not assigned
	b = NewBool(false, false, false)
	expect = []byte(nil)
	data, err = b.MarshalText()
	checkError(err)
	assert.Equal(t, data, expect, "MarshalText() fail")
}

func TestSetValidBool(t *testing.T) {
	// normal true
	b := Bool{}
	b.SetValid(true)
	expect := NewBool(true, false, true)
	assert.Equal(t, b, expect, "SetValid(true) fail")
	// normal false
	b = Bool{}
	b.SetValid(false)
	expect = NewBool(false, false, true)
	assert.Equal(t, b, expect, "SetValid(false) fail")
}

func TestPtrBool(t *testing.T) {
	// normal true
	b := NewBool(true, false, true)
	ptr := b.Ptr()
	expect := true
	assert.Equal(t, *ptr, expect, "Bool{Bool:true}.Ptr() fail")
	// normal null
	b = NewBool(false, true, true)
	ptr = b.Ptr()
	assert.Nil(t, ptr, "Bool{Null:true}.Ptr() fail")
	// normal not Valid
	b = NewBool(false, false, false)
	ptr = b.Ptr()
	assert.Nil(t, ptr, "Bool{Valid:false}.Ptr() fail")
}

func TestScanBool(t *testing.T) {
	// normal Scan(true)
	var b Bool
	expect := NewBool(true, false, true)
	err := b.Scan(true)
	checkError(err)
	assert.Equal(t, b, expect, "Scan(true) fail")
	// normal Scan(nil)
	b = Bool{}
	expect = NewBool(false, true, true)
	err = b.Scan(nil)
	checkError(err)
	assert.Equal(t, b, expect, "Scan(nil) fail")
}

func TestValueBool(t *testing.T) {
	// normal true
	b := NewBool(true, false, true)
	expect := true
	data, err := b.Value()
	checkError(err)
	assert.Equal(t, data, driver.Value(expect), "Value true fail")
	// normal null
	b = NewBool(false, true, true)
	data, err = b.Value()
	checkError(err)
	assert.Nil(t, data, "Value null fail")
	// normal not assigned
	b = NewBool(false, false, false)
	data, err = b.Value()
	checkError(err)
	assert.Nil(t, data, "Value not assigned fail")
}
