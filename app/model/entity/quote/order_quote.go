package entity

import (
	"encoding/json"
	"time"
)

type OrderQuote struct {
	ID                        uint64                `json:"id,omitempty"`
	QuoteCode                 string                `json:"quote_code,omitempty"`
	StoreID                   uint64                `json:"store_id,omitempty"`
	ContactID                 string                `json:"contact_id,omitempty"`
	CustomerEmail             string                `json:"customer_email,omitempty"`
	CustomerFirstname         string                `json:"customer_firstname,omitempty"`
	CustomerLastname          string                `json:"customer_lastname,omitempty"`
	Weight                    float64               `json:"weight,omitempty"`
	Subtotal                  float64               `json:"subtotal,omitempty"`
	DiscountAmount            float64               `json:"discount_amount,omitempty"`
	ShippingAmount            float64               `json:"shipping_amount,omitempty"`
	GrandTotal                float64               `json:"grand_total,omitempty"`
	CustomerNotes             string                `json:"customer_notes,omitempty"`
	DiscountDescription       string                `json:"discount_description,omitempty"`
	ShippingDiscountAmount    float64               `json:"shipping_discount_amount,omitempty"`
	TotalPointEarned          int                   `json:"total_point_earned,omitempty"`
	TotalPointSpent           int                   `json:"total_point_spent,omitempty"`
	TotalPointSpentConversion float64               `json:"total_point_spent_conversion,omitempty"`
	Currency                  string                `json:"currency,omitempty"`
	Status                    int                   `json:"status,omitempty"`
	CustomerTypeID            int                   `json:"customer_type_id,omitempty"`
	OrderTypeID               uint8                 `json:"order_type_id,omitempty"`
	TotalQuantity             int                   `json:"total_quantity,omitempty"`
	AppliedRuleIds            string                `json:"applied_rule_ids,omitempty"`
	CouponCode                string                `json:"coupon_code,omitempty"`
	DeviceID                  uint8                 `json:"device_id,omitempty"`
	CouponDiscountAmount      float64               `json:"coupon_discount_amount,omitempty"`
	SubsidizedAmount          float64               `json:"subsidized_amount,omitempty"`
	CreatedAt                 *time.Time            `json:"created_at,omitempty"`
	UpdatedAt                 *time.Time            `json:"updated_at,omitempty"`
	ConvertedAt               *time.Time            `json:"converted_at,omitempty"`
	PaymentMethodID           int                   `json:"payment_method_id,omitempty"`
	CustomerID                uint64                `json:"customer_id,omitempty"`
	CustomerGroupID           int                   `json:"customer_group_id,omitempty"`
	BaseSubtotal              float64               `json:"base_subtotal,omitempty"`
	SubtotalWithDiscount      float64               `json:"subtotal_with_discount,omitempty"`
	BaseSubtotalDiscount      float64               `json:"base_subtotal_discount,omitempty"`
	BaseDiscountAmount        float64               `json:"base_discount_amount,omitempty"`
	Redeem                    int                   `json:"redeem,omitempty"`
	BaseGrandTotal            float64               `json:"base_grand_total,omitempty"`
	PromoDescription          string                `json:"promo_description,omitempty"`
	TotalPointBonus           float64               `json:"total_point_bonus,omitempty"`
	TotalPointDiscount        int                   `json:"total_point_discount,omitempty"`
	Scope                     string                `json:"scope,omitempty"`
	AppID                     string                `json:"app_id,omitempty"`
	InsuranceAmount           int32                 `json:"insurance_amount,omitempty"`
	CustomerData              string                `json:"customer_data,omitempty"`
	AgentID                   string                `json:"agent_id,omitempty"`
	DataSource                int                   `json:"data_source,omitempty"`
	AdminFee                  int                   `json:"admin_fee,omitempty"`
	AdminFeeCalculation       int                   `json:"admin_fee_calculation,omitempty"`
	AdminFeeType              string                `json:"admin_fee_type,omitempty"`
	AdminFeeTypeID            int                   `json:"admin_fee_type_id,omitempty"`
	DataSourceValue           string                `json:"data_source_value,omitempty"`
	DataReceipt               json.RawMessage       `json:"data_receipt,omitempty"`
	OrderQuoteAddress         *OrderQuoteAddress    `gorm:"foreignKey:quote_id;references:id"`
	OrderQuotePayment         *OrderQuotePayment    `gorm:"foreignKey:quote_id;references:id"`
	OrderQuoteMerchant        *[]OrderQuoteMerchant `gorm:"foreignKey:quote_id;references:id;"`
}

func (o OrderQuote) TableName() string {
	return "order_quote"
}
