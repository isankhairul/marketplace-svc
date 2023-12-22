package elasticservice

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace-svc/app/model/base"
	modelelastic "marketplace-svc/app/model/elastic"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/app/repository"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"
	"math"
	"strconv"
	"time"

	"github.com/alitto/pond"
	"github.com/pkg/errors"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type ElasticProductService interface {
	Reindex(index string, productIDs string, storeID int, productType string, merchantID int, flush bool) error
	Search(ctx context.Context, input requestelastic.ProductRequest) ([]responseelastic.ProductResponse, base.Pagination, message.Message, error)
}

type elasticProductServiceImpl struct {
	config          config.Config
	logger          logger.Logger
	baseRepo        repository.BaseRepository
	productFlatRepo repository.ProductFlatRepository
	elasticClient   elastic.ElasticClient
}

func NewElasticProductService(
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	pfr repository.ProductFlatRepository,
	esc elastic.ElasticClient,
) ElasticProductService {
	return &elasticProductServiceImpl{config, lg, br, pfr, esc}
}

func (s elasticProductServiceImpl) Reindex(index string, productIDs string, storeID int, productType string, merchantID int, flush bool) error {
	if index == "" {
		index = "golang_product_store_" + fmt.Sprint(storeID)
	}

	titleLog := "[ES-PRODUCT-REINDEX-STORE-" + fmt.Sprint(storeID) + "]"
	sqlCount := s.productFlatRepo.GetQueryRawAllProduct(productIDs, storeID, productType, true, merchantID)

	var total int
	s.baseRepo.GetDB().Raw(sqlCount).Scan(&total)

	if total == 0 {
		s.logger.Info(titleLog + " EMPTY PRODUCTS")
		return errors.New("empty products")
	}

	chunk := 1000
	totalPage := int(math.Ceil(float64(total) / float64(chunk)))

	var sqlProduct = s.productFlatRepo.GetQueryRawAllProduct(productIDs, storeID, productType, false, merchantID)
	var arrBody []interface{}
	var chanBody = make(chan interface{})

	// This will create a pool with 20 running worker goroutines
	pool := pond.New(20, totalPage, pond.IdleTimeout(30*time.Second))

	for page := 1; page <= totalPage; page++ {
		n := page
		pool.Submit(func() {
			s.logger.Info(titleLog + "GET PRODUCTS - PAGE " + fmt.Sprint(n) + " OF " + fmt.Sprint(totalPage))
			offset := (n - 1) * chunk
			var values []map[string]interface{}
			_ = s.baseRepo.GetDB().Raw(sqlProduct + " LIMIT " + fmt.Sprint(chunk) + " OFFSET " + fmt.Sprint(offset)).Scan(&values)
			for _, value := range values {
				pid := fmt.Sprint(value["product_id"]) + "-" + value["merchant_uid"].(string)
				intVar, _ := strconv.ParseInt(fmt.Sprint(value["product_is_active"]), 0, 32)
				currentBody := modelelastic.EsProductFlat{
					ID:              pid,
					Name:            fmt.Sprint(value["product_name"]),
					Sku:             fmt.Sprint(value["product_sku"]),
					Slug:            fmt.Sprint(value["product_slug"]),
					BrandCode:       fmt.Sprint(value["product_brand_code"]),
					MetaTitle:       fmt.Sprint(value["product_meta_title"]),
					MetaDescription: fmt.Sprint(value["product_meta_description"]),
					IsActive:        int32(intVar),
					TypeID:          fmt.Sprint(value["product_type_id"]),
					Status:          value["product_status"].(int32),
					CreatedAt:       "2023-09-19 16:54:20",
					UpdatedAt:       "2023-09-19 16:54:20",
				}
				currentBody.Merchants = s.getMerchants(value)

				chanBody <- currentBody
			}
		})
	}

	// get chanBody
	for i := 0; i < total; i++ {
		item, ok := <-chanBody
		if ok {
			arrBody = append(arrBody, item)
		}
	}
	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()
	close(chanBody)

	_ = s.elasticClient.BulkIndex(arrBody, index, "product_store", flush)

	return nil
}

func (s elasticProductServiceImpl) getMerchants(item map[string]interface{}) modelelastic.EsProductFlatMerchant {
	var result = modelelastic.EsProductFlatMerchant{
		ID:             item["id"].(int32),
		UID:            item["merchant_uid"].(string),
		Code:           item["code"].(string),
		Name:           item["name"].(string),
		Slug:           item["slug"].(string),
		Status:         item["status"].(int32),
		TypeID:         item["type_id"].(int32),
		TypeSlug:       item["type_slug"].(string),
		Stock:          item["stock"].(int32),
		ReservedStock:  item["reserved_stock"].(int32),
		StockOnHand:    item["stock_on_hand"].(int32),
		MaxPurchaseQty: item["max_purchase_qty"].(int32),
		Sold:           item["sold"].(int32),
		ProvinceID:     item["province_id"].(int32),
		Province:       item["province"].(string),
		CityID:         item["city_id"].(int32),
		DistrictID:     item["district_id"].(int32),
		District:       item["district"].(string),
		SubdistrictID:  item["subdistrict_id"].(int32),
		Subdistrict:    item["subdistrict"].(string),
		PostalcodeID:   item["postalcode_id"].(int32),
		Zipcode:        item["zipcode"].(string),
		Location: modelelastic.EsProductFlatLocation{
			Lat: fmt.Sprint(item["latitude"]),
			Lon: fmt.Sprint(item["longitude"]),
		},
		Image:             item["image"].(string),
		Categories:        []string{},
		MerchantProductID: item["merchant_product_id"].(int32),
		MerchantSku:       item["merchant_sku"].(string),
		ProductStatus:     item["product_status"].(int32),
		Rating:            math.Round(item["rating"].(float64)),
		Review:            item["review"].(int32),
		SellingPrice:      item["selling_price"].(int32),
		HidePrice:         item["hide_price"].(bool),
	}
	var sp []modelelastic.EsProductFlatSpecialPrice
	var epfsp []map[string]interface{}
	json.Unmarshal([]byte(item["special_prices"].(string)), &epfsp)
	for _, val := range epfsp {
		sp = append(sp, modelelastic.EsProductFlatSpecialPrice{
			CustomerGroupID: int32(val["customer_group_id"].(float64)),
			Price:           int32(val["price"].(float64)),
			FromTime:        val["from_time"].(string),
			ToTime:          val["to_time"].(string),
		})
	}
	result.SpecialPrices = sp

	return result
}

// swagger:operation GET /products Products ProductRequest
// Product - List
//
// ---
// tags:
//   - "Products"
//
// security:
// - Bearer: []
// responses:
//
//	  '200':
//		   description: Product - List success response
//		   schema:
//		       properties:
//		           meta:
//		               $ref: '#/definitions/MetaResponse'
//		           data:
//		               type: object
//		               properties:
//		                   records:
//		                       type: array
//		                       items:
//		                           $ref: '#/definitions/ProductResponse'
func (s elasticProductServiceImpl) Search(_ context.Context, input requestelastic.ProductRequest) ([]responseelastic.ProductResponse, base.Pagination, message.Message, error) {
	var productResponse []map[string]interface{}
	var newProductResponse []responseelastic.ProductResponse
	var pagination base.Pagination
	msg := message.SuccessMsg

	if input.StoreID == nil {
		storeID := s.config.Elastic.DefaultStoreID
		input.StoreID = &storeID
	}

	indexName, err := s.getIndexName(*input.StoreID)
	if err != nil {
		return newProductResponse, pagination, message.ErrNoIndexName, err
	}

	params := s.buildQuerySearch(input)
	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return newProductResponse, pagination, message.ErrES, err
	}

	// requested fields
	arrFields := s.defaultFields()
	if input.Fields != "" {
		fields := util.StringExplode(input.Fields, ",")
		arrFields = append(arrFields, fields...)
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

	productResponse = s.transformSearch(responseElastic, arrFields, *input.StoreID)
	newProductResponse = s.transformResponse(productResponse)

	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return newProductResponse, pagination, msg, nil
}

func (s elasticProductServiceImpl) getIndexName(storeID int) (string, error) {
	indexName := s.config.Elastic.Index["index-products-flat"]
	if indexName == nil {
		return "", errors.New("config index-products-flat not defined")
	}

	return fmt.Sprint(indexName, "_", storeID), nil
}

func (s elasticProductServiceImpl) getIndexNameMerchant(storeID int) (string, error) {
	indexName := s.config.Elastic.Index["index-merchants-flat"]
	if indexName == nil {
		return "", errors.New("config index-merchants-flat not defined")
	}

	return fmt.Sprint(indexName, "_", storeID), nil
}

func (s elasticProductServiceImpl) buildQuerySearch(input requestelastic.ProductRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  input.Query,
					"fields": []string{"name"},
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
				"is_pharmacy": 1,
			},
		},
		{
			"term": map[string]interface{}{
				"is_active": 1,
			},
		},
		{
			"term": map[string]interface{}{
				"status": 1,
			},
		},
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

func (s elasticProductServiceImpl) buildQuerySearchMerchant(productID float64) map[string]interface{} {
	customerGroup := s.config.Elastic.DefaultCustomerGroup
	aggsSize := s.config.Elastic.AggsSize
	maxLimit := s.config.Elastic.MaxLimit

	params := map[string]interface{}{
		"size": maxLimit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"product_id": productID,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"is_pharmacy": 1,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": 1,
						},
					},
					map[string]interface{}{
						"terms": map[string]interface{}{
							"type_id": []int{1, 2},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"product_status": 1,
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"stock": map[string]interface{}{
								"gt": 0,
							},
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"selling_price": map[string]interface{}{
								"gt": 0,
							},
						},
					},
					map[string]interface{}{
						"nested": map[string]interface{}{
							"path": "special_prices",
							"query": map[string]interface{}{
								"bool": map[string]interface{}{
									"must": []interface{}{
										map[string]interface{}{
											"term": map[string]interface{}{
												"special_prices.customer_group_id": customerGroup,
											},
										},
										map[string]interface{}{
											"range": map[string]interface{}{
												"special_prices.price": map[string]interface{}{
													"gt": 0,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"sell_prices": map[string]interface{}{
				"min": map[string]interface{}{
					"field": "selling_price",
				},
			},
			"spec_prices": map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "special_prices",
				},
				"aggs": map[string]interface{}{
					"filtered": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"special_prices.customer_group_id": customerGroup,
							},
						},
						"aggs": map[string]interface{}{
							"group_by": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "special_prices.customer_group_id",
									"size":  aggsSize,
								},
								"aggs": map[string]interface{}{
									"min_price": map[string]interface{}{
										"min": map[string]interface{}{
											"field": "special_prices.price",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return params
}

func (s elasticProductServiceImpl) defaultFields() []string {
	return []string{
		"sku",
		"name",
		"uom",
		"uom_name",
		"weight",
		// "description",
		// "short_description",
		"images",
		// "principal_name",
		"price",
		"min_price",
		"proportional",
		"pharmacy_code",
	}
}

func (s elasticProductServiceImpl) transformResponse(response []map[string]interface{}) []responseelastic.ProductResponse {
	resp := []responseelastic.ProductResponse{}
	for _, val := range response {
		data := responseelastic.NewProductResponse(val)
		resp = append(resp, *data)
	}

	return resp
}

func (s elasticProductServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string, storeID int) []map[string]interface{} {
	var response []map[string]interface{}

	for _, item := range rs.Hits.Hits {
		var tmpResponse map[string]interface{}
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)

		// query to merchant flag
		indexName, err := s.getIndexNameMerchant(storeID)
		if err != nil {
			return response
		}

		productID := tmpResponse["id"].(float64)

		params := s.buildQuerySearchMerchant(productID)
		resp, err := s.elasticClient.Search(context.Background(), indexName, params)
		if err != nil {
			s.logger.Error(errors.New("error request elastic: " + err.Error()))
			return response
		}

		var responseElastic responseelastic.SearchResponse
		_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

		// sellingPrice := responseElastic.Aggregations.(map[string]interface{})["sell_prices"].(map[string]interface{})["value"].(float64)
		// specialPrice := responseElastic.Aggregations.(map[string]interface{})["spec_prices"].(map[string]interface{})["filtered"].(map[string]interface{})["group_by"].(map[string]interface{})["buckets"].([]interface{})[0].(map[string]interface{})["min_price"].(map[string]interface{})["value"].(float64)
		var sellingPrice float64
		if aggs, ok := responseElastic.Aggregations.(map[string]interface{}); ok {
			if sellPrices, ok := aggs["sell_prices"].(map[string]interface{}); ok {
				sellingPrice = sellPrices["value"].(float64)
			}
		}

		var specialPrice float64
		if aggs, ok := responseElastic.Aggregations.(map[string]interface{}); ok {
			if specPrices, ok := aggs["spec_prices"].(map[string]interface{}); ok {
				if filtered, ok := specPrices["filtered"].(map[string]interface{}); ok {
					if groupBy, ok := filtered["group_by"].(map[string]interface{}); ok {
						if buckets, ok := groupBy["buckets"].([]interface{}); ok && len(buckets) > 0 {
							if minPrice, ok := buckets[0].(map[string]interface{}); ok {
								specialPrice = minPrice["min_price"].(map[string]interface{})["value"].(float64)
							}
						}
					}
				}
			}
		}

		// selected field by request
		tmpResponseSelected := map[string]interface{}{}
		for _, field := range fields {
			if value, ok := tmpResponse[field]; ok {
				tmpResponseSelected[field] = value
			}
			if field == "price" {
				tmpResponseSelected[field] = sellingPrice
			} else if field == "min_price" {
				tmpResponseSelected[field] = specialPrice
			}
		}
		response = append(response, tmpResponseSelected)
	}

	return response
}
