package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/request"
	"net"
	stdHttp "net/http"
	"strings"

	"github.com/go-kit/kit/transport/http"
	"github.com/go-openapi/runtime/middleware/header"
)

func DecodeJSONBody(r *stdHttp.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			Message := "Content-Type header is not application/json"
			return request.MalformedRequest{Status: stdHttp.StatusUnsupportedMediaType, Message: Message}
		}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			Message := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}

		case errors.Is(err, io.ErrUnexpectedEOF):
			Message := "Request body contains badly-formed JSON"
			return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}

		case errors.As(err, &unmarshalTypeError):
			Message := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			Message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}

		case errors.Is(err, io.EOF):
			Message := "Request body must not be empty"
			return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}

		case err.Error() == "http: request body too large":
			Message := "Request body must not be larger than 1MB"
			return request.MalformedRequest{Status: stdHttp.StatusRequestEntityTooLarge, Message: Message}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		Message := "Request body must only contain a single JSON object"
		return request.MalformedRequest{Status: stdHttp.StatusBadRequest, Message: Message}
	}

	return nil
}

func RequestIPToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdHttp.Request) context.Context {
		return context.WithValue(ctx, base.RequestIPAddressContextKey, getIP(r))
	}
}

func getIP(r *stdHttp.Request) string {
	userIP := r.Header.Get("X-FORWARDED-FOR")
	if userIP == "" {
		userIP = r.RemoteAddr
	}

	// Handle multiple IP. ex: "66.96.247.62, 10.0.191.120"
	ipMultiple := strings.Split(userIP, ",")
	if len(ipMultiple) > 0 {
		userIP = ipMultiple[0]
	}

	ip, _, err := net.SplitHostPort(userIP)
	if err != nil {
		return userIP
	}

	// Handle localhost
	if ip == "::1" {
		return "127.0.0.1"
	}

	return ip
}
