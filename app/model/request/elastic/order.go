package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type OrderRequest struct {
	Query          string `json:"q" schema:"q" binding:"omitempty"`
	OrderNo        string `json:"order_no" schema:"order_no" binding:"omitempty"`
	CustomerID     string `json:"customer_id" schema:"customer_id"`
	Fields         string `json:"fields" schema:"fields" binding:"omitempty"`
	Status         string `json:"status" schema:"status" binding:"omitempty"`
	PaymentMethods string `json:"payment_methods" schema:"payment_methods" binding:"omitempty"`
	MerchantSlug   string `json:"merchant_slug" schema:"merchant_slug" binding:"omitempty"`
	IsReviewed     int    `json:"is_reviewed" schema:"is_reviewed" binding:"omitempty"`
	StoreID        int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Dir            string `json:"dir" schema:"dir" binding:"omitempty"`
	Sort           string `json:"sort" schema:"sort" binding:"omitempty"`
	Page           int    `json:"page" schema:"page" binding:"omitempty"`
	Limit          int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req OrderRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%v-%s-%d-%d", //nolint:govet
		req.Query, req.Status, req.Dir, req.Sort, req.Page, req.Limit)
}

func (req OrderRequest) DefaultPagination() OrderRequest {
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

	// set default storeID
	if req.StoreID == 0 {
		req.StoreID = 1
	}

	return req
}
