package entity

type PromotionScprCustomerCoupon struct {
	CustomerID      uint64 `gorm:"primary_key" json:"customer_id,omitempty"`
	PromotionScprID uint64 `gorm:"primary_key" json:"promotion_scpr_id,omitempty"`
	StoreID         int    `gorm:"primary_key" json:"store_id,omitempty"`
	MaxUsage        int    `json:"max_usage,omitempty"`
	CreatedAt       int    `json:"created_at,omitempty"`
}
