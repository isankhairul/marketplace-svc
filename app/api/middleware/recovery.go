package middleware

import (
	"fmt"
	"net/http"

	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/message"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// func ServeHTTP(h http.Handler, lg logger.Logger) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				base.ResponseWriter(w, http.StatusInternalServerError, base.SetDefaultResponse(req.Context(), message.FailedMsg))
// 				lg.Error(fmt.Errorf("%s", err))
// 			}
// 		}()
// 		h.ServeHTTP(w, req)
// 	})
// }

func ServeHTTP(h http.Handler, lg logger.Logger) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { // this one
		defer func() {
			if err := recover(); err != nil {
				base.ResponseWriter(w, http.StatusInternalServerError, base.SetDefaultResponse(req.Context(), message.FailedMsg))
				lg.Error(fmt.Errorf("%s", err))
			}
		}()
		h.ServeHTTP(w, req)
	}), "Klikmedis-Service")
}
