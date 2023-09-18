package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ApiRequestParams struct {
	Url         string
	Method      string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

func performRequest(apiReqParams *ApiRequestParams) (int, []byte, error) {
	client := http.Client{}

	// Method, URL, Body
	req, _ := http.NewRequest(apiReqParams.Method, apiReqParams.Url, bytes.NewBuffer(apiReqParams.Body))
	req.Header.Set("Content-Type", "application/json")

	// Headers
	for key, val := range apiReqParams.Headers {
		req.Header.Add(key, val)
	}

	// Query Params
	q := req.URL.Query()
	for key, val := range apiReqParams.QueryParams {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	// Start Request
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	// Check Response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, data, nil
}
