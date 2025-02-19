package elasticservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/app/repository"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type ElasticBannerService interface {
	Search(ctx context.Context, input requestelastic.BannerRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticBannerServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticBannerService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticBannerService {
	return &elasticBannerServiceImpl{config, lg, br, esc}
}

func (s elasticBannerServiceImpl) Search(ctx context.Context, input requestelastic.BannerRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var bannerResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg

	indexName, err := s.getIndexName()
	if err != nil {
		return bannerResponse, pagination, message.ErrNoIndexName, err
	}

	params := s.buildQuerySearch(input)
	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return bannerResponse, pagination, message.ErrES, err
	}

	// requested fields
	arrFields := s.defaultFields()
	if input.Fields != "" {
		fields := util.StringExplode(input.Fields, ",")
		arrFields = append(arrFields, fields...)
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	bannerResponse = s.transformSearch(responseElastic, arrFields)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return bannerResponse, pagination, msg, nil
}

func (s elasticBannerServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"title",
		"slug",
		"image",
	}
}

func (s elasticBannerServiceImpl) buildQuerySearch(input requestelastic.BannerRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": []interface{}{
				map[string]interface{}{
					"match_phrase": map[string]interface{}{
						"title": map[string]interface{}{
							"query": input.Query,
						},
					},
				},
			},
		}
	} else {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}

	// default filter status
	filters := []map[string]interface{}{
		{
			"term": map[string]interface{}{
				"status": 1,
			},
		},
	}

	// filter
	if input.CategorySlug != "" {
		filterCategorySlug := map[string]interface{}{
			"term": map[string]interface{}{
				"category_slug": input.CategorySlug,
			},
		}
		filters = append(filters, filterCategorySlug)
	}

	queryArray["bool"].(map[string]interface{})["filter"] = filters
	querySort := map[string]string{"created_at": "desc"}

	// pagination
	from := (input.Page - 1) * input.Limit

	params := map[string]interface{}{
		"query": queryArray,
		"from":  from,
		"size":  input.Limit,
		"sort":  querySort,
	}

	return params
}

func (s elasticBannerServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-banners"]
	if indexName == nil {
		return "", errors.New("config index-banners not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticBannerServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
	var response []map[string]interface{}

	for _, item := range rs.Hits.Hits {
		var tmpResponse map[string]interface{}
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)
		if tmpResponse["image"] != nil {
			tmpResponse["image"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpResponse["image"])
		}

		// selected field by request
		tmpResponseSelected := map[string]interface{}{}
		for _, field := range fields {
			if value, ok := tmpResponse[field]; ok {
				tmpResponseSelected[field] = value
			}
		}
		response = append(response, tmpResponseSelected)
	}

	return response
}
