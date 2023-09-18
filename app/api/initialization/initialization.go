package initialization

import (
	"marketplace-svc/app"
	"net/http"

	"marketplace-svc/helper/cache"
	"marketplace-svc/helper/config"
)

func InitRouting(app *app.Infra) *http.ServeMux {
	// Routing path
	mux := http.NewServeMux()
	mux.HandleFunc("/__kdhealth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	return mux
}

func InitKeyValueDatabase(cfg *config.CacheDBConfig) (cache.CacheDatabase, error) {
	return cache.SetupRedisConnection(cfg)
}
