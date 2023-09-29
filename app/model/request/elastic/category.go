package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
)

type CategoryRequest struct {
	Query      string `json:"q" schema:"q" binding:"omitempty"`
	Fields     string `json:"fields" schema:"fields" binding:"omitempty"`
	Level      *int   `json:"level" schema:"level" binding:"omitempty"`
	Position   *int   `json:"position" schema:"position" binding:"omitempty"`
	ParentID   *int   `json:"parent_id" schema:"parent_id" binding:"omitempty"`
	InHome     *int   `json:"in_home" schema:"in_home"`
	InHomepage *int   `json:"in_homepage" schema:"in_homepage"`
	InMenu     *int   `json:"in_menu" schema:"in_menu"`
	StoreID    int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Page       int    `json:"page" schema:"page" binding:"omitempty"`
	Limit      int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (b CategoryRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d-%d-%d-%d-%d-%d-%d-%d", //nolint:govet
		b.Query, b.Fields, b.Level, b.Position, b.ParentID, b.InHome, b.InHomepage, b.InMenu, b.StoreID, b.Page, b.Limit)
}

func (req CategoryRequest) DefaultPagination() CategoryRequest {
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

type CategoryTreeRequest struct {
	StoreID int `json:"store_id" schema:"store_id" binding:"omitempty"`
}

func (b CategoryTreeRequest) ToString() string {
	return fmt.Sprintf("%d", b.StoreID)
}

func (req CategoryTreeRequest) DefaultPagination() CategoryTreeRequest {
	if req.StoreID == 0 {
		req.StoreID = 1
	}

	return req
}
