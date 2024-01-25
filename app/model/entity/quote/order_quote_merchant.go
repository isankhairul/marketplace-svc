package entity

import (
	entity "marketplace-svc/app/model/entity/merchant"
	"time"
)

type OrderQuoteMerchant struct {
	ID                       uint64              `json:"id"`
	QuoteID                  uint64              `json:"quote_id"`
	MerchantID               uint64              `json:"merchant_id"`
	CreatedAt                time.Time           `json:"created_at"`
	UpdatedAt                time.Time           `json:"updated_at"`
	MerchantGrandTotal       float64             `json:"merchant_grand_total"`
	MerchantSubtotal         float64             `json:"merchant_subtotal"`
	MerchantTotalQuantity    int                 `json:"merchant_total_quantity"`
	MerchantTotalWeight      float64             `json:"merchant_total_weight"`
	MerchantTotalPointEarned float64             `json:"merchant_total_point_earned"`
	MerchantNotes            string              `json:"merchant_notes"`
	MerchantTypeID           uint64              `json:"merchant_type_id"`
	Selected                 bool                `json:"selected"`
	Redeem                   int                 `json:"redeem"`
	MerchantTotalPointSpent  int                 `json:"merchant_total_point_spent"`
	AppliedRuleIds           string              `json:"applied_rule_ids"`
	PromoDescription         string              `json:"promo_description"`
	DiscountAmount           float64             `json:"discount_amount"`
	ShippingDiscountAmount   float64             `json:"shipping_discount_amount"`
	Event                    int                 `json:"event"`
	AdminFeeCalculation      float64             `json:"admin_fee_calculation"`
	Merchant                 entity.Merchant     `gorm:"foreignKey:merchant_id;references:id;"`
	OrderQuoteItem           *[]OrderQuoteItem   `gorm:"foreignKey:quote_merchant_id;references:id;"`
	OrderQuoteShipping       *OrderQuoteShipping `gorm:"foreignKey:quote_merchant_id;references:id;"`
}

func (o OrderQuoteMerchant) TableName() string {
	return "order_quote_merchant"
}
