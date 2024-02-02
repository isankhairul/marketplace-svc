package helper

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"io"
	"time"
)

type HttpRC struct {
	RC     *resty.Client
	Logger logger.Logger
}

type FileReader struct {
	Param    string
	Filename string
	Reader   io.Reader
}

type RCResponse struct {
	Body     []byte          `json:"body"`
	Response *resty.Response `json:"response"`
}

type IHttpRC interface {
	Execute(
		ctx context.Context,
		method string, url string,
		headers *map[string]string,
		formData *map[string]string,
		fileReaders *[]FileReader,
		body *[]byte,
	) (*RCResponse, error)
}

func NewHttpRC(rc *resty.Client, log logger.Logger) IHttpRC {
	if rc == nil {
		rc = resty.New()
	}
	return &HttpRC{RC: rc, Logger: log}
}

func (rc HttpRC) Execute(
	ctx context.Context,
	method string, url string,
	headers *map[string]string,
	formData *map[string]string,
	fileReaders *[]FileReader,
	body *[]byte,
) (*RCResponse, error) {
	// set timeout
	httpRC := rc.RC
	httpRC.SetTimeout(3 * time.Minute)

	request := httpRC.R()

	strFormData := ""
	strBodyRequest := ""

	// set headers
	if headers != nil {
		request.SetHeaders(*headers)
	}

	// set formData
	if formData != nil {
		request.SetFormData(*formData)
		jsonFormData, _ := sonic.Marshal(*formData)
		strFormData = string(jsonFormData)
	}

	// set fileReaders
	if fileReaders != nil {
		for _, fr := range *fileReaders {
			request.SetFileReader(fr.Param, fr.Filename, fr.Reader)
		}
	}

	// set body
	if body != nil {
		request.SetBody(*body)
		strBodyRequest = string(*body)
	}

	// execute
	resp, err := request.Execute(method, url)
	responseBody := resp.Body()

	//Log CURL
	rc.Logger.WithContext(ctx).Info(fmt.Sprintf("url: %s, request: %s %s, response: %s", url, strBodyRequest, strFormData, string(responseBody)))

	return &RCResponse{Body: responseBody, Response: resp}, err
}
