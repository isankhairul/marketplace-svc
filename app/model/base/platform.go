package base

import "strings"

type ApiPlatform struct {
	AppVersion  string `json:"app_version"`
	Platform    string `json:"platform"`
	Brand       string `json:"brand"`
	PhoneType   string `json:"phone_type"`
	Version     string `json:"version"`
	Source      string `json:"source"`
	ScreeningID string `json:"screening_id"`
	TimeZone    string `json:"timezone"`
}

func (typePlatform *ApiPlatform) GetPLatformID() int {
	if strings.ToLower(typePlatform.Platform) == "android" {
		return 1
	} else if strings.ToLower(typePlatform.Platform) == "ios" {
		return 2
	} else {
		return 3
	}
}
