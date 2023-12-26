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

type ElasticMerchantService interface {
	Search(ctx context.Context, input requestelastic.MerchantRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
	Detail(ctx context.Context, input requestelastic.MerchantDetailRequest) (map[string]interface{}, message.Message, error)
	SearchByZipcode(ctx context.Context, input requestelastic.MerchantZipcodeRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticMerchantServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticMerchantService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticMerchantService {
	return &elasticMerchantServiceImpl{config, lg, br, esc}
}

func (s elasticMerchantServiceImpl) Search(_ context.Context, input requestelastic.MerchantRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var merchantResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg
	// set default storeID
	if input.StoreID == 0 {
		input.StoreID = 1
	}

	params := s.buildQuerySearch(input)

	indexName, err := s.getIndexName()
	if err != nil {
		return merchantResponse, pagination, message.ErrNoIndexName, err
	}

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return merchantResponse, pagination, message.ErrES, err
	}

	// requested fields
	arrFields := s.defaultFields()
	if input.Fields != "" {
		fields := util.StringExplode(input.Fields, ",")
		arrFields = append(arrFields, fields...)
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	merchantResponse = s.transformSearch(responseElastic, arrFields)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return merchantResponse, pagination, msg, nil
}

func (s elasticMerchantServiceImpl) buildQuerySearch(input requestelastic.MerchantRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"match_phrase_prefix": map[string]string{
					"suggestion_terms": strings.Trim(input.Query, ""),
				},
			},
			"filter": map[string]interface{}{
				"term": map[string]int{
					"merchant_store.id": input.StoreID,
				},
			},
		}
	} else {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
			"filter": map[string]interface{}{
				"term": map[string]int{
					"merchant_store.id": input.StoreID,
				},
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
	if input.PID != "" {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"province_id": util.StringExplode(input.PID, ","),
			},
		})
	}
	if input.CID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"city_id": util.StringExplode(input.CID, ","),
			},
		})
	}
	if input.Rating != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"rating": util.StringExplode(input.Rating, ","),
			},
		})
	}
	if input.Type != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"merchant_type": input.Type,
			},
		})
	}
	if input.Store != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"merchant_store.code": util.StringExplode(input.Store, ","),
			},
		})
	}
	if input.Zipcode != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"zipcode.keyword": input.Zipcode,
			},
		})
	}
	queryArray["bool"].(map[string]interface{})["filter"] = filters
	// end filter

	// sort
	fieldSort := "created_at"
	directionSort := "asc"
	if input.Dir != "" {
		directionSort = input.Dir
	}
	if input.Sort != "" {
		switch input.Sort {
		case "name":
			fieldSort = "name.raw"
		case "rating":
			fieldSort = "rating"
		case "relevance":
			fieldSort = "_score"
		}
	}
	querySort := map[string]string{fieldSort: directionSort}

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

func (s elasticMerchantServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"code",
		"name",
		"slug",
		"status",
		"image",
		"rating",
	}
}

func (s elasticMerchantServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-merchants"]
	if indexName == nil {
		return "", errors.New("config index-merchants not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticMerchantServiceImpl) getIndexNameZone() (string, error) {
	indexName := s.config.Elastic.Index["index-merchants-zone"]
	if indexName == nil {
		return "", errors.New("config index-merchants-zone not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticMerchantServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
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

func (s elasticMerchantServiceImpl) Detail(_ context.Context, input requestelastic.MerchantDetailRequest) (map[string]interface{}, message.Message, error) {
	var merchantResponse map[string]interface{}
	msg := message.SuccessMsg

	queryArray := map[string]interface{}{}

	// create query bool
	fieldsWhere := "slug"
	var fieldValue interface{} = input.Slug

	if input.ID != 0 {
		fieldsWhere = "id"
		fieldValue = input.ID
	}
	queryArray["bool"] = map[string]interface{}{
		"must": map[string]interface{}{
			"term": map[string]interface{}{
				fieldsWhere: fieldValue,
			},
		},
	}

	params := map[string]interface{}{
		"query": queryArray,
		"size":  1,
	}

	indexName, err := s.getIndexName()
	if err != nil {
		return merchantResponse, message.ErrNoIndexName, err
	}

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return merchantResponse, message.ErrES, err
	}

	// requested fields
	arrFields := s.defaultFields()
	if input.Fields != "" {
		fields := util.StringExplode(input.Fields, ",")
		arrFields = append(arrFields, fields...)
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	arrMerchantResponse := s.transformSearch(responseElastic, arrFields)
	if len(arrMerchantResponse) == 0 {
		return merchantResponse, message.ErrNoData, errors.New(message.ErrNoData.Message)
	}

	return arrMerchantResponse[0], msg, nil
}

func (s elasticMerchantServiceImpl) SearchByZipcode(ctx context.Context, input requestelastic.MerchantZipcodeRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var merchantResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg

	// create query bool
	queryArray := map[string]interface{}{}
	queryArray["bool"] = map[string]interface{}{
		"filter": map[string]interface{}{
			"term": map[string]string{
				"zipcode": strings.Trim(input.Zipcode, ""),
			},
		},
	}
	querySort := map[string]string{"priority": "asc"}

	// pagination
	from := (input.Page - 1) * input.Limit
	params := map[string]interface{}{
		"query": queryArray,
		"from":  from,
		"size":  input.Limit,
		"sort":  querySort,
	}
	indexName, err := s.getIndexNameZone()
	if err != nil {
		return merchantResponse, pagination, message.ErrNoIndexName, err
	}

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return merchantResponse, pagination, message.ErrES, err
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	// populate to detail merchant
	if len(responseElastic.Hits.Hits) > 0 {
		for _, val := range responseElastic.Hits.Hits {
			var tmpResponse map[string]interface{}
			jsonItem, _ := json.Marshal(val.Source)
			_ = json.Unmarshal(jsonItem, &tmpResponse)
			reqDetail := requestelastic.MerchantDetailRequest{ID: int(tmpResponse["merchant_id"].(float64))}
			dSource, _, err := s.Detail(ctx, reqDetail)

			if err == nil && dSource["id"] != nil {
				dSource["priority"] = tmpResponse["priority"]
				merchantResponse = append(merchantResponse, dSource)
			}
		}
	}

	return merchantResponse, pagination, msg, nil
}
