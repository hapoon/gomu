package gomu

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Validate use tags for fields.
// result will be equal to `false` if there are any errors.
func Validate(s interface{}) (result bool, err error) {
	result = true
	if s == nil {
		return
	}
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		result = false
		err = fmt.Errorf("function only accepts structs; got %s", val.Kind())
		return
	}
	var errs Errors
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if typeField.PkgPath != "" {
			continue
		}
		resultField, err2 := typeCheck(valueField, typeField, val)
		if err2 != nil {
			errs = append(errs, err2)
		}
		result = result && resultField
	}
	if len(errs) > 0 {
		err = errs
	}
	return
}

func typeCheck(v reflect.Value, t reflect.StructField, o reflect.Value) (bool, error) {
	if !v.IsValid() {
		return false, nil
	}

	tag := t.Tag.Get(tagName)
	switch tag {
	case "":
		return true, nil
	case "-":
		return true, nil
	}

	options := parseTagIntoMap(tag)
	var customTypeErrors Errors
	var customTypeValidatorsExist bool
	for validatorName, customErrorMessage := range options {
		if validatefunc, ok := CustomTypeTagMap.Get(validatorName); ok {
			customTypeValidatorsExist = true
			if result := validatefunc(v.Interface(), o.Interface()); !result {
				if len(customErrorMessage) > 0 {
					customTypeErrors = append(customTypeErrors, Error{Name: t.Name, Err: fmt.Errorf(customErrorMessage), CustomErrorMessageExists: true})
					continue
				}
				customTypeErrors = append(customTypeErrors, Error{Name: t.Name, Err: fmt.Errorf("%s does not validate as %s", fmt.Sprint(v), validatorName), CustomErrorMessageExists: false})
			}
		}
	}
	if customTypeValidatorsExist {
		if len(customTypeErrors.Errors()) > 0 {
			return false, customTypeErrors
		}
		return true, nil
	}

	if isEmptyValue(v) {
		return checkRequired(v, t, options)
	}

	switch v.Type() {
	case reflect.TypeOf(String{}), reflect.TypeOf(Int{}), reflect.TypeOf(Bool{}):
		for validator, customErrorMessage := range options {
			var negate bool
			customMsgExists := (len(customErrorMessage) > 0)
			if validator[0] == '!' {
				validator = string(validator[1:])
				negate = true
			}

			for key, value := range ParamTagRegexMap {
				ps := value.FindStringSubmatch(validator)
				if len(ps) > 0 {
					if validatefunc, ok := ParamTagMap[key]; ok {
						switch v.Type() {
						case reflect.TypeOf(String{}):
							field := fmt.Sprint(v.FieldByName("String"))
							if result := validatefunc(field, ps[1:]...); (!result && !negate) || (result && negate) {
								var err error
								if !negate {
									if customMsgExists {
										err = fmt.Errorf(customErrorMessage)
									} else {
										err = fmt.Errorf("%s does not validate as %s", field, validator)
									}
								} else {
									if customMsgExists {
										err = fmt.Errorf(customErrorMessage)
									} else {
										err = fmt.Errorf("%s does validate as %s", field, validator)
									}
								}
								return false, Error{Name: t.Name, Err: err, CustomErrorMessageExists: customMsgExists}
							}
						default:
							return false, Error{t.Name, fmt.Errorf("Validator %s doesn't support type %s", validator, v.Type()), false}
						}
					}
				}
			}

			if validatefunc, ok := TagMap[validator]; ok {
				switch v.Type() {
				case reflect.TypeOf(String{}):
					field := fmt.Sprint(v.FieldByName("String"))
					if result := validatefunc(field); !result && !negate || result && negate {
						var err error
						if !negate {
							if customMsgExists {
								err = fmt.Errorf(customErrorMessage)
							} else {
								err = fmt.Errorf("%s does not validate as %s", field, validator)
							}
						} else {
							if customMsgExists {
								err = fmt.Errorf(customErrorMessage)
							} else {
								err = fmt.Errorf("%s does validate as %s", field, validator)
							}
						}
						return false, Error{t.Name, err, customMsgExists}
					}
				default:
					return false, Error{t.Name, fmt.Errorf("Validator %s doesn't support type %s", validator, v.Type()), false}
				}
			}
		}
		return true, nil
	}
	switch v.Kind() {
	case reflect.Struct:
		return Validate(v.Interface())
	default:
		return false, nil
	}
}

func parseTagIntoMap(tag string) tagOptionsMap {
	optionsMap := make(tagOptionsMap)
	options := strings.SplitN(tag, ",", -1)
	for _, option := range options {
		validationOptions := strings.Split(option, "~")
		if !isValidTag(validationOptions[0]) {
			continue
		}
		if len(validationOptions) == 2 {
			optionsMap[validationOptions[0]] = validationOptions[1]
		} else {
			optionsMap[validationOptions[0]] = ""
		}
	}
	return optionsMap
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func checkRequired(v reflect.Value, t reflect.StructField, options tagOptionsMap) (bool, error) {
	if requiredOption, isRequired := options["required"]; isRequired {
		if len(requiredOption) > 0 {
			return false, Error{t.Name, fmt.Errorf(requiredOption), true}
		}
		return false, Error{t.Name, fmt.Errorf("non zero value required"), true}
	}
	return true, nil
}

// IsURL check if the string is an URL.
func IsURL(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}

// IsRequestURL check if the string rawurl, assuming it was recieved in an HTTP request,
// is a valid URL confirm to RFC 3986.
func IsRequestURL(rawurl string) bool {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false
	}
	if len(url.Scheme) == 0 {
		return false
	}
	return true
}

// IsRequestURI check if the string rawurl, assuming it was recieved in an HTTP request,
// is an absolute URI or an absolute path.
func IsRequestURI(rawurl string) bool {
	_, err := url.ParseRequestURI(rawurl)
	return err == nil
}

// StringLength check string's length(including multi byte strings)
func StringLength(str string, params ...string) (result bool) {
	if len(params) == 2 {
		length := utf8.RuneCountInString(str)
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		result = length >= int(min) && length <= int(max)
	}
	return
}

// ByteLength check string's length.
func ByteLength(str string, params ...string) (result bool) {
	if len(params) == 2 {
		length := len(str)
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		result = length >= int(min) && length <= int(max)
	}
	return
}
