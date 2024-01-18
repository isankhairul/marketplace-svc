package transportelastic

import (
	"context"
	"encoding/json"
	"errors"
	"marketplace-svc/app"
	endpointelastic "marketplace-svc/app/api/endpoint/elastic"
	"marketplace-svc/app/model/base/encoder"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/_struct"
	"marketplace-svc/helper/global"
	"marketplace-svc/helper/logger"
	"net/http"
	"strings"

	"github.com/go-kit/kit/auth/jwt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func EsMerchantHttpHandler(s elasticservice.ElasticMerchantService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsMerchantEndpoints(s, app.Config)
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

	pr.Methods(http.MethodPost).Path(app.URLWithPrefix("merchant-product/")).Handler(httpTransport.NewServer(
		ep.SearchMerchantProduct,
		decodeRequestESMerchantProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix("pharmacies/")).Handler(httpTransport.NewServer(
		ep.SearchPharmacies,
		decodeRequestESPharmacies,
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

func decodeRequestESMerchantProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.MerchantProductRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
		return nil, err
	}

	global.HtmlEscape(&req.Body)

	// Extract and set the bearer token in JwtPayload
	token, err := extractBearerToken(r)
	if err != nil {
		return nil, err
	}
	req.Token = token

	// Default and max LIMIT
	req = req.DefaultPagination()

	return req, nil
}

func decodeRequestESPharmacies(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.PharmaciesRequest
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

func extractBearerToken(r *http.Request) (string, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header not found")
	}

	// Check if the header has the Bearer prefix
	if strings.HasPrefix(authHeader, "Bearer ") {
		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return token, nil
	}

	return "", errors.New("Invalid Authorization header format")
}
