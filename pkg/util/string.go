package util

import "strings"

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
