package helper

import (
	"io"

	"github.com/go-resty/resty/v2"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type HttpRC struct {
	RC     *resty.Client
	Logger logger.Logger
}

type FileReader struct {
	//param, fileName string, reader io.Reader
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

func (rc *HttpRC) Execute(
	method string, url string,
	headers *map[string]string,
	formData *map[string]string,
	fileReaders *[]FileReader,
	body *[]byte,
) (*RCResponse, error) {

	request := rc.RC.R()

	// set headers
	if headers != nil {
		request.SetHeaders(*headers)
	}

	// set formData
	if formData != nil {
		request.SetFormData(*formData)
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
	}

	// execute
	resp, err := request.Execute(method, url)

	return &RCResponse{Body: resp.Body(), Response: resp}, err
}
