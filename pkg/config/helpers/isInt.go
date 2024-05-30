package helpers

import "strconv"

func IsInt(str string) (interface{}, bool) {
	value, err := strconv.Atoi(str)
	return value, err == nil
}
