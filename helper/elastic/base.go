package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/config"
)

type ElasticClient interface {
	Search(ctx context.Context, collection string, query map[string]interface{}) (*esapi.Response, error)
	Pagination(dataEl map[string]interface{}, page, limit int) base.Pagination
}

type elasticClient struct {
	Ec *elasticsearch.Client
}

func NewElasticClient(cfg *config.ElasticSearch) (ElasticClient, error) {
	elasticCfg := elasticsearch.Config{
		Addresses: []string{cfg.Host},
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(elasticCfg)
	if err != nil {
		return nil, err
	}

	return &elasticClient{Ec: client}, nil
}
