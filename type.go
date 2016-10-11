package gomu

import (
	"regexp"
	"sync"
)

// Validator is a wrapper for functions that return bool and accept string.
type Validator func(str string) bool

// ParamValidator is a wrapper for validator functions.
type ParamValidator func(str string, params ...string) bool

// CustomTypeValidator is a wrapper for validator functions that returns bool and accept any type.
// The second parameter should be the context (in the case of validating a struct: the whole object being validated)
type CustomTypeValidator func(i interface{}, o interface{}) bool

// TagMap is a map of functions, that can be used as tags for Validate function.
var TagMap = map[string]Validator{}

type tagOptionsMap map[string]string

// ParamTagMap is a map of functions accept variants parameters.
var ParamTagMap = map[string]ParamValidator{
	"length":       ByteLength,
	"stringlength": StringLength,
}

// ParamTagRegexMap maps param tags to their respective regexes.
var ParamTagRegexMap = map[string]*regexp.Regexp{
	"length":       regexp.MustCompile("^length\\((\\d+)\\|(\\d+)\\)$"),
	"stringlength": regexp.MustCompile("^stringlength\\((\\d+)\\|(\\d+)\\)$"),
}

type customTypeTagMap struct {
	validators map[string]CustomTypeValidator

	sync.RWMutex
}

func (tm *customTypeTagMap) Get(name string) (CustomTypeValidator, bool) {
	tm.RLock()
	defer tm.RUnlock()
	v, ok := tm.validators[name]
	return v, ok
}

func (tm *customTypeTagMap) Set(name string, ctv CustomTypeValidator) {
	tm.Lock()
	defer tm.Unlock()
	tm.validators[name] = ctv
}

// CustomTypeTagMap is a map of functions that can be used as tags for Validate function.
var CustomTypeTagMap = &customTypeTagMap{validators: make(map[string]CustomTypeValidator)}
