package entity

import "time"

type PromotionCatalogProductPrice struct {
	ID                        uint64    `json:"id"`
	RuleDate                  string    `json:"rule_date"`
	CustomerGroupID           int       `json:"customer_group_id"`
	ProductID                 uint64    `json:"product_id"`
	RulePrice                 float64   `json:"rule_price"`
	StoreID                   int       `json:"store_id"`
	MerchantID                uint64    `json:"merchant_id"`
	PromotionCatalogID        uint64    `json:"promotion_catalog_id"`
	PromotionCatalogProductID uint64    `json:"promotion_catalog_product_id"`
	LatestStartDate           time.Time `json:"latest_start_date"`
	EarliestEndDate           time.Time `json:"earliest_end_date"`
	FromTime                  time.Time `json:"from_time"`
	ToTime                    time.Time `json:"to_time"`
	CreatedAt                 time.Time `json:"created_at"`
}

func (m PromotionCatalogProductPrice) TableName() string {
	return "promotion_catalog_product_price"
}
