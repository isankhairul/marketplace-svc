package elasticservice

import (
	"context"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/message"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type mockElasticMerchantService struct{}

func (m *mockElasticMerchantService) SearchMerchantProduct(ctx context.Context, cfg *config.KalcareAPI, input requestelastic.MerchantProductRequest) ([]responseelastic.MerchantProductResponse, base.Pagination, message.Message, error) {
	return nil, base.Pagination{}, message.SuccessMsg, nil
}

var svcMerchant elasticservice.ElasticMerchantService

func init() {
	// loadConfig()
	// cfg := config.Init()
	lg, _ = logger.NewLogger(logger.NewGoKitLog(&logger.LogConfig{}), "")
	cfg := config.Config{}
	app := &app.Infra{DB: nil, Log: lg, Config: &config.Config{}}
	svcMerchant = elasticservice.NewElasticMerchantService(app, cfg, lg, baseRepo, baseElasticRepo)
	// ctx = context.Background()
}

func loadConfig() {
	// Load configuration
	viper.SetConfigType("yaml")
	var profile string = "prd"
	// if os.Getenv("KD_ENV") != "" {
	//    profile = "prd"
	// }

	var configFileName []string
	configFileName = append(configFileName, "config-")
	configFileName = append(configFileName, profile)

	viper.SetConfigName(strings.Join(configFileName, ""))
	viper.AddConfigPath("../../../../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}

	// override with env vars
	viper.AutomaticEnv()
	viper.SetEnvPrefix("KD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func TestSearchMerchantProductSuccess(t *testing.T) {
	service := &mockElasticMerchantService{}

	_, _, msg, err := service.SearchMerchantProduct(context.Background(), &config.KalcareAPI{}, requestelastic.MerchantProductRequest{})

	assert.NoError(t, err)
	assert.Equal(t, message.SuccessMsg, msg)

	// storeID := 1

	// input := requestelastic.MerchantProductRequest{
	// 	StoreID: &storeID,
	// 	Sort:    "recommendation",
	// 	Body: requestelastic.BodyMerchantProduct{
	// 		Lat: -6.258486,
	// 		Lon: 106.79229,
	// 		Items: []requestelastic.BodyReceiptsItems{
	// 			{
	// 				SKU: "T0420",
	// 				QTY: 6,
	// 			},
	// 			{
	// 				SKU: "ACLC02TAB",
	// 				QTY: 10,
	// 			},
	// 		},
	// 	},
	// }

	// aggJson := []byte(`
	// {
	// 	"group_merchant" : {
	// 		"doc_count_error_upper_bound" : 0,
	// 		"sum_other_doc_count" : 0,
	// 		"buckets" : [
	// 			{
	// 			"key" : 4514,
	// 			"doc_count" : 193,
	// 			"data" : {
	// 				"hits" : {
	// 				"total" : {
	// 					"value" : 193,
	// 					"relation" : "eq"
	// 				},
	// 				"max_score" : 9.542351,
	// 				"hits" : [
	// 					{
	// 					"_index" : "staging-merchant_flat_1",
	// 					"_type" : "_doc",
	// 					"_id" : "4514-40713",
	// 					"_score" : 9.542351,
	// 					"_source" : {
	// 						"id" : "4514-40713",
	// 						"merchant_id" : 4514,
	// 						"merchant_uid" : "611021786e18100a5b6c3",
	// 						"product_id" : 40713,
	// 						"product_sku" : "ACLC02TAB",
	// 						"merchant_sku" : "ACLC02TAB-moshkal002",
	// 						"product_status" : 1,
	// 						"stock" : 32,
	// 						"stock_on_hand" : 0,
	// 						"rating" : 0.0,
	// 						"review" : 0,
	// 						"categories" : [ ],
	// 						"max_purchase_qty" : 1,
	// 						"selling_price" : 11500,
	// 						"special_prices" : [
	// 						{
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 5
	// 						},
	// 						{
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 4
	// 						},
	// 						{
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 1
	// 						},
	// 						{
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 52
	// 						}
	// 						],
	// 						"status" : 1,
	// 						"type_id" : 2,
	// 						"updated_at" : "2023-12-20 21:18:42",
	// 						"is_pharmacy" : 1,
	// 						"location" : {
	// 						"lat" : "-6.3022015485334",
	// 						"lon" : "106.73174620979"
	// 						}
	// 					}
	// 					},
	// 					{
	// 					"_index" : "staging-merchant_flat_1",
	// 					"_type" : "_doc",
	// 					"_id" : "4514-40715",
	// 					"_score" : 9.274088,
	// 					"_source" : {
	// 						"id" : "4514-40713",
	// 						"merchant_id" : 4514,
	// 						"merchant_uid" : "611021786e18100a5b6c3",
	// 						"product_id" : 40713,
	// 						"product_sku" : "ACLC02TAB",
	// 						"merchant_sku" : "ACLC02TAB-moshkal002",
	// 						"product_status" : 1,
	// 						"stock" : 32,
	// 						"stock_on_hand" : 0,
	// 						"rating" : 0.0,
	// 						"review" : 0,
	// 						"categories" : [ ],
	// 						"max_purchase_qty" : 1,
	// 						"selling_price" : 11500,
	// 						"special_prices" : [
	// 						  {
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 5
	// 						  },
	// 						  {
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 4
	// 						  },
	// 						  {
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 1
	// 						  },
	// 						  {
	// 							"price" : 11500,
	// 							"to_time" : "2023-12-20 23:59:59",
	// 							"from_time" : "2023-12-20 00:00:00",
	// 							"customer_group_id" : 52
	// 						  }
	// 						],
	// 						"status" : 1,
	// 						"type_id" : 2,
	// 						"updated_at" : "2023-12-20 21:18:42",
	// 						"is_pharmacy" : 1,
	// 						"location" : {
	// 						  "lat" : "-6.3022015485334",
	// 						  "lon" : "106.73174620979"
	// 						}
	// 					  }
	// 					}
	// 				]
	// 				}
	// 			},
	// 			"distance" : {
	// 				"value" : 8687.867264645245
	// 			},
	// 			"product_ACLC02TAB" : {
	// 				"value" : 1.0
	// 			},
	// 			"product_T0420" : {
	// 				"value" : 1.0
	// 			},
	// 			"fulfill" : {
	// 				"value" : 2.0
	// 			}
	// 			}
	// 		]
	// 	}
	// }`)

	// var aggData map[string]interface{}
	// err := json.Unmarshal(aggJson, &aggData)

	// responseElastic := responseelastic.SearchResponse{
	// 	Took:     1,
	// 	TimedOut: false,
	// 	Shards: struct {
	// 		Total      int `json:"total"`
	// 		Successful int `json:"successful"`
	// 		Skipped    int `json:"skipped"`
	// 		Failed     int `json:"failed"`
	// 	}{
	// 		Total:      5,
	// 		Successful: 4,
	// 		Skipped:    0,
	// 		Failed:     1,
	// 	},
	// 	Hits: struct {
	// 		Total struct {
	// 			Value    int    `json:"value"`
	// 			Relation string `json:"relation"`
	// 		} `json:"total"`
	// 		MaxScore float64                              `json:"max_score"`
	// 		Hits     []responseelastic.SearchResponseHits `json:"hits"`
	// 	}{
	// 		Total: struct {
	// 			Value    int    `json:"value"`
	// 			Relation string `json:"relation"`
	// 		}{
	// 			Value:    2,
	// 			Relation: "eq",
	// 		},
	// 		MaxScore: 10.5,
	// 		Hits:     nil,
	// 	},
	// 	Aggregations: aggData,
	// }

	// responseJSON, _ := json.Marshal(responseElastic)

	// body := ioutil.NopCloser(strings.NewReader(string(responseJSON)))

	// baseElasticRepo.Mock.On("Search", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(
	// 	&esapi.Response{
	// 		StatusCode: http.StatusOK,
	// 		Body:       body,
	// 	},
	// 	nil,
	// )

	// lat := input.Body.Lat
	// lon := input.Body.Lon
	// var orderBy string

	// if input.Sort == "" {
	// 	orderBy = "recommendation"
	// } else {
	// 	orderBy = input.Sort
	// }

	// productSkus := []map[string]interface{}{}
	// for _, sku := range input.Body.Items {
	// 	skuMap := map[string]interface{}{
	// 		"sku":   sku.SKU,
	// 		"stock": sku.QTY,
	// 	}
	// 	productSkus = append(productSkus, skuMap)
	// }

	// params := buildQuerySearchMerchantProduct(input, productSkus, lat, lon, orderBy)

	// baseElasticRepo.Mock.On("Search", context.Background(), "mock.Anything", params).Return(responseElastic, nil).Once()

	// merchantProduct, _, msg, err := svcMerchant.SearchMerchantProduct(context.Background(), &config.KalcareAPI{}, input)

	// assert.NoError(t, err, "Unexpected error")
	// assert.Equal(t, message.SuccessMsg, msg)
	// assert.Equal(t, 1, len(merchantProduct))
}

func buildQuerySearchMerchantProduct(input requestelastic.MerchantProductRequest, productSkus []map[string]interface{}, lat float64, lon float64, orderBy string) map[string]interface{} {
	productVariables := make(map[string]string)

	for i, sku := range productSkus {
		productVariables[fmt.Sprintf("productSku%d", i+1)] = sku["sku"].(string)
		productVariables[fmt.Sprintf("product%d", i+1)] = fmt.Sprintf("product_%s", sku["sku"].(string))
		productVariables[fmt.Sprintf("stock%d", i+1)] = fmt.Sprintf("%d", sku["stock"].(int))
	}

	// pagination
	// from := (input.Page - 1) * input.Limit

	params := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": func() []map[string]interface{} {
					shouldClauses := make([]map[string]interface{}, len(productSkus))

					for i, sku := range productSkus {
						termClause := map[string]interface{}{
							"term": map[string]interface{}{
								"product_sku": sku["sku"],
							},
						}
						shouldClauses[i] = termClause
					}

					return shouldClauses
				}(),
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"is_pharmacy": 1,
						},
					},
					{
						"term": map[string]interface{}{
							"status": 1,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"group_merchant": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "merchant_id",
					"size":  10000,
				},
				"aggs": func() map[string]interface{} {
					aggs := make(map[string]interface{})

					for i, sku := range productSkus {
						aggs[productVariables[fmt.Sprintf("product%d", i+1)]] = map[string]interface{}{
							"sum": map[string]interface{}{
								"script": map[string]interface{}{
									"lang": "painless",
									"source": fmt.Sprintf("(doc['product_sku'].value == '%s' && doc['stock'].value >= %s) ? 1 : 0",
										sku["sku"], productVariables[fmt.Sprintf("stock%d", i+1)]),
								},
							},
						}
					}

					var scriptParts []string
					var scriptParts2 []string
					var bucketsPathKeys []string

					for _, variable := range productVariables {
						if strings.HasPrefix(variable, "product_") {
							scriptParts = append(scriptParts, fmt.Sprintf("params.%s >= 1", variable))
							scriptParts2 = append(scriptParts2, fmt.Sprintf("params.%s", variable))
							bucketsPathKeys = append(bucketsPathKeys, variable)
						}
					}

					script := strings.Join(scriptParts, " || ")
					script2 := strings.Join(scriptParts2, " + ")

					bucketsPath := make(map[string]string)
					for _, key := range bucketsPathKeys {
						bucketsPath[key] = key
					}

					aggs["filter_buckets"] = map[string]interface{}{
						"bucket_selector": map[string]interface{}{
							"buckets_path": bucketsPath,
							"script":       script,
						},
					}

					aggs["data"] = map[string]interface{}{
						"top_hits": map[string]interface{}{
							"size": len(input.Body.Items),
						},
					}

					aggs["fulfill"] = map[string]interface{}{
						"bucket_script": map[string]interface{}{
							"buckets_path": bucketsPath,
							"script":       script2,
						},
					}

					aggs["distance"] = map[string]interface{}{
						"min": map[string]interface{}{
							"script": map[string]interface{}{
								"lang":   "painless",
								"source": fmt.Sprintf("Math.round(doc['location'].arcDistance(%f, %f) * 0.001 * 10.0) / 10.0", lat, lon),
							},
						},
					}

					if orderBy == "distance" {
						aggs["sort"] = map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": []map[string]interface{}{
									{
										"distance.value": map[string]interface{}{
											"order": "asc",
										},
									},
								},
								"from": 0,
								"size": input.Limit,
							},
						}
					} else if orderBy == "fulfill" {
						aggs["sort"] = map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": []map[string]interface{}{
									{
										"fulfill.value": map[string]interface{}{
											"order": "desc",
										},
									},
								},
								"from": 0,
								"size": input.Limit,
							},
						}
					} else if orderBy == "recommendation" {
						aggs["sort"] = map[string]interface{}{
							"bucket_sort": map[string]interface{}{
								"sort": []map[string]interface{}{
									{
										"fulfill.value": map[string]interface{}{
											"order": "desc",
										},
									},
									{
										"distance.value": map[string]interface{}{
											"order": "asc",
										},
									},
								},
								"from": 0,
								"size": input.Limit,
							},
						}
					}

					return aggs
				}(),
			},
		},
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"order": "desc",
					"unit":  "km",
					"mode":  "min",
				},
			},
		},
		"size": 0,
	}

	return params
}
