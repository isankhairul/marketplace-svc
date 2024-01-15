package entity

import "time"

type PromotionCatalogProductPrice struct {
	ID                        int       `json:"id"`
	RuleDate                  string    `json:"rule_date"`
	CustomerGroupID           int       `json:"customer_group_id"`
	ProductID                 int       `json:"product_id"`
	RulePrice                 float64   `json:"rule_price"`
	StoreID                   int       `json:"store_id"`
	MerchantID                int       `json:"merchant_id"`
	PromotionCatalogID        int       `json:"promotion_catalog_id"`
	PromotionCatalogProductID int       `json:"promotion_catalog_product_id"`
	LatestStartDate           time.Time `json:"latest_start_date"`
	EarliestEndDate           time.Time `json:"earliest_end_date"`
	FromTime                  time.Time `json:"from_time"`
	ToTime                    time.Time `json:"to_time"`
	CreatedAt                 time.Time `json:"created_at"`
}
