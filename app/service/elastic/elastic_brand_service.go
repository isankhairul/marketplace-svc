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
	"strings"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type ElasticBrandService interface {
	Search(ctx context.Context, input requestelastic.BrandRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticBrandServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticBrandService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticBrandService {
	return &elasticBrandServiceImpl{config, lg, br, esc}
}

func (s elasticBrandServiceImpl) Search(_ context.Context, input requestelastic.BrandRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var brandResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg
	// set default storeID
	if input.StoreID == 0 {
		input.StoreID = 1
	}

	params := s.buildQuerySearch(input)

	indexName, err := s.getIndexName()
	if err != nil {
		return brandResponse, pagination, message.ErrNoIndexName, err
	}
	indexName = indexName + "_" + fmt.Sprint(input.StoreID)

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return brandResponse, pagination, message.ErrES, err
	}

	// requested fields
	arrFields := s.defaultFields()
	if input.Fields != "" {
		fields := util.StringExplode(input.Fields, ",")
		arrFields = append(arrFields, fields...)
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	brandResponse = s.transformSearch(responseElastic, arrFields)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return brandResponse, pagination, msg, nil
}

func (s elasticBrandServiceImpl) buildQuerySearch(input requestelastic.BrandRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  strings.Trim(input.Query, ""),
					"fields": []string{"name", "code", "principal_code"},
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
	if input.PrefixName != "" {
		filters = append(filters,
			map[string]interface{}{
				"prefix": map[string]interface{}{
					"code": strings.ToLower(input.PrefixName),
				},
			})
	}
	if input.ShowOfficial != nil {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"show_official": *input.ShowOfficial,
				},
			})
	}
	if input.PrincipalCode != "" {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"principal_code": util.StringExplode(input.PrincipalCode, ","),
				},
			})
	}
	// end filter

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

func (s elasticBrandServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"code",
		"name",
		"slug",
		"image",
	}
}

func (s elasticBrandServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-brands"]
	if indexName == nil {
		return "", errors.New("config index-brands not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticBrandServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
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
