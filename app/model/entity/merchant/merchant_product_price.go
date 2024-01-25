package entity

import "time"

type MerchantProductPrice struct {
	ID                            uint64    `json:"id"`
	MerchantProductID             uint64    `json:"merchant_product_id"`
	SellingPrice                  float64   `json:"selling_price"`
	SpecialPrice                  int64     `json:"special_price"`
	BundlingQuantity              int64     `json:"bundling_quantity"`
	BundlingPrice                 int64     `json:"bundling_price"`
	BundlingSpecialPrice          int64     `json:"bundling_special_price"`
	SpecialPriceStartTime         time.Time `json:"special_price_start_time"`
	SpecialPriceEndTime           time.Time `json:"special_price_end_time"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
	StoreID                       uint64    `json:"store_id"`
	MerchantID                    uint64    `json:"merchant_id"`
	MerchantSpecialPrice          int64     `json:"merchant_special_price"`
	MerchantSpecialPriceStartTime time.Time `json:"merchant_special_price_start_time"`
	MerchantSpecialPriceEndTime   time.Time `json:"merchant_special_price_end_time"`
	MarketplaceStoreID            int       `json:"marketplace_store_id"`
}

func (m MerchantProductPrice) TableName() string {
	return "merchant_product_price"
}
