package base

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"marketplace-svc/helper/logger"
	"marketplace-svc/helper/message"
)

// swagger:model SuccessResponse
type responseHttp struct {
	// Meta is the API response information
	// in: struct{}
	Meta MetaResponse `json:"meta"`

	// Data is our data
	// in: struct{}
	Data data `json:"data"`
	// Errors is the response message
	// in: string
	Errors interface{} `json:"errors,omitempty"`
}

// swagger:model MetaResponse
type MetaResponse struct {
	// CorrelationId is the response correlation_id
	//in: string
	CorrelationId string `json:"correlation_id"`
	// Code is the response code
	// example: 1000
	Code int `json:"code"`
	// Message is the response message
	// example: Success
	Message string `json:"message"`

	// Pagination of to paginate response
	// in: struct{}
	Pagination *Pagination `json:"pagination,omitempty"`
}

type emptyStruct struct{}

type data struct {
	Records interface{} `json:"records,omitempty"`
	Record  interface{} `json:"record,omitempty"`
}

func SetHttpResponse(ctx context.Context, msg message.Message, result interface{}, paging *Pagination) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(result)
	isNil := false
	if result != nil {
		if resultType.Kind() == reflect.Ptr && !resultValue.IsNil() {
			// If it is a pointer, get the underlying element type
			elemType := resultType.Elem()
			elemValue := resultValue.Elem()
			// Check if the underlying type is a slice
			if elemType.Kind() == reflect.Slice {
				// The result is a pointer to a slice
				isSlice = true
				if elemValue.IsNil() || elemValue.Len() == 0 {
					isNil = true
				}
			} else {
				// The result is a pointer to a non-slice object
				isSlice = false
				if elemValue.IsValid() && elemValue.IsZero() {
					isNil = true
				}
			}
		} else if resultType.Kind() == reflect.Ptr && resultValue.IsNil() {
			isNil = true
		}
	} else {
		isNil = true
	}

	if isSlice {
		if isNil {
			result = []emptyStruct{}
		}
		dt.Records = result
		dt.Record = nil
	} else {
		if isNil {
			result = emptyStruct{}
		}
		dt.Records = nil
		dt.Record = result
	}

	return responseHttp{
		Meta: MetaResponse{
			CorrelationId: logger.GetTraceIdentifier(ctx),
			Code:          msg.Code,
			Message:       msg.Message,
			Pagination:    paging,
		},

		Data: dt,
	}
}

func GetHttpResponse(resp interface{}) *responseHttp {
	result, ok := resp.(responseHttp)
	if ok {
		return &result
	}
	return nil
}

func SetDefaultResponse(ctx context.Context, msg message.Message) interface{} {
	return responseHttp{
		Meta: MetaResponse{
			CorrelationId: logger.GetTraceIdentifier(ctx),
			Code:          msg.Code,
			Message:       msg.Message,
		},
	}
}

func SetErrorResponse(ctx context.Context, msg message.Message, errs error) interface{} {
	return responseHttp{
		Meta: MetaResponse{
			CorrelationId: logger.GetTraceIdentifier(ctx),
			Code:          msg.Code,
			Message:       msg.Message,
		},
		Errors: errs,
	}
}

func ResponseWriter(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

// swagger:response ErrorResponseBody
type ErrorResponseBody struct {
	// in: body
	Body struct {
		Meta   MetaResponse `json:"meta"`
		Errors struct {
			FieldName string `json:"field_name"`
		} `json:"errors"`
		Data struct{} `json:"data"`
	}
}
