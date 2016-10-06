package gomu

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// String is a nullable string.
type String struct {
	String string
	Null   bool
	Valid  bool
}

// NewString creates a new String.
func NewString(s string, n bool, valid bool) String {
	return String{
		String: s,
		Null:   n,
		Valid:  valid,
	}
}

// StringFrom creates a new String that will never be blank.
func StringFrom(s string) String {
	return NewString(s, false, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", true, true)
	}
	return NewString(*s, false, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.String = x
	case nil:
		s.Null = true
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type nzgo.String", reflect.TypeOf(v).Name())
	}
	s.Valid = err == nil
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *String) UnmarshalText(text []byte) (err error) {
	s.String = string(text)
	s.Valid = true
	return
}

// MarshalJSON implements json.Marshaler.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nil, nil
	}
	if s.Null {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// MarshalText implements encoding.TextMarshaler.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	if s.Null {
		return []byte("null"), nil
	}
	return []byte(s.String), nil
}

// SetValid changes this String's value and also sets Valid to be true.
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null or not valid.
func (s String) Ptr() *string {
	if !s.Valid || s.Null {
		return nil
	}
	return &s.String
}
