//	 Marketplace Service:
//	  version: 0.1
//	  title: Marketplace Service API
//	 Schemes: https
//	 Host:
//	 BasePath: /marketplace-svc/api/v1/
//		Consumes:
//		- application/json
//	 Produces:
//	 - application/json
//	 SecurityDefinitions:
//	  Bearer:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//	 swagger:meta
package main

import (
	"context"
	"flag"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/api/initialization"
	"marketplace-svc/app/api/middleware"
	"marketplace-svc/helper/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	_ "github.com/spf13/viper/remote"
)

func main() {
	cfg := config.Init()
	// pass to infra value
	infra := &app.Infra{
		Config: cfg,
	}
	// enable configuration
	infra.WithLog().WithDB().WithKafkaProducer().WithElasticClient()

	log := infra.Log
	log.Info("sentry: " + fmt.Sprint(viper.GetBool("sentry.is-active")))

	// sentry
	if viper.GetBool("sentry.is-active") {
		err := initialization.InitSentry()
		if err != nil {
			log.Error(err)
		}
		log.Info("Connection Sentry Success")
	}

	// Routing initialization
	mux := initialization.InitRouting(infra)
	// Wrap the entire router with gzip compression middleware
	handler := handlers.CompressHandler(mux)
	// Wrap with additional middleware if needed
	finalHandler := middleware.ServeHTTP(handler, infra.Log)

	address := flag.String("listen", ":"+strconv.Itoa(cfg.Server.Port), "Listen address.")
	httpServer := http.Server{
		Addr:              *address,
		Handler:           finalHandler,
		ReadHeaderTimeout: time.Second * cfg.Server.Timeout,
	}

	// Setup graceful shutdown
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Info("start graceful shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			panic(err)
		}
		close(idleConnectionsClosed)
	}()

	log.Info(fmt.Sprintf("Listening at port %s", strconv.Itoa(cfg.Server.Port)))
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	<-idleConnectionsClosed
}
