package requestelastic

import "fmt"

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
