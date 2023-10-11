package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type OrderRequest struct {
	Query  string `json:"q" schema:"q" binding:"omitempty"`
	Status string `json:"status" schema:"status" binding:"omitempty"`
	Event  bool   `json:"event" schema:"event" binding:"omitempty"`
	Dir    string `json:"dir" schema:"dir" binding:"omitempty"`
	Sort   string `json:"sort" schema:"sort" binding:"omitempty"`
	Page   int    `json:"page" schema:"page" binding:"omitempty"`
	Limit  int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req OrderRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%v-%s-%s-%d-%d", //nolint:govet
		req.Query, req.Status, req.Event, req.Dir, req.Sort, req.Page, req.Limit)
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
	return req
}
