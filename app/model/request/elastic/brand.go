package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type BrandRequest struct {
	Query         string `json:"q" schema:"q" binding:"omitempty"`
	Fields        string `json:"fields" schema:"fields" binding:"omitempty"`
	PrefixName    string `json:"prefix_name" schema:"prefix_name" binding:"omitempty"`
	PrincipalCode string `json:"principal_code" schema:"principal_code" binding:"omitempty"`
	ShowOfficial  *int   `json:"show_official" schema:"show_official" binding:"omitempty"`
	StoreID       int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Page          int    `json:"page" schema:"page" binding:"omitempty"`
	Limit         int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (b BrandRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%s-%s-%d-%d-%d-%d", //nolint:govet
		b.Query, b.Fields, b.PrefixName, b.PrincipalCode, b.ShowOfficial, b.StoreID, b.Page, b.Limit)
}

func (req BrandRequest) DefaultPagination() BrandRequest {
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

	// default storeID
	if req.StoreID == 0 {
		req.StoreID = 1
	}

	return req
}
