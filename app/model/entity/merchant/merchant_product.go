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

type DetailMerchantProduct struct {
	ID                     uint64     `json:"id"`
	Sku                    string     `json:"sku"`
	MerchantSku            string     `json:"merchant_sku"`
	MerchantIncludedItem   bool       `json:"merchant_included_item"`
	Name                   string     `json:"name"`
	Slug                   string     `json:"slug"`
	Weight                 float64    `json:"weight"`
	BrandCode              string     `json:"brand_code"`
	BasePoint              int        `json:"base_point"`
	BasePrice              float64    `json:"base_price"`
	ProductType            string     `json:"product_type"`
	AttributeSetID         int        `json:"attribute_set_id"`
	RewardPointSellProduct int        `json:"reward_point_sell_product"`
	SpecialPriceStartTime  *time.Time `json:"special_price_start_time"`
	SpecialPriceEndTime    *time.Time `json:"special_price_end_time"`
	SellingPrice           float64    `json:"selling_price"`
	SpecialPrice           float64    `json:"special_price"`
}
