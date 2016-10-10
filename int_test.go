package gomu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStructInt struct {
	ID Int `json:"id"`
}

func TestIntFrom(t *testing.T) {
	// value > 0
	target := IntFrom(12345)
	expect := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "IntFrom(12345) fail")
	// value = 0
	target = IntFrom(0)
	expect = Int{
		Int64: 0,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "IntFrom(0) fail")
}

func TestIntFromPtr(t *testing.T) {
	// int64 pointer
	n := int64(12345)
	target := IntFromPtr(&n)
	expect := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "IntFromPtr() fail")
	// nil
	target = IntFromPtr(nil)
	expect = Int{
		Int64: 0,
		Null:  true,
		Valid: true,
	}
	assert.Equal(t, target, expect, "IntFromPtr(nil) fail")
}

func TestUnmarshalJSONInt(t *testing.T) {
	// key and value assigned
	j := []byte(`{"id":12345}`)
	target := testStructInt{}
	json.Unmarshal(j, &target)
	expect := testStructInt{
		ID: Int{
			Int64: 12345,
			Null:  false,
			Valid: true,
		},
	}
	assert.Equal(t, target, expect, "UnmarshalJSON(int) fail")
	// value is null
	j = []byte(`{"id":null}`)
	target = testStructInt{}
	json.Unmarshal(j, &target)
	expect = testStructInt{
		ID: Int{
			Int64: 0,
			Null:  true,
			Valid: true,
		},
	}
	assert.Equal(t, target, expect, "UnmarshalJSON(null) fail")
	// key is not assigned
	j = []byte("{}")
	target = testStructInt{}
	json.Unmarshal(j, &target)
	expect = testStructInt{
		ID: Int{
			Int64: 0,
			Null:  false,
			Valid: false,
		},
	}
	assert.Equal(t, target, expect, "UnmarshalJSON(key is not assigned) fail")
	// key and value assigned(value is float)
	j = []byte(`{"id":1.23}`)
	target = testStructInt{}
	json.Unmarshal(j, &target)
	expect = testStructInt{
		ID: Int{
			Int64: 0,
			Null:  false,
			Valid: false,
		},
	}
	assert.Equal(t, target, expect, "UnmarshalJSON(float) fail")
}

func TestUnmarshalTextInt(t *testing.T) {
	// key and value assigned
	target := Int{}
	err := target.UnmarshalText([]byte("12345"))
	checkError(err)
	expect := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "UnmarshalText() fail")
	// value is ""
	target = Int{}
	err = target.UnmarshalText([]byte(""))
	checkError(err)
	expect = Int{
		Int64: 0,
		Null:  true,
		Valid: true,
	}
	assert.Equal(t, target, expect, "UnmarshalText() fail")
	// value is "null"
	target = Int{}
	err = target.UnmarshalText([]byte("null"))
	checkError(err)
	expect = Int{
		Int64: 0,
		Null:  true,
		Valid: true,
	}
	assert.Equal(t, target, expect, `UnmarshalText("null") fail`)
	// value is null
	target = Int{}
	err = target.UnmarshalText(nil)
	checkError(err)
	expect = Int{
		Int64: 0,
		Null:  false,
		Valid: false,
	}
	assert.Equal(t, target, expect, "UnmarshalText(null) fail")
}

func TestMarshalJSONInt(t *testing.T) {
	// key=ID value=12345
	ts := testStructInt{
		ID: IntFrom(12345),
	}
	target, err := json.Marshal(ts)
	checkError(err)
	expect := []byte(`{"id":12345}`)
	assert.Equal(t, target, expect, "MarshalJSON(12345) fail")
	// key=ID value=null
	ts = testStructInt{
		ID: Int{
			Int64: 0,
			Null:  true,
			Valid: true,
		},
	}
	target, err = json.Marshal(ts)
	checkError(err)
	expect = []byte(`{"id":null}`)
	assert.Equal(t, target, expect, "MarshalJSON(null) fail")
	// key is not assigned
	ts = testStructInt{
		ID: Int{
			Int64: 0,
			Null:  false,
			Valid: false,
		},
	}
	target, err = json.Marshal(ts)
	checkError(err)
	expect = []byte(`{"id":null}`)
	assert.Equal(t, target, expect, "MarshalJSON(key is not assigned) fail")
}

func TestMarshalTextInt(t *testing.T) {
	// 12345
	ts := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	target, err := ts.MarshalText()
	checkError(err)
	expect := []byte("12345")
	assert.Equal(t, target, expect, "MarshalText(12345) fail")
	// null
	ts = Int{
		Int64: 0,
		Null:  true,
		Valid: true,
	}
	target, err = ts.MarshalText()
	checkError(err)
	expect = []byte("null")
	assert.Equal(t, target, expect, "MarshalText(null) fail")
	// key is not assigned
	ts = Int{
		Int64: 0,
		Null:  false,
		Valid: false,
	}
	target, err = ts.MarshalText()
	checkError(err)
	expect = []byte(nil)
	assert.Equal(t, target, expect, "MarshalText() fail")
}

func TestSetValidInt(t *testing.T) {
	// 12345
	target := Int{}
	target.SetValid(12345)
	expect := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "SetValid(12345) fail")
}

func TestPtrInt(t *testing.T) {
	// 12345
	ts := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	target := ts.Ptr()
	expect := int64(12345)
	assert.Equal(t, *target, expect, "Ptr() fail")
	// null
	ts = Int{
		Int64: 0,
		Null:  true,
		Valid: true,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() fail")
	// not Valid
	ts = Int{
		Int64: 0,
		Null:  false,
		Valid: false,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() fail")
}

func TestScanInt(t *testing.T) {
	// 12345
	target := Int{}
	expect := Int{
		Int64: 12345,
		Null:  false,
		Valid: true,
	}
	err := target.Scan(12345)
	checkError(err)
	assert.Equal(t, target, expect, `Scan(12345) fail`)
}
