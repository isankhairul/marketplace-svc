package entity

import (
	entity "marketplace-svc/app/model/entity/catalog"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	"time"
)

type OrderQuoteItem struct {
	ID                   uint64                          `json:"id"`
	QuoteMerchantID      uint64                          `json:"quote_merchant_id"`
	ProductID            uint64                          `json:"product_id"`
	ItemTypeID           uint64                          `json:"item_type_id"`
	ProductSku           string                          `json:"product_sku"`
	MerchantSku          string                          `json:"merchant_sku"`
	MerchantCategoryID   int                             `json:"merchant_category_id"`
	CategoryID           *uint64                         `json:"category_id"`
	BrandID              *uint64                         `json:"brand_id"`
	Name                 string                          `json:"name"`
	ItemNotes            string                          `json:"item_notes"`
	Weight               float64                         `json:"weight"`
	Quantity             int32                           `json:"quantity"`
	Price                float64                         `json:"price"`
	DiscountPercentage   float64                         `json:"discount_percentage"`
	DiscountAmount       float64                         `json:"discount_amount"`
	RowWeight            float64                         `json:"row_weight"`
	RowTotal             float64                         `json:"row_total"`
	OriginalPrice        float64                         `json:"original_price"`
	RowOriginalPrice     float64                         `json:"row_original_price"`
	PointEarned          int                             `json:"point_earned"`
	PointSpent           int                             `json:"point_spent"`
	PointSpentConversion float64                         `json:"point_spent_conversion"`
	CreatedAt            time.Time                       `json:"created_at"`
	UpdatedAt            time.Time                       `json:"updated_at"`
	RowPointEarned       int                             `json:"row_point_earned"`
	BasePrice            float64                         `json:"base_price"`
	BaseDiscountAmount   float64                         `json:"base_discount_amount"`
	BaseRowTotal         float64                         `json:"base_row_total"`
	PromoDescription     string                          `json:"promo_description"`
	BonusPoint           float64                         `json:"bonus_point"`
	DiscountPoint        int                             `json:"discount_point"`
	DiscountWeight       float64                         `json:"discount_weight"`
	RowPointSpent        int                             `json:"row_point_spent"`
	Redeem               int                             `json:"redeem"`
	MinPrice             float64                         `json:"min_price"`
	MaxPrice             float64                         `json:"max_price"`
	ProductKn            int                             `json:"product_kn"`
	ProductKalbe         int                             `json:"product_kalbe"`
	BrandName            string                          `json:"brand_name"`
	BrandCode            string                          `json:"brand_code"`
	Event                int                             `json:"event"`
	Selected             bool                            `json:"selected"`
	Location             string                          `json:"location"`
	AdditionalInfo       string                          `json:"additional_info"`
	ParentInfo           string                          `json:"parent_info"`
	StartDate            time.Time                       `json:"start_date"`
	EndDate              time.Time                       `json:"end_date"`
	EventOnline          int                             `json:"event_online"`
	ProductType          string                          `json:"product_type"`
	QuoteID              uint64                          `json:"quote_id"`
	AttributeSetID       int                             `json:"attribute_set_id"`
	FreeProduct          bool                            `json:"free_product"`
	FreeProductCommit    bool                            `json:"free_product_commit"`
	FreeProductRuleID    int                             `json:"free_product_rule_id"`
	MerchantSpecialPrice float64                         `json:"merchant_special_price"`
	NonChangeableItem    bool                            `json:"non_changeable_item"`
	MerchantIncludedItem bool                            `json:"merchant_included_item"`
	IsKliknow            int                             `json:"is_kliknow"`
	ProductFlat          *entity.ProductFlat             `gorm:"foreignKey:product_sku;references:sku;" json:"-"`
	MerchantProduct      *entitymerchant.MerchantProduct `gorm:"foreignKey:product_sku;references:product_sku;" json:"-"`
	Product              *entity.Product                 `gorm:"foreignKey:product_id;references:id;" json:"-"`
	ProductCategory      *[]entity.ProductCategory       `gorm:"foreignKey:product_id;references:product_id;" json:"-"`
	Merchant             *entitymerchant.Merchant        `gorm:"-"`
}

func (o OrderQuoteItem) TableName() string {
	return "order_quote_item"
}

type OrderQuoteItemJoin struct {
	ID                   uint64     `json:"id"`
	QuoteMerchantID      uint64     `json:"quote_merchant_id"`
	ProductID            uint64     `json:"product_id"`
	ProductSku           string     `json:"product_sku"`
	MerchantSku          string     `json:"merchant_sku"`
	MerchantCategoryID   uint64     `json:"merchant_category_id"`
	CategoryID           *uint64    `json:"category_id"`
	BrandID              *uint64    `json:"brand_id"`
	Name                 string     `json:"name"`
	ItemNotes            string     `json:"item_notes"`
	Weight               float64    `json:"weight"`
	Quantity             int32      `json:"quantity"`
	Price                float64    `json:"price"`
	DiscountPercentage   float64    `json:"discount_percentage"`
	DiscountAmount       float64    `json:"discount_amount"`
	RowWeight            float64    `json:"row_weight"`
	RowTotal             float64    `json:"row_total"`
	OriginalPrice        float64    `json:"original_price"`
	RowOriginalPrice     float64    `json:"row_original_price"`
	BasePrice            float64    `json:"base_price"`
	BaseDiscountAmount   float64    `json:"base_discount_amount"`
	BaseRowTotal         float64    `json:"base_row_total"`
	PromoDescription     string     `json:"promo_description"`
	DiscountWeight       float64    `json:"discount_weight"`
	MinPrice             float64    `json:"min_price"`
	MaxPrice             float64    `json:"max_price"`
	BrandName            string     `json:"brand_name"`
	BrandCode            string     `json:"brand_code"`
	Selected             bool       `json:"selected"`
	Location             string     `json:"location"`
	AdditionalInfo       string     `json:"additional_info"`
	ParentInfo           string     `json:"parent_info"`
	StartDate            *time.Time `json:"start_date"`
	EndDate              *time.Time `json:"end_date"`
	ProductType          string     `json:"product_type"`
	QuoteID              uint64     `json:"quote_id"`
	AttributeSetID       int        `json:"attribute_set_id"`
	MerchantSpecialPrice float64    `json:"merchant_special_price"`
	PfIsActive           int        `json:"pf_is_active"`
	PfStatus             int        `json:"pf_status"`
	PID                  uint64     `json:"p_id"`
	PSku                 string     `json:"p_sku"`
	PName                string     `json:"p_name"`
	PSlug                string     `json:"p_slug"`
	OqmMerchantID        uint64     `json:"oqm_merchant_id"`
	MID                  uint64     `json:"m_id"`
	MMerchantUID         string     `json:"m_merchant_uid"`
	MMerchantCode        string     `json:"m_merchant_code"`
	MMerchantName        string     `json:"m_merchant_name"`
	MSlug                string     `json:"m_slug"`
}
