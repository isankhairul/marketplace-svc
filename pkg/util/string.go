package util

import (
	"strconv"
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringExplode(data string, delimiter string) []string {
	data = strings.Trim(data, "")
	return strings.Split(data, delimiter)
}

func StringToInt(strInt string) int {
	intStr, err := strconv.ParseInt(strInt, 10, 16)
	if err != nil {
		return 0
	}
	return int(intStr)
}

func IntToString(a []uint64) string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(int(v))
	}

	return strings.Join(b, ",")
}
