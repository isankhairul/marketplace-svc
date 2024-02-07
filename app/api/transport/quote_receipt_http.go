package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/auth/jwt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"marketplace-svc/app"
	"marketplace-svc/app/api/endpoint"
	"marketplace-svc/app/model/base/encoder"
	"marketplace-svc/app/model/request"
	"marketplace-svc/app/service"
	"marketplace-svc/helper/global"
	"marketplace-svc/helper/logger"
	"net/http"
)

func QuoteReceiptHttpHandler(s service.QuoteReceiptService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()
	ep := endpoint.MakeQuoteReceiptEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
		httpTransport.ServerBefore(global.HTTPToContextJWTClaims()),
		httpTransport.ServerBefore(global.HTTPHeaderToContext()),
	}

	pr.Methods(http.MethodPost).Path(app.URLWithPrefix("quote-receipt")).Handler(httpTransport.NewServer(
		ep.Create,
		decodeCreateQuoteReceipt,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix("quote-receipt/{quote_code}")).Handler(httpTransport.NewServer(
		ep.Find,
		decodeRequestQuoteReceipt,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodPut).Path(app.URLWithPrefix("quote-receipt/{quote_code}")).Handler(httpTransport.NewServer(
		ep.Save,
		decodeSaveQuoteReceipt,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix("quote-receipt/{quote_code}/validate")).Handler(httpTransport.NewServer(
		ep.Validate,
		decodeRequestQuoteReceiptValidate,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestQuoteReceipt(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	quotCode := mux.Vars(r)["quote_code"]
	return quotCode, nil
}

func decodeRequestQuoteReceiptValidate(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	quotCode := mux.Vars(r)["quote_code"]
	return quotCode, nil
}

func decodeSaveQuoteReceipt(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.QuoteReceiptRq
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.QuoteCode = fmt.Sprint(mux.Vars(r)["quote_code"])
	return req, nil
}

func decodeCreateQuoteReceipt(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.QuoteReceiptRq
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}
	// get body if existing
	_ = json.NewDecoder(r.Body).Decode(&req)

	return req, nil
}
