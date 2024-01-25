package entity

import "time"

type MerchantStores struct {
	ID         uint64    `json:"id"`
	StoreID    uint64    `json:"store_id"`
	MerchantID uint64    `json:"merchant_id"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (m MerchantStores) TableName() string {
	return "merchant_stores"
}
