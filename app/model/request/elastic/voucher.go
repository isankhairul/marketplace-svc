package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type VoucherRequest struct {
	Fields   string `json:"fields" schema:"fields" binding:"omitempty"`
	Criteria string `json:"criteria" schema:"criteria" binding:"omitempty"`
	Value    string `json:"value" schema:"value" binding:"omitempty"`
	StoreID  int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Page     int    `json:"page" schema:"page" binding:"omitempty"`
	Limit    int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req VoucherRequest) DefaultPagination() VoucherRequest {
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

func (req VoucherRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%s-%d-%d-%d", req.Fields, req.Criteria, req.Value, req.StoreID, req.Page, req.Limit)
}
