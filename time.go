package gomu

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Time is a nullable Time.
type Time struct {
	Time  time.Time
	Null  bool
	Valid bool
}

// NewTime creates a new Time.
func NewTime(t time.Time, null bool, valid bool) Time {
	return Time{
		Time:  t,
		Null:  null,
		Valid: valid,
	}
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, false, true)
}

// TimeFromPtr creates a new Time that will be null if t is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, true, true)
	}
	return NewTime(*t, false, true)
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
// It supports string, object, and null input.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case string:
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		nu, nuOK := x["Null"].(bool)
		va, vaOK := x["Valid"].(bool)
		if !tiOK || !nuOK || !vaOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type gomu.Time requires key "Time" to be of type string and key "Null" to be of type bool and key "Valid" to be of type bool; found %T and %T and %T, respectively`, x["Time"], x["Null"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Null = nu
		t.Valid = va
		return
	case nil:
		t.Null = true
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type gomu.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Time) UnmarshalText(text []byte) (err error) {
	if text == nil {
		return
	}
	str := string(text)
	if str == "" || str == "null" {
		t.Null = true
		t.Valid = true
		return
	}
	if err = t.Time.UnmarshalText(text); err != nil {
		return
	}
	t.Valid = true
	return
}

// MarshalJSON implements encoding/json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	if t.Null {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

// MarshalText implements encoding.TextMarshaler.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return nil, nil
	}
	if t.Null {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

// SetValid changes this Time value and sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Null = false
	t.Valid = true
}

// Ptr returns a pointer to this Time value, or a nil pointer if this Time is null or not valid.
func (t Time) Ptr() *time.Time {
	if !t.Valid || t.Null {
		return nil
	}
	return &t.Time
}

// Scan implements database/sql.Scanner.
func (t *Time) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Null = true
	default:
		err = fmt.Errorf("gomu: cannot scan type %T into gomu.Time: %v", value, value)
	}
	t.Valid = err == nil
	return
}

// Value implements database/sql.Valuer.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid || t.Null {
		return nil, nil
	}
	return driver.Value(t.Time), nil
}
