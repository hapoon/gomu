package gomu

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	timeString = "2016-10-06T10:00:00Z"
	timeJSON   = []byte(`"` + timeString + `"`)
	timeObj, _ = time.Parse(time.RFC3339, timeString)
)

type testStructTime struct {
	CreateTime Time `json:"createTime"`
}

func TestTimeFrom(t *testing.T) {
	// normal value=time.Time
	target := TimeFrom(timeObj)
	expect := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "TimeFrom(time.Time) fail")
}

func TestTimeFromPtr(t *testing.T) {
	// normal value=time.Time pointer
	target := TimeFromPtr(&timeObj)
	expect := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "TimeFromPtr(&time.Time) fail")
	// normal value=nil pointer
	target = TimeFromPtr(nil)
	expect = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	assert.Equal(t, target, expect, "TimeFromPtr(nil) fail")
}

func TestUnmarshalJSONTime(t *testing.T) {
	// normal(key and value assigned)
	j := []byte(`{"createTime":"` + timeString + `"}`)
	target := testStructTime{}
	expect := testStructTime{
		CreateTime: Time{
			Time:  timeObj,
			Null:  false,
			Valid: true,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON() fail")
	// value is null
	j = []byte(`{"createTime":null}`)
	target = testStructTime{}
	expect = testStructTime{
		CreateTime: Time{
			Time:  time.Time{},
			Null:  true,
			Valid: true,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON() fail")
	// key is not assigned
	j = []byte(`{}`)
	target = testStructTime{}
	expect = testStructTime{
		CreateTime: Time{
			Time:  time.Time{},
			Null:  false,
			Valid: false,
		},
	}
	json.Unmarshal(j, &target)
	assert.Equal(t, target, expect, "UnmarshalJSON() fail")
}

func TestUnmarshalTextTime(t *testing.T) {
	// normal
	var target Time
	expect := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	err := target.UnmarshalText([]byte(timeString))
	checkError(err)
	assert.Equal(t, target, expect, `UnmarshalText([]byte(`+timeString+`)) fail`)
	// normal "null"
	target = Time{}
	expect = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	err = target.UnmarshalText([]byte("null"))
	checkError(err)
	assert.Equal(t, target, expect, `UnmarshalText([]byte("null")) fail`)
	// normal ""
	target = Time{}
	expect = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	err = target.UnmarshalText([]byte(""))
	checkError(err)
	assert.Equal(t, target, expect, `UnmarshalText([]byte("")) fail`)
	// normal null
	target = Time{}
	expect = Time{
		Time:  time.Time{},
		Null:  false,
		Valid: false,
	}
	err = target.UnmarshalText(nil)
	checkError(err)
	assert.Equal(t, target, expect, "UnmarshalText(nil) fail")
}

func TestMarshalJSONTime(t *testing.T) {
	// normal key=createTime value=timeString
	ts := testStructTime{
		CreateTime: Time{
			Time:  timeObj,
			Null:  false,
			Valid: true,
		},
	}
	expect := []byte(`{"createTime":"` + timeString + `"}`)
	target, err := json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON(`+string(expect)+`) fail`)
	// normal key=createTime value=null
	ts = testStructTime{
		CreateTime: Time{
			Time:  time.Time{},
			Null:  true,
			Valid: true,
		},
	}
	expect = []byte(`{"createTime":null}`)
	target, err = json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON(`+string(expect)+`) fail`)
	// normal key not assigned
	ts = testStructTime{
		CreateTime: Time{
			Time:  time.Time{},
			Null:  false,
			Valid: false,
		},
	}
	expect = []byte(`{"createTime":null}`)
	target, err = json.Marshal(ts)
	checkError(err)
	assert.Equal(t, target, expect, `MarshalJSON(`+string(expect)+`) fail`)
}

func TestMarshalTextTime(t *testing.T) {
	// normal timeString
	ts := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	expect := []byte(timeString)
	target, err := ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, "MarshalText("+timeString+") fail")
	// normal null
	ts = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	expect = []byte("null")
	target, err = ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, "MarshalText(null) fail")
	// normal key is not assigned
	ts = Time{
		Time:  time.Time{},
		Null:  false,
		Valid: false,
	}
	expect = []byte(nil)
	target, err = ts.MarshalText()
	checkError(err)
	assert.Equal(t, target, expect, "MarshalText() fail")
}

func TestSetValidTime(t *testing.T) {
	// normal
	target := Time{}
	target.SetValid(timeObj)
	expect := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	assert.Equal(t, target, expect, "SetValid("+timeString+")")
}

func TestPtrTime(t *testing.T) {
	// normal timeString
	ts := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	target := ts.Ptr()
	assert.Equal(t, *target, timeObj, "Ptr() fail")
	// normal null
	ts = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() fail")
	// normal not valid
	ts = Time{
		Time:  time.Time{},
		Null:  false,
		Valid: false,
	}
	target = ts.Ptr()
	assert.Nil(t, target, "Ptr() fail")
}

func TestScanTime(t *testing.T) {
	// normal timeString
	target := Time{}
	expect := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	err := target.Scan(timeObj)
	checkError(err)
	assert.Equal(t, target, expect, "Scan("+timeString+") fail")
	// normal nil
	target = Time{}
	expect = Time{
		Time:  time.Time{},
		Null:  true,
		Valid: true,
	}
	err = target.Scan(nil)
	checkError(err)
	assert.Equal(t, target, expect, "Scan(nil) fail")
}

func TestValueTime(t *testing.T) {
	// normal timeString
	ts := Time{
		Time:  timeObj,
		Null:  false,
		Valid: true,
	}
	expect := timeObj
	target, err := ts.Value()
	checkError(err)
	assert.Equal(t, target, expect, "Value "+timeString+"fail")
}
