package transportelastic

import (
	"context"
	"marketplace-svc/app"
	endpointelastic "marketplace-svc/app/api/endpoint/elastic"
	"marketplace-svc/app/model/base/encoder"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/_struct"
	"marketplace-svc/helper/logger"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func EsProductHttpHandler(s elasticservice.ElasticProductService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsProductEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
	}

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/products/")).Handler(httpTransport.NewServer(
		ep.Search,
		decodeRequestESProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestESProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.ProductRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}

	// Default and max LIMIT
	req = req.DefaultPagination()

	return req, nil
}
