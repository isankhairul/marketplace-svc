package util

import "strings"

func PointToCoordinates(point string) map[string]string {
	response := map[string]string{
		"latitude":  "",
		"longitude": "",
	}
	point = strings.ReplaceAll(strings.ReplaceAll(point, "(", ""), ")", "")
	arrPoint := strings.Split(point, ",")
	if len(arrPoint) == 2 {
		response["latitude"] = arrPoint[0]
		response["longitude"] = arrPoint[1]
	}

	return response
}
