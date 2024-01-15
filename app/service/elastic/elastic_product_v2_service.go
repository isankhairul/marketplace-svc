package elasticservice

import (
	"encoding/json"
	"fmt"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	modelelastic "marketplace-svc/app/model/elastic"
	"marketplace-svc/app/repository"
	"marketplace-svc/helper/elastic"
	"math"
)

type ElasticProductServiceV2 interface {
	//Reindex(index string, productIDs string, storeID int, productType string, merchantID int, flush bool) error
	//Search(ctx context.Context, input requestelastic.BannerRequest) ([]map[string]interface{}, base.Pagination, message.Message, error)
}

type elasticProductServiceV2Impl struct {
	logger           logger.Logger
	baseRepo         repository.BaseRepository
	productFlatRepo  repository.ProductFlatRepository
	merchantFlatRepo repository.MerchantFlatRepository
	elasticClient    elastic.ElasticClient
}

func NewElasticProductServiceV2(
	lg logger.Logger,
	br repository.BaseRepository,
	pfr repository.ProductFlatRepository,
	mf repository.MerchantFlatRepository,
	esc elastic.ElasticClient,
) ElasticProductServiceV2 {
	return &elasticProductServiceV2Impl{lg, br, pfr, mf, esc}
}

/*
func (s elasticProductServiceV2Impl) Reindex(index string, productIDs string, storeID int, productType string, merchantID int, flush bool) error {
	if index == "" {
		index = "golang_merchant_flat_" + fmt.Sprint(storeID)
	}

	titleLog := "[ES-MERCHANT-FLAT-REINDEX-STORE-" + fmt.Sprint(storeID) + "]"
	sqlCount := s.merchantFlatRepo.GetQueryRawAll(productIDs, storeID, merchantID, productType, true)

	var total int
	s.baseRepo.GetDB().Raw(sqlCount).Scan(&total)

	if total == 0 {
		s.logger.Info(titleLog + " EMPTY PRODUCTS")
		return errors.New("empty products")
	}

	chunk := 1000
	totalPage := int(math.Ceil(float64(total) / float64(chunk)))

	var sqlProduct = s.merchantFlatRepo.GetQueryRawAll(productIDs, storeID, merchantID, productType, false)
	var arrBody []interface{}
	var chanBody = make(chan interface{})

	// This will create a pool with 20 running worker goroutines
	pool := pond.New(20, totalPage, pond.IdleTimeout(30*time.Second))

	for page := 1; page <= totalPage; page++ {
		n := page
		//pool.Submit(func() {
		s.logger.Info(titleLog + "GET MERCHANT_FLAT - PAGE " + fmt.Sprint(n) + " OF " + fmt.Sprint(totalPage))
		offset := (n - 1) * chunk
		var values []map[string]interface{}
		_ = s.baseRepo.GetDB().Raw(sqlProduct + " LIMIT " + fmt.Sprint(chunk) + " OFFSET " + fmt.Sprint(offset)).Scan(&values)
		for _, value := range values {
			pid := fmt.Sprint(value["merchant_id"]) + "-" + fmt.Sprint(value["product_id"])
			//intVar, _ := strconv.ParseInt(fmt.Sprint(value["product_is_active"]), 0, 32)
			currentBody := modelelastic.EsMerchantFlat{
				ID:            pid,
				MerchantID:    value["merchant_id"].(int),
				MerchantUID:   fmt.Sprint(value["merchant_uid"]),
				ProductID:     value["product_id"].(int),
				ProductSKU:    fmt.Sprint(value["product_sku"]),
				MerchantSKU:   fmt.Sprint(value["merchant_sku"]),
				ProductStatus: value["product_status"].(int),
				Stock:         value["stock"].(int),
				StockOnHand:   value["stock_on_hand"].(int),
				Review:        value["review"].(int),
				Rating:        value["rating"].(float64),
				TypeID:        value["product_type_id"].(int),
				Status:        value["status"].(int),
				UpdatedAt:     "2023-09-19 16:54:20",
			}
			jsonCurrentBody, _ := json.Marshal(currentBody)
			fmt.Println("jsonCurrentBody", string(jsonCurrentBody))
			panic("here")

			//chanBody <- currentBody
		}
		//})
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
*/

func (s elasticProductServiceV2Impl) getMerchants(item map[string]interface{}) modelelastic.EsProductFlatMerchant {
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
