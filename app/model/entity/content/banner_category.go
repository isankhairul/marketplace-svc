package entity

import "time"

type BannerCategory struct {
	ID                      int       `json:"id"`
	CategoryName            string    `json:"category_name"`
	CategorySlug            string    `json:"category_slug"`
	Status                  int       `json:"status"`
	MappingBannerCategoryID string    `json:"mapping_banner_category_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
