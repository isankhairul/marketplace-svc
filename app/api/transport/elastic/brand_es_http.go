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

func EsBrandHttpHandler(s elasticservice.ElasticBrandService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsBrandEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
	}

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/brand")).Handler(httpTransport.NewServer(
		ep.Search,
		decodeRequestESBrand,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestESBrand(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.BrandRequest
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
