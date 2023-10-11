package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type MerchantRequest struct {
	Query   string `json:"q" schema:"q" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	PID     string `json:"pid" schema:"pid"`
	CID     string `json:"cid" schema:"cid"`
	Rating  string `json:"rating" schema:"rating"`
	Type    string `json:"type" schema:"type"`
	Zipcode string `json:"zipcode" schema:"zipcode"`
	StoreID int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Store   string `json:"store" schema:"store" binding:"omitempty"`
	Sort    string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir     string `json:"dir" schema:"dir" binding:"omitempty"`
	Page    int    `json:"page" schema:"page" binding:"omitempty"`
	Limit   int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req MerchantRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d-%s-%s-%s-%d-%d-%d", //nolint:govet
		req.Query, req.Fields, req.PID, req.CID, req.Rating, req.Type, req.Zipcode, req.StoreID, req.Page, req.Limit)
}

func (req MerchantRequest) DefaultPagination() MerchantRequest {
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

type MerchantDetailRequest struct {
	Slug    string `json:"slug" schema:"slug" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	StoreID int    `json:"store_id" schema:"store_id" binding:"omitempty"`
}

func (req MerchantDetailRequest) ToString() string {
	return fmt.Sprintf("%s-%d", req.Fields, req.StoreID)
}
