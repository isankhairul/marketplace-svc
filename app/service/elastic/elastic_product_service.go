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

// swagger:operation GET /es/products Products ProductRequest
// Product - List
//
// ---
// tags:
//   - "Elastic - Products"
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

	indexName, err := s.getIndexName()
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
	productResponse = s.transformSearch(responseElastic, arrFields)
	newProductResponse = s.transformResponse(productResponse)
	pagination = s.elasticClient.Pagination(responseElastic, input.Page, input.Limit)

	return newProductResponse, pagination, msg, nil
}

func (s elasticProductServiceImpl) getIndexName() (string, error) {
	indexName := s.config.Elastic.Index["index-products"]
	if indexName == nil {
		return "", errors.New("config index-products not defined")
	}

	return fmt.Sprint(indexName), nil
}

func (s elasticProductServiceImpl) buildQuerySearch(input requestelastic.ProductRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}

	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  input.Query,
					"fields": []string{"title"},
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
		map[string]interface{}{
			"term": map[string]interface{}{
				"is_pharmacy": 1,
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

func (s elasticProductServiceImpl) defaultFields() []string {
	return []string{
		"sku",
		"name",
		"uom",
		"uom_name",
		"weight",
		"description",
		"short_description",
		"images",
		"principal_name",
		"price",
		"min_price",
		"max_price",
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

func (s elasticProductServiceImpl) transformSearch(rs responseelastic.SearchResponse, fields []string) []map[string]interface{} {
	var response []map[string]interface{}

	for _, item := range rs.Hits.Hits {
		var tmpResponse map[string]interface{}
		jsonItem, _ := json.Marshal(item.Source)
		_ = json.Unmarshal(jsonItem, &tmpResponse)

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
