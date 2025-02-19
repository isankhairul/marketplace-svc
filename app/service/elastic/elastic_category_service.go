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

type ElasticCategoryService interface {
	Search(ctx context.Context, input requestelastic.CategoryRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
	SearchTree(ctx context.Context, input requestelastic.CategoryTreeRequest) ([]map[string]interface{}, message.Message, error)
}

type elasticCategoryServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticCategoryService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticCategoryService {
	return &elasticCategoryServiceImpl{config, lg, br, esc}
}

func (s elasticCategoryServiceImpl) Search(_ context.Context, input requestelastic.CategoryRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var bannerResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg

	indexName, err := s.getIndexName()
	if err != nil {
		return bannerResponse, pagination, message.ErrNoIndexName, err
	}
	indexName = indexName + "_" + fmt.Sprint(input.StoreID)

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

func (s elasticCategoryServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"name",
		"slug",
		"image",
		"in_menu",
		"in_home",
		"in_homepage",
	}
}

func (s elasticCategoryServiceImpl) buildQuerySearch(input requestelastic.CategoryRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}
	queryArray["bool"] = map[string]interface{}{
		"must": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
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
	if input.Level != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"level": *input.Level,
				},
			})
	}
	if input.Position != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"position": *input.Position,
				},
			})
	}
	if input.ParentID != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"parent_id": *input.ParentID,
				},
			})
	}
	if input.InHome != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"in_home": *input.ParentID,
				},
			})
	}
	if input.InHomepage != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"in_homepage": *input.ParentID,
				},
			})
	}
	if input.InMenu != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"in_menu": *input.ParentID,
				},
			})
	}

	queryArray["bool"].(map[string]interface{})["filter"] = filters

	// pagination
	from := (input.Page - 1) * input.Limit

	params := map[string]interface{}{
		"query": queryArray,
		"from":  from,
		"size":  input.Limit,
	}

	return params
}

func (s elasticCategoryServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-category-store"]
	if indexName == nil {
		return "", errors.New("config index-category-store not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticCategoryServiceImpl) getIndexNameTree() (string, error) {
	indexName := s.config.Elastic.Index["index-category-tree"]
	if indexName == nil {
		return "", errors.New("config index-category-store not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticCategoryServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
	var searchResponse []map[string]interface{}

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
		searchResponse = append(searchResponse, tmpResponseSelected)
	}

	return searchResponse
}

func (s elasticCategoryServiceImpl) SearchTree(_ context.Context, input requestelastic.CategoryTreeRequest) ([]map[string]interface{}, message.Message, error) {
	var searchResponse []map[string]interface{}
	msg := message.SuccessMsg

	indexName, err := s.getIndexNameTree()
	if err != nil {
		return searchResponse, message.ErrNoIndexName, err
	}
	indexName = indexName + "_" + fmt.Sprint(input.StoreID)

	queryArray := map[string]interface{}{}
	queryArray["bool"] = map[string]interface{}{
		"must": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	params := map[string]interface{}{
		"query": queryArray,
	}
	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return searchResponse, message.ErrES, err
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	searchResponse = s.transformSearchTree(responseElastic)

	return searchResponse, msg, nil
}

func (s elasticCategoryServiceImpl) transformSearchTree(rs responseelastic.SearchResponse) []map[string]interface{} {
	var searchResponse []map[string]interface{}
	for _, item := range rs.Hits.Hits {
		var tmpResponse map[string]interface{}
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)
		if tmpResponse["sub"] != nil {
			sub := tmpResponse["sub"].([]interface{})
			for i := 0; i < len(sub); i++ {
				tmpSub := sub[i].(map[string]interface{})
				if tmpSub["icon"] != nil {
					tmpSub["icon"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpSub["icon"])
				}
				if tmpSub["image"] != nil {
					tmpSub["image"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpSub["image"])
				}
			}
		}

		if tmpResponse["image"] != nil {
			tmpResponse["image"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpResponse["image"])
			tmpResponse["icon"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpResponse["icon"])
		}

		searchResponse = append(searchResponse, tmpResponse)
	}

	return searchResponse
}
