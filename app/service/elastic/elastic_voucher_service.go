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
	promorepo "marketplace-svc/app/repository/promotion"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"
	"strconv"
	"strings"
)

type ElasticVoucherService interface {
	Search(ctx context.Context, input requestelastic.VoucherRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticVoucherServiceImpl struct {
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	voucherRepo   promorepo.PromotionScprCustomerCouponRepository
	elasticClient elastic.ElasticClient
}

func NewElasticVoucherService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	vr promorepo.PromotionScprCustomerCouponRepository,
	esc elastic.ElasticClient,
) ElasticVoucherService {
	return &elasticVoucherServiceImpl{config, lg, br, vr, esc}
}

func (s elasticVoucherServiceImpl) Search(ctx context.Context, input requestelastic.VoucherRequest) ([]map[string]interface{}, base.Pagination, message.Message, error) {
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
	bannerResponse = s.transformSearch(ctx, responseElastic, arrFields, input)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return bannerResponse, pagination, msg, nil
}

func (s elasticVoucherServiceImpl) defaultFields() []string {
	return []string{
		"id",
		"name",
		"description",
		"coupon_code",
		"category",
		"images",
		"usages",
		"from_date",
		"to_date",
		"promotion_scpr_type_id",
		"times_used",
		"uses_per_customer",
		"uses_per_coupon",
		"code",
		"merchant_ids",
		"brand",
		"payment_method",
		"category_slug",
		"claimed",
		"apply_to_claim",
		"slug",
	}
}

func (s elasticVoucherServiceImpl) buildQuerySearch(input requestelastic.VoucherRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Criteria == "category" && strings.ToLower(input.Value) == "all" {
		queryArray["bool"] = map[string]interface{}{
			"should": []map[string]interface{}{
				{
					"match": map[string]bool{
						"category.merchant": true,
					},
				},
				{
					"match": map[string]bool{
						"category.principal": true,
					},
				},
				{
					"match": map[string]bool{
						"category.payment": true,
					},
				},
			},
			"minimum_should_match": 1,
			"filter": []map[string]interface{}{
				{
					"term": map[string]interface{}{
						"status": 1,
					},
				},
				{
					"term": map[string]interface{}{
						"is_show": 1,
					},
				},
				{
					"range": map[string]interface{}{
						"from_date": map[string]string{"lte": "now"},
					},
				},
				{
					"range": map[string]interface{}{
						"to_date": map[string]string{"gte": "now"},
					},
				},
			},
		}
	} else if input.Criteria == "category" {
		queryArray["bool"] = map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"term": map[string]interface{}{"action_rule": 0},
				},
				{
					"term": map[string]interface{}{"category." + input.Value: true},
				},
				{
					"term": map[string]interface{}{"status": 1},
				},
				{
					"term": map[string]interface{}{"is_show": 1},
				},
				{
					"range": map[string]interface{}{"from_date": map[string]string{"lte": "now"}},
				},
				{
					"range": map[string]interface{}{"to_date": map[string]string{"gte": "now"}},
				},
			},
		}
	}

	if input.Criteria == "merchant" {
		queryArray["bool"] = map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"term": map[string]interface{}{"action_rule": map[string]int{"value": 0}},
				},
				{
					"term": map[string]interface{}{"coupon_type": map[string]int{"value": 2}},
				},
			},
			"filter": []map[string]interface{}{
				{
					"term": map[string]interface{}{"status": 1},
				},
				{
					"term": map[string]interface{}{"is_show": 1},
				},
				{
					"range": map[string]interface{}{"from_date": map[string]string{"lte": "now"}},
				},
				{
					"range": map[string]interface{}{"to_date": map[string]string{"gte": "now"}},
				},
				{
					"nested": map[string]interface{}{
						"path": "merchant_ids",
						"query": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": map[string]interface{}{
									"term": map[string]interface{}{
										"merchant_ids": input.Value,
									},
								},
							},
						},
					},
				},
			},
		}
	}
	querySort := map[string]string{"id": "desc"}

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

func (s elasticVoucherServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-voucher"]
	if indexName == nil {
		return "", errors.New("config index-banners not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticVoucherServiceImpl) transformSearch(ctx context.Context, rs responseelastic.SearchResponse, fields []string, input requestelastic.VoucherRequest) []map[string]interface{} {
	var response []map[string]interface{}
	var customerID uint64
	user, isValid := middleware.IsAuthContext(ctx)
	if isValid && user.CustomerID != 0 {
		customerID = user.CustomerID
	}

	for _, item := range rs.Hits.Hits {
		var tmpResponse map[string]interface{}
		var newTmpResponse = map[string]interface{}{}
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)
		if tmpResponse["image"] != nil {
			tmpResponse["image"] = s.config.URL.BaseImageURL + fmt.Sprint(tmpResponse["image"])
		}

		// add baseUrl image
		var images map[string]interface{}
		if tmpResponse["images"] != nil && len(tmpResponse["images"].([]interface{})) > 0 {
			rawImage := tmpResponse["images"].([]interface{})[0]
			images = rawImage.(map[string]interface{})
		}
		for _, key := range []string{"default", "original", "thumbnail"} {
			if images[key] != nil {
				images[key] = s.config.URL.BaseImageURL + fmt.Sprint(images[key])
			}
		}

		// get customer claimed
		dbc := repository.NewDBContext(s.baseRepo.GetDB(), ctx)
		intID, _ := strconv.Atoi(fmt.Sprint(tmpResponse["id"]))
		isClaimed := s.voucherRepo.CheckClaimed(dbc, intID, int(customerID), input.StoreID)

		// format new response
		newTmpResponse["id"] = tmpResponse["id"]
		newTmpResponse["coupon_code"] = tmpResponse["coupon_code"]
		newTmpResponse["name"] = tmpResponse["name"]
		newTmpResponse["slug"] = tmpResponse["slug"]
		newTmpResponse["description"] = tmpResponse["description"]
		newTmpResponse["merchants"] = tmpResponse["merchant_ids"]
		newTmpResponse["brand"] = tmpResponse["brand"]
		newTmpResponse["payment_method"] = tmpResponse["payment_method"]
		newTmpResponse["category_slug"] = tmpResponse["category_slug"]
		newTmpResponse["apply_to_claims"] = tmpResponse["apply_to_claims"]
		newTmpResponse["date"] = map[string]interface{}{
			"from": tmpResponse["from_date"],
			"to":   tmpResponse["to_date"],
		}
		newTmpResponse["claimed"] = isClaimed
		newTmpResponse["images"] = images

		// selected field by request
		newTmpResponseSelected := map[string]interface{}{}
		for _, field := range fields {
			if value, ok := newTmpResponse[field]; ok {
				newTmpResponseSelected[field] = value
			}
		}
		response = append(response, newTmpResponseSelected)
	}

	return response
}
