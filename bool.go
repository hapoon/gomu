package gomu

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// Bool is a nullable bool.
// If key is assigned and value is not null, Bool is specified value, Null is false, Valid is true.
// If value is null, Null is true.
// If key is not assigned, Valid is false.
type Bool struct {
	Bool  bool
	Null  bool
	Valid bool
}

// NewBool creates a new Bool.
func NewBool(b bool, n bool, valid bool) Bool {
	return Bool{
		Bool:  b,
		Null:  n,
		Valid: valid,
	}
}

// BoolFrom creates a new Bool that will always be valid
func BoolFrom(b bool) Bool {
	return NewBool(b, false, true)
}

// BoolFromPtr creates a new Bool that will be null if b is nil
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, true, true)
	}
	return NewBool(*b, false, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bool) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case nil:
		b.Null = true
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type gomu.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = err == nil
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bool) UnmarshalText(text []byte) (err error) {
	if text == nil {
		return
	}
	switch string(text) {
	case "", "null":
		b.Null = true
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		err = fmt.Errorf("invalid input: %s", string(text))
	}
	b.Valid = err == nil
	return
}

// MarshalJSON implements json.Marshaler.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if b.Null {
		return []byte("null"), nil
	}
	if b.Bool {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// MarshalText implements encoding.TextMarshaler
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	if b.Null {
		return []byte("null"), nil
	}
	if b.Bool {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// SetValid changes Bool value and also sets Valid to be true
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is null or not valid.
func (b Bool) Ptr() *bool {
	if !b.Valid || b.Null {
		return nil
	}
	return &b.Bool
}

// Scan implements database/sql.Scanner.
func (b *Bool) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case bool:
		b.Bool = x
	case nil:
		b.Null = true
	default:
		err = fmt.Errorf("gomu: cannot scan type %T into gomu.Bool: %v", value, value)
	}
	b.Valid = err == nil
	return
}

// Value implements database/sql.Valuer.
func (b Bool) Value() (driver.Value, error) {
	if !b.Valid || b.Null {
		return nil, nil
	}
	return b.Bool, nil
}
