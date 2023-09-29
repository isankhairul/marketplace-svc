package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type MerchantRequest struct {
	Query   string `json:"q" schema:"q" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	PID     *int   `json:"pid" schema:"pid"`
	CID     *int   `json:"cid" schema:"cid"`
	Rating  string `json:"rating" schema:"rating"`
	Type    string `json:"type" schema:"type"`
	Zipcode string `json:"zipcode" schema:"zipcode"`
	StoreID int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Page    int    `json:"page" schema:"page" binding:"omitempty"`
	Limit   int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (b MerchantRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d-%s-%s-%s-%d-%d-%d", //nolint:govet
		b.Query, b.Fields, b.PID, b.CID, b.Rating, b.Type, b.Zipcode, b.StoreID, b.Page, b.Limit)
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
