package entity

import "time"

type OrderQuoteShipping struct {
	ID                          int64     `json:"id"`
	QuoteMerchantID             int       `json:"quote_merchant_id"`
	FreeShipping                bool      `json:"free_shipping"`
	ShippingProviderID          int       `json:"shipping_provider_id"`
	ShippingProviderDurationID  int       `json:"shipping_provider_duration_id"`
	ShippingCostActual          int       `json:"shipping_cost_actual"`
	ShippingCostSubsidized      int       `json:"shipping_cost_subsidized"`
	ShippingRate                int       `json:"shipping_rate"`
	ShippingQuantifier          string    `json:"shipping_quantifier"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	ShippingDiscountAmount      float64   `json:"shipping_discount_amount"`
	TotalDiscountWeight         float64   `json:"total_discount_weight"`
	ShippingInsurance           int       `json:"shipping_insurance"`
	InstanceDelivery            int       `json:"instance_delivery"`
	InsuranceCost               int       `json:"insurance_cost"`
	InsuranceDiscount           int       `json:"insurance_discount"`
	InsuranceAmount             int       `json:"insurance_amount"`
	MandatoryShippingInsurance  bool      `json:"mandatory_shipping_insurance"`
	SubsidizedShippingInsurance bool      `json:"subsidized_shipping_insurance"`
	InsuranceFeeIncluded        bool      `json:"insurance_fee_included"`
	DeliveryDate                time.Time `json:"delivery_date"`
	ShippingWeight              float64   `json:"shipping_weight"`
	HasCod                      int       `json:"has_cod"`
	ShippingStartTime           string    `json:"shipping_start_time"`
	ShippingEndTime             string    `json:"shipping_end_time"`
}

func (o OrderQuoteShipping) TableName() string {
	return "order_quote_shipping"
}
