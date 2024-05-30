package helpers

import "strings"

func IsString(str string) (interface{}, bool) {
	lowerStr := strings.TrimSpace(str)
	return lowerStr, lowerStr != ""
}
