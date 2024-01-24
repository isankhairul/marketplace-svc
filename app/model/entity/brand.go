package entity

import "time"

type Brand struct {
	ID               int64     `json:"id"`
	BrandCode        string    `json:"brand_code"`
	BrandName        string    `json:"brand_name"`
	ProductGroupKnID int       `json:"product_group_kn_id"`
	Status           int       `json:"status"`
	KnsBrandID       int       `json:"kns_brand_id"`
	KlaBrandID       int       `json:"kla_brand_id"`
	PrincipalCode    string    `json:"principal_code"`
	StoreID          int       `json:"store_id"`
	Image            string    `json:"image"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Slug             string    `json:"slug"`
	SortOrder        int       `json:"sort_order"`
	MetaTitle        string    `json:"meta_title"`
	MetaDescription  string    `json:"meta_description"`
	ShowOfficial     int       `json:"show_official"`
}

func (b Brand) TableName() string {
	return "brand"
}
