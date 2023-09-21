package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"marketplace-svc/app/model/base"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/helper/config"
)

type ElasticClient interface {
	Search(ctx context.Context, collection string, query map[string]interface{}) (*esapi.Response, error)
	Pagination(rs responseelastic.SearchResponse, page int, limit int) base.Pagination
	BulkIndex(body []interface{}, indexName string, filename string, flush bool) error
	GetClient() *elasticsearch.Client
}

type elasticClient struct {
	Ec  *elasticsearch.Client
	Log logger.Logger
}

func (e elasticClient) GetClient() *elasticsearch.Client {
	return e.Ec
}

func (e elasticClient) GetLogger() logger.Logger {
	return e.Log
}

func NewElasticClient(cfg *config.ElasticConfig, log logger.Logger) (ElasticClient, error) {
	elasticCfg := elasticsearch.Config{
		Addresses: []string{cfg.Host},
		Username:  cfg.Username,
		Password:  cfg.Password,
		// Retry on 429 TooManyRequests statuses
		//
		RetryOnStatus: []int{502, 503, 504, 429},
		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	}

	client, err := elasticsearch.NewClient(elasticCfg)
	if err != nil {
		return nil, err
	}

	return elasticClient{Ec: client, Log: log}, nil
}
