package elasticservice

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-kit/log"
	"marketplace-svc/app/repository"
)

type ElasticProductService interface {
	//GetArticle(input publicrequest.ArticlePublicRequest) (*publicresponse.ArticlePublicResponse, message.Message, interface{})
}

type elasticProductServiceImpl struct {
	logger        log.Logger
	baseRepo      repository.BaseRepository
	elasticClient *elasticsearch.Client
}

func NewElasticProductService(
	lg log.Logger,
	br repository.BaseRepository,
	esc *elasticsearch.Client,
) ElasticProductService {
	return &elasticProductServiceImpl{lg, br, esc}
}
