package encoder

import (
	"context"
	"encoding/json"
	helpervalidation "marketplace-svc/helper/validation"
	"net/http"

	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/request"
	"marketplace-svc/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

type encodeError interface {
	error() error
}

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if err, ok := resp.(encodeError); ok && err.error() != nil {
		EncodeError(ctx, err.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	result := base.GetHttpResponse(resp)
	switch result.Meta.Code {
	case message.ErrPageNotFound.Code:
		w.WriteHeader(http.StatusNotFound)
	case message.ErrNoAuth.Code, message.UnauthorizedError.Code, message.AuthenticationFailed.Code, message.UnauthorizedTokenDevice.Code:
		w.WriteHeader(http.StatusUnauthorized)
	case message.ErrDB.Code, message.ErrBadRouting.Code, message.ErrReq.Code, message.ErrNoData.Code, message.BannedLogin.Code, message.WaitingOTPRelease.Code:
		w.WriteHeader(http.StatusBadRequest)
	case message.SuccessMsg.Code:
		w.WriteHeader(http.StatusOK)
	case message.NotAllowed.Code:
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(resp)
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	msgResponse := message.ValidationError
	statusResponse := http.StatusBadRequest

	switch e := err.(type) {
	case validation.Errors, helpervalidation.ErrorsWithoutKey:
		msgResponse = message.Message{Message: e.Error(), Code: message.ValidationError.Code}
		statusResponse = http.StatusBadRequest
	case request.MalformedRequest:
		statusResponse = e.Status
		if e.Message != "" {
			msgResponse = message.Message{Message: e.Message, Code: e.Status}
		}
	default:
		statusResponse = http.StatusInternalServerError
	}

	base.ResponseWriter(w, statusResponse, base.SetErrorResponse(ctx, msgResponse, err))
}
