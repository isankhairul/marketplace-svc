package initialization

import (
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"marketplace-svc/app"
	"marketplace-svc/app/api/middleware"
	"marketplace-svc/app/api/transport"
	transportelastic "marketplace-svc/app/api/transport/elastic"
	elasticregistry "marketplace-svc/app/registry/elastic"
	"marketplace-svc/helper/_struct"
	"net/http"

	"marketplace-svc/helper/cache"
	"marketplace-svc/helper/config"
)

func InitRouting(app *app.Infra) *http.ServeMux {
	// logging
	loggingMiddleware := logger.LoggingMiddleware(app.Log)

	// Elastic Service registry
	esBannerSvc := elasticregistry.RegisterEsBannerService(app)
	esBrandSvc := elasticregistry.RegisterEsBrandService(app)
	esCategorySvc := elasticregistry.RegisterEsCategoryService(app)
	esVoucherSvc := elasticregistry.RegisterEsVoucherService(app)

	//  Transport initialization
	swagHttp := transport.SwaggerHttpHandler(app.Config.URL) //don't delete or change this !!

	// Elastic Transport initialization
	esBannerHttp := transportelastic.EsBannerHttpHandler(esBannerSvc, app)
	esBrandHttp := transportelastic.EsBrandHttpHandler(esBrandSvc, app)
	esCategoryHttp := transportelastic.EsCategoryHttpHandler(esCategorySvc, app)
	esVoucherHttp := transportelastic.EsVoucherHttpHandler(esVoucherSvc, app)

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.HandleFunc("/__kdhealth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	mux.Handle(app.URLWithPrefix(_struct.PrefixES+"/banner/"), middleware.Adapt(esBannerHttp, loggingMiddleware))
	mux.Handle(app.URLWithPrefix(_struct.PrefixES+"/brand/"), middleware.Adapt(esBrandHttp, loggingMiddleware))
	mux.Handle(app.URLWithPrefix(_struct.PrefixES+"/categories/"), middleware.Adapt(esCategoryHttp, loggingMiddleware))
	mux.Handle(app.URLWithPrefix(_struct.PrefixES+"/voucher/"), middleware.Adapt(esVoucherHttp, loggingMiddleware))

	return mux
}

func InitKeyValueDatabase(cfg *config.CacheDBConfig) (cache.CacheDatabase, error) {
	return cache.SetupRedisConnection(cfg)
}
