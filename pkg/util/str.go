package util

import (
	"github.com/samber/lo"
	"strconv"
)

func Substr(str string, length uint) string {
	return lo.Substring(str, 0, length)
}

func IsEmpty(str string) bool {
	return lo.IsEmpty(str)
}

func IsNotEmpty(str string) bool {
	return lo.IsNotEmpty(str)
}

func ToInt(str string) int {
	atoi, _ := strconv.Atoi(str)
	return atoi
}
