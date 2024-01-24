package elasticservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	"marketplace-svc/app/api/middleware"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/app/repository"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"
	"strings"
)

type ElasticOrderService interface {
	Search(ctx context.Context, input requestelastic.OrderRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticOrderServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticOrderService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticOrderService {
	return &elasticOrderServiceImpl{config, lg, br, esc}
}

func (s elasticOrderServiceImpl) Search(ctx context.Context, input requestelastic.OrderRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
	var brandResponse []map[string]interface{}
	var pagination base.Pagination
	msg := message.SuccessMsg
	var customerID int64
	user, isValid := middleware.IsAuthContext(ctx)
	if isValid && user.CustomerID != 0 {
		customerID = user.CustomerID
	}
	// return no auth if customer_id not found
	if customerID == 0 {
		return brandResponse, pagination, message.ErrNoAuth, errors.New(message.ErrNoAuth.Message)
	}

	input.CustomerID = fmt.Sprint(customerID)
	params := s.buildQuerySearch(input)

	indexName, err := s.getIndexName()
	if err != nil {
		return brandResponse, pagination, message.ErrNoIndexName, err
	}

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

func (s elasticOrderServiceImpl) buildQuerySearch(input requestelastic.OrderRequest) map[string]interface{} {
	// create query bool
	queryArray := map[string]interface{}{}
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"match": map[string]interface{}{
					"order_no": strings.Trim(input.OrderNo, ""),
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
				"customer.id": input.CustomerID,
			},
		},
	}

	// filter
	if input.Status != "" {
		filters = append(filters,
			map[string]interface{}{
				"terms": map[string]interface{}{
					"status_id": util.StringExplode(input.Status, ","),
				},
			})
		if fmt.Sprint(base.ORDER_STATUS_CONFIRMED) == input.Status && input.IsReviewed == 1 {
			filters = append(filters,
				map[string]interface{}{
					"terms": map[string]interface{}{
						"is_reviewed": input.IsReviewed,
					},
				})
		}
	}
	if input.PaymentMethods != "" {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"payment_method_id": util.StringExplode(input.PaymentMethods, ","),
				},
			})
	}
	if input.StoreID != 0 {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"store_id": input.StoreID,
				},
			})
	}
	if input.MerchantSlug != "" {
		filters = append(filters,
			map[string]interface{}{
				"term": map[string]interface{}{
					"merchant.slug": input.MerchantSlug,
				},
			})
	}
	// end filter

	queryArray["bool"].(map[string]interface{})["filter"] = filters
	querySort := map[string]string{"order_date": "desc"}

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

func (s elasticOrderServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"order_no",
		"order_date",
		"status_label",
		"store",
	}
}

func (s elasticOrderServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-order"]
	if indexName == nil {
		return "", errors.New("config index-order not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticOrderServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
	var response []map[string]interface{}

	for _, item := range rs.Hits.Hits {
		var tmpResponse responseelastic.OrderResponse
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)
		if len(tmpResponse.OrderItems) > 0 {
			for i := 0; i < len(tmpResponse.OrderItems); i++ {
				tmpResponse.OrderItems[i].ProductImage = s.config.URL.BaseImageURL + fmt.Sprint(tmpResponse.OrderItems[i].ProductImage)
			}
		}

		// to map string for filter selected field by request
		tmpResponseMap := map[string]interface{}{}
		jsonTmpResponse, _ := json.Marshal(tmpResponse)
		_ = json.Unmarshal(jsonTmpResponse, &tmpResponseMap)

		// selected field by request
		tmpResponseSelected := map[string]interface{}{}
		for _, field := range fields {
			if value, ok := tmpResponseMap[field]; ok {
				tmpResponseSelected[field] = value
			}
		}
		response = append(response, tmpResponseSelected)
	}

	return response
}
