package app

import (
	"context"
	"fmt"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/queue"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"marketplace-svc/helper/cache"
	"marketplace-svc/helper/config"
	db "marketplace-svc/helper/database"
)

type Infra struct {
	Config        *config.Config
	Log           logger.Logger
	DB            *db.Database
	Redis         cache.CacheDatabase
	KafkaProducer *queue.KafkaBase
	ElasticClient elastic.ElasticClient
}

func (app *Infra) URLWithPrefix(url string) string {
	return fmt.Sprintf("%s%s", app.Config.BasePrefix(), url)
}

func (app *Infra) LogWithContext(ctx context.Context) logger.Logger {
	return app.Log.WithContext(ctx)
}

func (app *Infra) EnableLog() *Infra {
	log, err := logger.NewLogger(
		logger.NewGoKitLog(&logger.LogConfig{
			Level: app.Config.Server.LogConfig.Level,
		}), "Marketplace Service",
	)
	if err != nil {
		log.Error(err)
		panic(err.Error())
	}
	app.Log = log
	return app
}

func (app *Infra) EnableDB() *Infra {
	//Init DB Connection
	migration := &db.Option{Migrate: []interface{}{}}
	initDB, err := db.NewDatabaseConnect(&app.Config.DB, migration)
	if err != nil {
		panic(err.Error())
	}
	app.DB = &initDB
	return app
}

func (app *Infra) EnableKafkaProducer() *Infra {
	kafkaProducer, err := queue.NewKafkaProducer(*app.Config)
	if err != nil {
		panic(err.Error())
	}
	app.KafkaProducer = kafkaProducer
	return app
}

func (app *Infra) EnableElasticClient() *Infra {
	ec, err := elastic.NewElasticClient(&app.Config.ElasticSearch)
	if err != nil {
		panic(err.Error())
	}
	app.ElasticClient = ec
	return app
}
