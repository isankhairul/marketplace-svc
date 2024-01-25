package entity

import "time"

type Category struct {
	ID              uint64    `json:"id"`
	ParentID        uint64    `json:"parent_id"`
	StoreID         uint64    `json:"store_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	MetaTitle       string    `json:"meta_title"`
	MetaDescription string    `json:"meta_description"`
	MetaKeyword     string    `json:"meta_keyword"`
	Status          int       `json:"status"`
	Position        int       `json:"position"`
	Level           int       `json:"level"`
	ChildrenCount   int       `json:"children_count"`
	SortOrder       int       `json:"sort_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Slug            string    `json:"slug"`
	FullSlug        string    `json:"full_slug"`
	Image           string    `json:"image"`
	InMenu          int       `json:"in_menu"`
	URLKey          string    `json:"url_key"`
	TermCondition   string    `json:"term_condition"`
	LandingPage     int       `json:"landing_page"`
	InHome          int       `json:"in_home"`
	Icon            string    `json:"icon"`
	IsCategoryFee   int       `json:"is_category_fee"`
	InHomepage      int       `json:"in_homepage"`
	MetaTitleH1     string    `json:"meta_title_h1"`
	IsCategoryFree  int       `json:"is_category_free"`
}
