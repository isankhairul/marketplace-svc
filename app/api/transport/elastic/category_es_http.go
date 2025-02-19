package transportelastic

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"marketplace-svc/app"
	endpointelastic "marketplace-svc/app/api/endpoint/elastic"
	"marketplace-svc/app/model/base/encoder"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/_struct"
	"marketplace-svc/helper/logger"
	"net/http"
)

func EsCategoryHttpHandler(s elasticservice.ElasticCategoryService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsCategoryEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
	}

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/categories/")).Handler(httpTransport.NewServer(
		ep.Search,
		decodeRequestESCategory,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/categories/tree")).Handler(httpTransport.NewServer(
		ep.SearchTree,
		decodeRequestESCategoryTree,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestESCategory(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.CategoryRequest
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

func decodeRequestESCategoryTree(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.CategoryTreeRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}

	// default storeID
	req = req.DefaultPagination()

	return req, nil
}
