package elasticservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/app/repository"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/elastic"
	"marketplace-svc/helper/global"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"
	"strconv"
	"strings"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

type ElasticMerchantService interface {
	Search(ctx context.Context, input requestelastic.MerchantRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
	Detail(ctx context.Context, input requestelastic.MerchantDetailRequest) (map[string]interface{}, message.Message, error)
	SearchByZipcode(ctx context.Context, input requestelastic.MerchantZipcodeRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
	SearchMerchantProduct(ctx context.Context, cfg *config.KalcareAPI, input requestelastic.MerchantProductRequest) ([]responseelastic.MerchantProductResponse, base.Pagination, message.Message, error)
}

type elasticMerchantServiceImpl struct {
	app           *app.Infra
	config        config.Config
	logger        logger.Logger
	baseRepo      repository.BaseRepository
	elasticClient elastic.ElasticClient
}

func NewElasticMerchantService(
	app *app.Infra,
	config config.Config,
	lg logger.Logger,
	br repository.BaseRepository,
	esc elastic.ElasticClient,
) ElasticMerchantService {
	return &elasticMerchantServiceImpl{app, config, lg, br, esc}
}

func (s elasticMerchantServiceImpl) SearchMerchantProduct(ctx context.Context, cfg *config.KalcareAPI, input requestelastic.MerchantProductRequest) ([]responseelastic.MerchantProductResponse, base.Pagination, message.Message, error) {
	log := s.app.LogWithContext(ctx)

	var newMerchantProductResponse []responseelastic.MerchantProductResponse
	var pagination base.Pagination
	msg := message.SuccessMsg

	// validate payload
	if err := input.Validate(); err != nil {
		return newMerchantProductResponse, pagination, message.ErrInvalidReqBody, err
	}

	if input.StoreID == nil {
		storeID := s.config.Elastic.DefaultStoreID
		input.StoreID = &storeID
	}

	indexName, err := s.getIndexFlatName(*input.StoreID)
	if err != nil {
		return newMerchantProductResponse, pagination, message.ErrNoIndexName, err
	}

	// hit customer info marketplace
	// customerInfo, err := communicates.GetCustomerInfo(log, cfg, input.Token)
	// if err != nil {
	// 	log.Error(err)
	// 	return newMerchantProductResponse, pagination, message.FailedMsg, err
	// }
	// if customerInfo.Message != "" {
	// 	return newMerchantProductResponse, pagination, message.Message{Code: 400, Message: customerInfo.Message}, err
	// }
	// if customerInfo.Data.Record.Address == nil || customerInfo.Data.Record.Address.Latitude == nil || customerInfo.Data.Record.Address.Longitude == nil {
	// 	return newMerchantProductResponse, pagination, message.AddressNotFound, err
	// }

	lat := input.Body.Lat
	lon := input.Body.Lon
	var orderBy string

	if input.Sort == "" {
		orderBy = "recommendation"
	} else {
		orderBy = input.Sort
	}

	if orderBy != "distance" && orderBy != "fulfill" && orderBy != "recommendation" {
		return newMerchantProductResponse, pagination, message.MerchantProductSearchSortNotFound, err
	}

	productSkus := []map[string]interface{}{}
	for _, sku := range input.Body.Items {
		skuMap := map[string]interface{}{
			"sku":   sku.SKU,
			"stock": sku.QTY,
		}
		productSkus = append(productSkus, skuMap)
	}

	params := s.buildQuerySearchMerchantProduct(input, productSkus, lat, lon, orderBy)
	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		log.Error(errors.New("error request elastic: " + err.Error()))
		return newMerchantProductResponse, pagination, message.ErrES, err
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

	newMerchantProductResponse = s.transformSearchMerchantProduct(responseElastic, productSkus, lat, lon, *input.StoreID)

	// pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return newMerchantProductResponse, pagination, msg, nil

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

func (s elasticMerchantServiceImpl) buildQuerySearchMerchantProduct(input requestelastic.MerchantProductRequest, productSkus []map[string]interface{}, lat float64, lon float64, orderBy string) map[string]interface{} {
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

func (s elasticMerchantServiceImpl) buildQuerySearchMerchant(merchantID float64) map[string]interface{} {
	queryArray := map[string]interface{}{}
	// create query bool
	queryArray["bool"] = map[string]interface{}{
		"must": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	// default filter status
	filters := []map[string]interface{}{
		{
			"term": map[string]interface{}{
				"id": merchantID,
			},
		},
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
	}

	queryArray["bool"].(map[string]interface{})["filter"] = filters

	params := map[string]interface{}{
		"query": queryArray,
	}

	return params
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

func (s elasticMerchantServiceImpl) getIndexFlatName(storeID int) (string, error) {
	indexName := s.config.Elastic.Index["index-merchants-flat"]
	if indexName == nil {
		return "", errors.New("config index-merchants-flat not defined")
	}

	return fmt.Sprint(indexName, "_", storeID), nil
}

func (s elasticMerchantServiceImpl) getIndexFlatNameProduct(storeID int) (string, error) {
	indexName := s.config.Elastic.Index["index-products-flat"]
	if indexName == nil {
		return "", errors.New("config index-products-flat not defined")
	}

	return fmt.Sprint(indexName, "_", storeID), nil
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

func (s elasticMerchantServiceImpl) transformSearchMerchantProduct(rs responseelastic.SearchResponse, productSkus []map[string]interface{}, lat1 float64, lon1 float64, storeID int) []responseelastic.MerchantProductResponse {
	var allMerchantResponses []responseelastic.MerchantProductResponse
	// var response responseelastic.MerchantProductResponse
	var response2 map[string]interface{}

	buckets, _ := rs.Aggregations.(map[string]interface{})["group_merchant"].(map[string]interface{})["buckets"].([]interface{})

	var merchantID float64
	var distance float64
	var fulfill float64
	var totalProduct float64

	if len(buckets) > 0 {

		for _, buck := range buckets {
			// firstBucket, ok := buckets[0].(map[string]interface{})
			buck := buck.(map[string]interface{})
			merchantID = buck["key"].(float64)
			distance = buck["distance"].(map[string]interface{})["value"].(float64)
			fulfill = buck["fulfill"].(map[string]interface{})["value"].(float64)

			var productAvailable []responseelastic.ProductsAvailable
			var productOrdered []responseelastic.ProductsOrdered

			for _, sku := range productSkus {
				productSkuName := "product_" + sku["sku"].(string)
				available := buck[productSkuName].(map[string]interface{})["value"].(float64)
				productOrdered = append(productOrdered, responseelastic.ProductsOrdered{
					SKU:       sku["sku"].(string),
					QTY:       sku["stock"].(int),
					Available: available,
				})
			}

			bucket_hits, ok := buck["data"].(map[string]interface{})
			if ok {
				hits, ok := bucket_hits["hits"].(map[string]interface{})
				if ok {
					total, ok := hits["total"].(map[string]interface{})
					if ok {
						// total jumlah product yang dijual
						value, _ := total["value"].(float64)
						totalProduct = value
					}
					hits2, ok := hits["hits"].([]interface{})
					if ok {
						// productIDs yang jual
						for _, item := range hits2 {
							value, _ := item.(map[string]interface{})["_source"].(map[string]interface{})
							productSKU := value["product_sku"].(string)
							merchantSKU := value["merchant_sku"].(string)
							productStock := value["stock"].(float64)
							// add selling_price and special_price
							price := s.GetPriceProduct(productSKU, storeID)
							// add uom and uom_name
							productDetail := s.GetProductDetail(productSKU, storeID)
							productAvailable = append(productAvailable, responseelastic.ProductsAvailable{
								SKU:          productSKU,
								MerchantSKU:  merchantSKU,
								Name:         productDetail["name"].(string),
								QTY:          productStock,
								UOM:          productDetail["uom"].(string),
								UOMName:      productDetail["uom_name"].(string),
								Image:        productDetail["image"].(string),
								SellingPrice: price[0],
								SpecialPrice: price[1],
							})
						}
					}
				}
			}

			for i := range productOrdered {
				orderedSKU := productOrdered[i].SKU
				orderedAvailable := productOrdered[i].Available

				// Check if the product is available
				found := false
				for j := range productAvailable {
					if productAvailable[j].SKU == orderedSKU {
						found = true
						//added selling_price and special_price in prduct_ordered
						productOrdered[i].SellingPrice = productAvailable[j].SellingPrice
						productOrdered[i].SpecialPrice = productAvailable[j].SpecialPrice
						//added name, uom and uom_name in product_ordered
						productOrdered[i].Name = productAvailable[j].Name
						productOrdered[i].UOM = productAvailable[j].UOM
						productOrdered[i].UOMName = productAvailable[j].UOMName
						//added image in product_ordered
						productOrdered[i].Image = productAvailable[j].Image
						//added stock merchant
						productOrdered[i].QTYAvailable = productAvailable[j].QTY
						//added merchant_sku
						productOrdered[i].MerchantSKU = productAvailable[j].MerchantSKU
						break
					}
				}

				if found && orderedAvailable == 1 {
					productOrdered[i].Status = "available"
					productOrdered[i].IsAvailable = true
				} else if found && orderedAvailable == 0 {
					productOrdered[i].Status = "oos"
					productOrdered[i].IsAvailable = false

				} else {
					productOrdered[i].Status = "notavailable"
				}
			}

			indexName, _ := s.getIndexName()
			params := s.buildQuerySearchMerchant(merchantID)
			resp, err := s.elasticClient.Search(context.Background(), indexName, params)
			if err != nil {
				s.logger.Error(errors.New("error request elastic: " + err.Error()))
				return allMerchantResponses
			}

			var responseElastic responseelastic.SearchResponse
			_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

			merchantID := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["id"]
			merchantUID := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["uid"]
			merchantName := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["name"]
			lat2, _ := strconv.ParseFloat(responseElastic.Hits.Hits[0].Source.(map[string]interface{})["location"].(map[string]interface{})["lat"].(string), 64)
			lon2, _ := strconv.ParseFloat(responseElastic.Hits.Hits[0].Source.(map[string]interface{})["location"].(map[string]interface{})["lon"].(string), 64)
			shippingDuration := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["shipping_duration"].([]interface{})

			if distance > 30 {
				shippingDuration = []interface{}{"Reguler"}
			}

			// calculate distance
			coordinates1 := responseelastic.Coordinates{Lat: lat1, Lon: lon1}
			coordinates2 := responseelastic.Coordinates{Lat: lat2, Lon: lon2}

			response2 = map[string]interface{}{
				"merchant_id":        merchantID,
				"merchant_uid":       merchantUID,
				"merchant_name":      merchantName,
				"distance":           distance,
				"distance2":          global.CalculateDistance(coordinates1.Lat, coordinates1.Lon, coordinates2.Lat, coordinates2.Lon, "km"),
				"fulfill":            fulfill,
				"total_product":      totalProduct,
				"products_available": productAvailable, // kalo product nya ga dijual ga bakal muncul disini
				"products_ordered":   productOrdered,
			}

			fmt.Println(response2)

			// mapping product_ordered
			var productItems []responseelastic.MerchantProductItems
			var totalPrice float64
			var qtyAvailable float64
			var availableItems int
			for _, item := range productOrdered {
				if item.Status != "notavailable" { // hanya menampilkan product available dan oos

					if item.Status == "available" { // hanya menambahkan price di product yang status nya available
						totalPrice = item.SpecialPrice * float64(item.QTY)
					}

					if item.QTYAvailable <= float64(item.QTY) {
						qtyAvailable = item.QTYAvailable
					} else {
						qtyAvailable = float64(item.QTY)
					}

					if item.IsAvailable {
						availableItems++
					}

					productItems = append(productItems, responseelastic.MerchantProductItems{
						SKU:          item.SKU,
						MerchantSKU:  item.MerchantSKU,
						Name:         item.Name,
						QTY:          item.QTY,
						QTYAvailable: qtyAvailable,
						UOM:          item.UOM,
						UOMName:      item.UOMName,
						SellingPrice: item.SellingPrice,
						SpecialPrice: item.SpecialPrice,
						TotalPrice:   totalPrice,
						Image:        item.Image,
						IsAvailable:  item.IsAvailable,
						// Status:       item.Status,
					})
				}
			}

			response := responseelastic.MerchantProductResponse{
				ID:             merchantID.(float64),
				UID:            merchantUID.(string),
				Name:           merchantName.(string),
				Distance:       distance,
				TotalPrice:     CalculateTotalPriceItems(productItems),
				Shippings:      shippingDuration,
				AvailableItems: availableItems,
				TotalItems:     len(productSkus),
				Items:          productItems,
			}

			allMerchantResponses = append(allMerchantResponses, response)
		}

	}

	return allMerchantResponses
}

func (s elasticMerchantServiceImpl) GetProductDetail(productSKU string, storeID int) map[string]interface{} {
	var response map[string]interface{}

	maxLimit := s.config.Elastic.MaxLimit

	params := map[string]interface{}{
		"size": maxLimit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"sku": productSKU,
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
				},
			},
		},
	}

	indexName, _ := s.getIndexFlatNameProduct(storeID)

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return response
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

	name := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["name"].(string)
	uom := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["uom"].(string)
	uomName := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["uom_name"].(string)
	images, _ := responseElastic.Hits.Hits[0].Source.(map[string]interface{})["images"].([]interface{})
	var thumbnailURL string
	if len(images) > 0 {
		imageInfo, _ := images[0].(map[string]interface{})
		thumbnailURL, _ = imageInfo["thumbnail"].(string)
	}

	response = map[string]interface{}{
		"name":     name,
		"uom":      uom,
		"uom_name": uomName,
		"image":    thumbnailURL,
	}

	return response
}

func (s elasticMerchantServiceImpl) GetPriceProduct(productSKU string, storeID int) [2]float64 {
	var response [2]float64

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
							"product_sku": productSKU,
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

	indexName, _ := s.getIndexFlatName(storeID)

	resp, err := s.elasticClient.Search(context.Background(), indexName, params)
	if err != nil {
		s.logger.Error(errors.New("error request elastic: " + err.Error()))
		return response
	}

	var responseElastic responseelastic.SearchResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseElastic)

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

	response[0] = sellingPrice
	response[1] = specialPrice

	return response
}

func CalculateTotalPriceItems(items []responseelastic.MerchantProductItems) float64 {
	total := 0.0
	for _, item := range items {
		total += item.TotalPrice
	}
	return total
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
