package helpers

import (
	"strconv"
)

func IsBool(str string) (interface{}, bool) {
	value, err := strconv.ParseBool(str)
	return value, err == nil
}
