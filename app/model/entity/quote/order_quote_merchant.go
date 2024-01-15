package entity

import "time"

type OrderQuoteMerchant struct {
	ID                       int              `json:"id"`
	QuoteID                  int64            `json:"quote_id"`
	MerchantID               int              `json:"merchant_id"`
	CreatedAt                time.Time        `json:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at"`
	MerchantGrandTotal       int              `json:"merchant_grand_total"`
	MerchantSubtotal         int              `json:"merchant_subtotal"`
	MerchantTotalQuantity    int              `json:"merchant_total_quantity"`
	MerchantTotalWeight      float64          `json:"merchant_total_weight"`
	MerchantTotalPointEarned int              `json:"merchant_total_point_earned"`
	MerchantNotes            string           `json:"merchant_notes"`
	MerchantTypeID           int              `json:"merchant_type_id"`
	Selected                 bool             `json:"selected"`
	Redeem                   int              `json:"redeem"`
	MerchantTotalPointSpent  int              `json:"merchant_total_point_spent"`
	AppliedRuleIds           string           `json:"applied_rule_ids"`
	PromoDescription         string           `json:"promo_description"`
	DiscountAmount           float64          `json:"discount_amount"`
	ShippingDiscountAmount   float64          `json:"shipping_discount_amount"`
	Event                    int              `json:"event"`
	AdminFeeCalculation      int              `json:"admin_fee_calculation"`
	OrderQuoteItem           []OrderQuoteItem `gorm:"foreignKey:quote_merchant_id;references:id;"`
}

func (o OrderQuoteMerchant) TableName() string {
	return "order_quote_merchant"
}
