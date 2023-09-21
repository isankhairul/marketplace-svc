package transport

import (
	"net/http"
	"os"
	"regexp"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"

	"marketplace-svc/helper/config"
)

func SwaggerHttpHandler(cfg config.URLConfig) http.Handler {
	pr := mux.NewRouter()
	basePath := cfg.BasePrefix
	basePrefix := cfg.BasePrefix

	// Handling & Manipulate swagger.yaml basePath with config-val
	pr.HandleFunc(basePath+"swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		fileBytes, err := os.ReadFile("swagger.yaml")
		if err != nil {
			panic(err)
		}

		regex, _ := regexp.Compile(`^basePath\s*:\s+.*`)
		fileBytes = regex.ReplaceAll(fileBytes, []byte("basePath: "+basePrefix))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/yaml")
		_, _ = w.Write(fileBytes)
	})
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml", BasePath: basePath}
	sh := middleware.SwaggerUI(opts, nil)
	pr.Handle(basePath+"docs", sh)

	//// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "swagger.yaml", BasePath: basePath, Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	pr.Handle(basePath+"doc", sh1)

	return pr
}
