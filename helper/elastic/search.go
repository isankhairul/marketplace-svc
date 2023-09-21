package elastic

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"golang.org/x/exp/slices"
	"marketplace-svc/app/model/base"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"math"
)

type Search interface {
	Search(ctx context.Context, collection string, query map[string]interface{}) (*esapi.Response, error)
	Pagination(dataEl map[string]interface{}, page, limit int) base.Pagination
}

func (elastic elasticClient) Search(ctx context.Context, collection string, query map[string]interface{}) (*esapi.Response, error) {
	request := esapi.SearchRequest{
		Index: []string{collection}, // Replace with your collection name
		Body:  esutil.NewJSONReader(query),
	}

	// Perform the search request
	responseEl, err := request.Do(ctx, elastic.Ec)
	if err != nil {
		return nil, err
	}
	if !slices.Contains([]int{200, 201}, responseEl.StatusCode) {
		return nil, errors.New("status code not 200, with message: " + responseEl.String())
	}

	return responseEl, nil
}

func (elastic elasticClient) Pagination(rs responseelastic.SearchResponse, page int, limit int) base.Pagination {
	totalRecords := rs.Hits.Total.Value
	var pagination base.Pagination
	pagination.Limit = limit
	pagination.Page = page
	pagination.TotalRecords = int64(totalRecords)
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(limit)))
	var records int64
	records = (int64(pagination.Limit * pagination.Page)) / int64(pagination.Page)
	if pagination.Page == pagination.TotalPage {
		records = int64(totalRecords) - ((int64(pagination.TotalPage) - 1) * int64(pagination.Limit))
	}
	if pagination.TotalPage == 0 {
		records = 0
	}
	pagination.Records = records
	return pagination
}
