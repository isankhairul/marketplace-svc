package requestelastic

import (
	"encoding/json"
	"fmt"
	"marketplace-svc/app/model/base"
)

// swagger:parameters ProductRequest
type ProductRequest struct {
	//Global Search
	// in: query
	Query string `json:"q" schema:"q" binding:"omitempty"`
	// Additional Fields
	// Example: "description,short_description"
	Fields string `json:"fields" schema:"fields" binding:"omitempty"`
	// StoreID
	StoreID *int `json:"store_id" schema:"store_id" binding:"omitempty"`
	// Page number
	Page int `json:"page" schema:"page" binding:"omitempty"`
	// Maximum records per page
	Limit int `json:"limit" schema:"limit" binding:"omitempty"`
	// JSON value for filtering. The format is {"field_name": "filter_keyword", ...}
	// Example {"name":"Dr Ican","str_no":"XXXX-XXXX"]}
	Filter string `schema:"filter"`
	// swagger:ignore
	ParsedFilter ProductFilter `schema:"-"`
}

func (b ProductRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d-%d", //nolint:govet
		b.Query, b.Fields, b.StoreID, b.Page, b.Limit)
}

func (req ProductRequest) DefaultPagination() ProductRequest {
	if req.Limit == 0 {
		req.Limit = base.PAGINATION_MIN_LIMIT
	}
	if req.Limit > base.PAGINATION_MAX_LIMIT {
		req.Limit = base.PAGINATION_MAX_LIMIT
	}

	// Default page 1
	if req.Page == 0 {
		req.Page = 1
	}
	return req
}

type ProductFilter struct {
	ProdCode []string `json:"prod_code"`
}

func (req *ProductRequest) ParseFilterAndSetDefault() error {

	var filter ProductFilter
	if req.Filter != "" {
		err := json.Unmarshal([]byte(req.Filter), &filter)
		if err != nil {
			return err
		}
	}

	req.ParsedFilter = filter
	return nil
}
