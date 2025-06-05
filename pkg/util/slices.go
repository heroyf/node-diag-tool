package util

import "github.com/samber/lo"

// InSlice 判断字符串是否在slice中
func InSlice(slices []string, str string) bool {
	return lo.Contains(slices, str)
}

func Contains(slices []string, str string) bool {
	return InSlice(slices, str)
}
