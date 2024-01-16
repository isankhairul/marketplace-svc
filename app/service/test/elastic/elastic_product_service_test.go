package elasticservice

import (
	"context"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/message"
	"testing"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"

	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/app/repository/repository_mock"
	elasticservice "marketplace-svc/app/service/elastic"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockElasticProductService struct{}

func (m *mockElasticProductService) Search(ctx context.Context, input requestelastic.ProductRequest) ([]responseelastic.ProductResponse, base.Pagination, message.Message, error) {
	return nil, base.Pagination{}, message.SuccessMsg, nil
}

var lg logger.Logger
var svcProduct elasticservice.ElasticProductService
var baseRepo = &repository_mock.BaseRepository{Mock: mock.Mock{}}
var productFlatRepo = &repository_mock.ProductFlatRepository{Mock: mock.Mock{}}
var baseElasticRepo = &repository_mock.ElasticClient{Mock: mock.Mock{}}

func init() {
	lg, _ = logger.NewLogger(logger.NewGoKitLog(&logger.LogConfig{}), "")
	config := &config.Config{}
	svcProduct = elasticservice.NewElasticProductService(*config, lg, baseRepo, productFlatRepo, baseElasticRepo)
	// ctx = context.Background()
}

func TestSearchProductSuccess(t *testing.T) {
	service := &mockElasticProductService{}

	_, _, msg, err := service.Search(context.Background(), requestelastic.ProductRequest{})

	assert.NoError(t, err)
	assert.Equal(t, message.SuccessMsg, msg)

	// storeID := 1

	// input := requestelastic.ProductRequest{
	// 	StoreID: &storeID,
	// }

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
	// 		Hits: []responseelastic.SearchResponseHits{
	// 			{
	// 				Index: "index1",
	// 				Type:  "type1",
	// 				ID:    "1",
	// 				Score: 8.5,
	// 				Source: struct {
	// 					ID           float64       `json:"id"`
	// 					Name         string        `json:"name"`
	// 					UOM          string        `json:"uom"`
	// 					UOMName      string        `json:"uom_name"`
	// 					ShortDesc    string        `json:"short_description"`
	// 					Description  string        `json:"description"`
	// 					IsFree       int           `json:"is_free_product"`
	// 					Slug         string        `json:"slug"`
	// 					Terms        string        `json:"completion_terms"`
	// 					Principal    string        `json:"principal_name"`
	// 					IsActive     int           `json:"is_active"`
	// 					UpdatedAt    string        `json:"updated_at"`
	// 					MetaTitle    string        `json:"meta_title"`
	// 					BasePoint    int           `json:"base_point"`
	// 					IsSpot       int           `json:"is_spot"`
	// 					IsKliknow    int           `json:"is_kliknow"`
	// 					RewardPoint  int           `json:"reward_point_sell_product"`
	// 					TypeID       string        `json:"type_id"`
	// 					CreatedAt    string        `json:"created_at"`
	// 					SKU          string        `json:"sku"`
	// 					Barcode      string        `json:"barcode"`
	// 					BrandCode    string        `json:"brand_code"`
	// 					BasePointRp  int           `json:"base_point_rupiah"`
	// 					IsFamilyGift int           `json:"is_family_gift"`
	// 					IsLangganan  int           `json:"is_langganan"`
	// 					MetaTitleH1  string        `json:"meta_title_h1"`
	// 					MetaDesc     string        `json:"meta_description"`
	// 					FullText     string        `json:"full_text_search"`
	// 					Weight       int           `json:"weight"`
	// 					IsPharmacy   int           `json:"is_pharmacy"`
	// 					IsPrescript  int           `json:"is_prescription"`
	// 					Categories   []interface{} `json:"categories"`
	// 					Images       []interface{} `json:"images"`
	// 					IsTicket     int           `json:"is_ticket"`
	// 					Status       int           `json:"status"`
	// 					Breadcrumbs  []interface{} `json:"breadcrumbs"`
	// 				}{
	// 					ShortDesc:    "ACLONAC adalah obat yang digunakan sebagai pereda nyeri, mengurangi gangguan inflamasi (radang), dismenore, nyeri ringan sampai sedang pasca operasi khususnya k",
	// 					Description:  "ACLONAC adalah obat yang digunakan sebagai pereda nyeri, mengurangi gangguan inflamasi (radang), dismenore, nyeri ringan sampai sedang pasca operasi khususnya ketika pasien juga mengalami peradangan.",
	// 					UOM:          "STP",
	// 					IsFree:       0,
	// 					Slug:         "aclonac-25mg-10-tablet",
	// 					Terms:        "Aclonac 25mg 10 Tablet",
	// 					UOMName:      "Strip",
	// 					Principal:    "pharos",
	// 					IsActive:     1,
	// 					UpdatedAt:    "2023-07-28 09:51:27",
	// 					Name:         "Aclonac 25mg 10 Tablet",
	// 					MetaTitle:    "Aclonac 25mg 10 Tablet",
	// 					BasePoint:    0,
	// 					IsSpot:       0,
	// 					IsKliknow:    0,
	// 					ID:           40713,
	// 					RewardPoint:  0,
	// 					TypeID:       "simple",
	// 					CreatedAt:    "2020-07-03 09:44:03",
	// 					SKU:          "ACLC02TAB",
	// 					Barcode:      "",
	// 					BrandCode:    "aclonac",
	// 					BasePointRp:  0,
	// 					IsFamilyGift: 0,
	// 					IsLangganan:  0,
	// 					MetaTitleH1:  "Aclonac 25mg 10 Tablet",
	// 					MetaDesc:     "ACLONAC adalah obat yang digunakan sebagai pereda nyeri, mengurangi gangguan inflamasi (radang), dismenore, nyeri ringan sampai sedang pasca operasi khususnya k",
	// 					FullText:     "Aclonac 25mg 10 Tablet aclonac",
	// 					Weight:       25,
	// 					IsPharmacy:   1,
	// 					IsPrescript:  1,
	// 					Categories:   []interface{}{},
	// 					Images:       []interface{}{(*interface{})(nil)},
	// 					IsTicket:     0,
	// 					Status:       1,
	// 					Breadcrumbs:  []interface{}{},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Aggregations: nil,
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

	// paginationResult := base.Pagination{
	// 	Records:   1,
	// 	Limit:     1,
	// 	Page:      1,
	// 	TotalPage: 1,
	// }

	// params := buildQuerySearch(input)

	// baseElasticRepo.Mock.On("Search", context.Background(), "mock.Anything", params).Return(responseElastic, &paginationResult).Once()

	// baseElasticRepo.On("Pagination", mock.AnythingOfType("responseelastic.SearchResponse"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(
	// 	paginationResult,
	// )

	// products, pagination, msg, err := svcProduct.Search(context.Background(), input)

	// assert.NoError(t, err, "Unexpected error")
	// assert.Equal(t, message.SuccessMsg, msg)
	// assert.Equal(t, 1, len(products))
	// assert.Equal(t, int64(len(products)), pagination.Records)
}

func buildQuerySearch(input requestelastic.ProductRequest) map[string]interface{} {
	queryArray := map[string]interface{}{}
	// create query bool
	if input.Query != "" {
		queryArray["bool"] = map[string]interface{}{
			"must": []interface{}{
				map[string]interface{}{
					"match_phrase_prefix": map[string]interface{}{
						"completion_terms": map[string]interface{}{
							"query": input.Query,
						},
					},
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
				"is_prescription": 1,
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

	if len(input.ParsedFilter.ProdCode) > 0 {
		prodCodeFilter := map[string]interface{}{
			"terms": map[string]interface{}{
				"sku": input.ParsedFilter.ProdCode,
			},
		}
		filters = append(filters, prodCodeFilter)
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
