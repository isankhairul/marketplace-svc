package elastic

import (
	"context"
	"marketplace-svc/app/model/base"
	"math"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
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

	return responseEl, nil
}

func (elastic elasticClient) Pagination(dataEl map[string]interface{}, page, limit int) base.Pagination {
	hits := dataEl["hits"].(map[string]interface{})
	hitsTotal := hits["total"].(map[string]interface{})
	totalRecords := hitsTotal["value"].(float64)
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
