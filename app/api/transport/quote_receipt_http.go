package transport

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"marketplace-svc/app"
	"marketplace-svc/app/api/endpoint"
	"marketplace-svc/app/model/base/encoder"
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

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix("quote-receipt/{quote_code}")).Handler(httpTransport.NewServer(
		ep.Find,
		decodeRequestQuoteReceipt,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestQuoteReceipt(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	quotCode := mux.Vars(r)["quote_code"]
	return quotCode, nil
}
