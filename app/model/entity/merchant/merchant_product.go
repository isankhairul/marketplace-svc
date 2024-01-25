package entity

import "time"

type MerchantProduct struct {
	ID                   uint64                 `json:"id"`
	ProductSku           string                 `json:"product_sku"`
	MerchantSku          string                 `json:"merchant_sku"`
	Status               int                    `json:"status"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	CategoryID           int                    `json:"category_id"`
	Stock                int                    `json:"stock"`
	ReservedStock        int                    `json:"reserved_stock"`
	StockOnHand          int                    `json:"stock_on_hand"`
	BufferStock          int                    `json:"buffer_stock"`
	MerchantID           uint64                 `json:"merchant_id"`
	ProductID            uint64                 `json:"product_id"`
	MaxPurchaseQty       int                    `json:"max_purchase_qty"`
	ParentReservedStock  int                    `json:"parent_reserved_stock"`
	MerchantIncludedItem bool                   `json:"merchant_included_item"`
	OldStatus            int                    `json:"old_status"`
	UpdatedBy            string                 `json:"updated_by"`
	MerchantProductPrice []MerchantProductPrice `gorm:"foreignKey:merchant_product_id;references:id;"`
}

func (m MerchantProduct) TableName() string {
	return "merchant_product"
}
