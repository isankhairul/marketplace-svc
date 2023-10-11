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

func EsMerchantHttpHandler(s elasticservice.ElasticMerchantService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsMerchantEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
	}

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/merchant/")).Handler(httpTransport.NewServer(
		ep.Search,
		decodeRequestESMerchant,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/merchant/{slug}")).Handler(httpTransport.NewServer(
		ep.Detail,
		decodeRequestESMerchantDetail,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/merchant/zipcode/{zipcode}")).Handler(httpTransport.NewServer(
		ep.SearchByZipcode,
		decodeRequestESMerchantByZipcode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeRequestESMerchant(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.MerchantRequest
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

func decodeRequestESMerchantDetail(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.MerchantDetailRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}

	// default storeID
	if req.StoreID == 0 {
		req.StoreID = 1
	}
	slug := mux.Vars(r)["slug"]
	req.Slug = slug

	return req, nil
}

func decodeRequestESMerchantByZipcode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.MerchantZipcodeRequest
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
	zipcode := mux.Vars(r)["zipcode"]
	req.Zipcode = zipcode

	return req, nil
}
