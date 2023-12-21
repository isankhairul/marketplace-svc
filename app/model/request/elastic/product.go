package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type ProductRequest struct {
	Query   string `json:"q" schema:"q" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	StoreID *int   `json:"store_id" schema:"store_id" binding:"omitempty"`
	Page    int    `json:"page" schema:"page" binding:"omitempty"`
	Limit   int    `json:"limit" schema:"limit" binding:"omitempty"`
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
