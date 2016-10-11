package gomu

import "strconv"

// ToInt convert the input string to an integer, or 0 if the input is not an integer.
func ToInt(str string) (result int64, err error) {
	result, err = strconv.ParseInt(str, 0, 64)
	if err != nil {
		result = 0
	}
	return
}
