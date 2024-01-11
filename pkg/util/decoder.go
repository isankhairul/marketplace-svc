package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

func Base64ToMap(base64Str string) (map[string]interface{}, error) {
	arr, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(arr, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Base64ToString(base64Str string) (*string, error) {
	arr, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	strVal := string(arr)
	return &strVal, nil
}

func GetTokenFromAuth(auth string) string {
	token := strings.Split(auth, " ")
	if len(token) == 2 {
		return token[1]
	}
	return auth
}
