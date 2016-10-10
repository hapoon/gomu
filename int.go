package gomu

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int is a nullable int.
type Int struct {
	Int64 int64
	Null  bool
	Valid bool
}

// NewInt creates a new Int.
func NewInt(i int64, n bool, valid bool) Int {
	return Int{
		Int64: i,
		Null:  n,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will never be blank.
func IntFrom(i int64) Int {
	return NewInt(i, false, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, true, true)
	}
	return NewInt(*i, false, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Int64)
	case nil:
		i.Null = true
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type gomu.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int) UnmarshalText(text []byte) (err error) {
	if text == nil {
		return
	}
	str := string(text)
	if str == "" || str == "null" {
		i.Null = true
		i.Valid = true
		return
	}
	i.Int64, err = strconv.ParseInt(str, 10, 64)
	i.Valid = err == nil
	return
}

// MarshalJSON implements json.Marshaler.
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Null || !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return nil, nil
	}
	if i.Null {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// SetValid changes this Int value and also sets Valid to be true.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Null = false
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null or not valid.
func (i Int) Ptr() *int64 {
	if i.Null || !i.Valid {
		return nil
	}
	return &i.Int64
}

// Scan implements database/sql.Scanner.
func (i *Int) Scan(value interface{}) (err error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i.Int64 = rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i.Int64 = int64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		i.Int64 = int64(rv.Float())
	default:
		err = fmt.Errorf("gomu: cannot scan type %T into gomu.Int: %v", value, value)
	}
	i.Valid = err == nil
	return
}
