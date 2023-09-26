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
	"strconv"
)

func EsVoucherHttpHandler(s elasticservice.ElasticVoucherService, app *app.Infra) http.Handler {
	pr := mux.NewRouter()

	ep := endpointelastic.MakeEsVoucherEndpoints(s)
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorHandler(app.Log),
		httpTransport.ServerErrorEncoder(encoder.EncodeError),
		httpTransport.ServerBefore(jwt.HTTPToContext(), logger.TraceIdentifier()),
	}

	pr.Methods(http.MethodGet).Path(app.URLWithPrefix(_struct.PrefixES + "/voucher/{criteria}/{value}/store/{storeID}")).Handler(httpTransport.NewServer(
		ep.Search,
		decodeRequestESVoucher,
		encoder.EncodeResponseHTTP,
		append(options, httpTransport.ServerBefore(jwt.HTTPToContext()))...,
	))

	return pr
}

func decodeRequestESVoucher(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req requestelastic.VoucherRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	schDecoder := schema.NewDecoder()
	schDecoder.IgnoreUnknownKeys(true)
	if err = schDecoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}

	criteria := mux.Vars(r)["criteria"]
	value := mux.Vars(r)["value"]
	storeID := 0
	if mux.Vars(r)["storeID"] != "" {
		intStoreID, _ := strconv.Atoi(mux.Vars(r)["storeID"])
		storeID = intStoreID
	}
	
	req.Criteria = criteria
	req.Value = value
	req.StoreID = storeID

	// Set Default
	req = req.DefaultPagination()

	return req, nil
}
