package utils

import "strconv"

func StrToBool(str string) bool {
	b, _ := strconv.ParseBool(str)
	return b
}
