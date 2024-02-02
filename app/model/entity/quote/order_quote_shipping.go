package entity

import (
	"marketplace-svc/app/model/entity"
	"time"
)

type OrderQuoteShipping struct {
	ID                          uint64                  `json:"id,omitempty"`
	QuoteMerchantID             uint64                  `json:"quote_merchant_id,omitempty"`
	FreeShipping                bool                    `json:"free_shipping,omitempty"`
	ShippingProviderID          uint64                  `json:"shipping_provider_id,omitempty"`
	ShippingProviderDurationID  uint64                  `json:"shipping_provider_duration_id,omitempty"`
	ShippingCostActual          float64                 `json:"shipping_cost_actual,omitempty"`
	ShippingCostSubsidized      float64                 `json:"shipping_cost_subsidized,omitempty"`
	ShippingRate                float64                 `json:"shipping_rate,omitempty"`
	ShippingQuantifier          string                  `json:"shipping_quantifier,omitempty"`
	CreatedAt                   time.Time               `json:"created_at,omitempty"`
	UpdatedAt                   time.Time               `json:"updated_at,omitempty"`
	ShippingDiscountAmount      float64                 `json:"shipping_discount_amount,omitempty"`
	TotalDiscountWeight         float64                 `json:"total_discount_weight,omitempty"`
	ShippingInsurance           int                     `json:"shipping_insurance,omitempty"`
	InstanceDelivery            int                     `json:"instance_delivery,omitempty"`
	InsuranceCost               float64                 `json:"insurance_cost,omitempty"`
	InsuranceDiscount           float64                 `json:"insurance_discount,omitempty"`
	InsuranceAmount             float64                 `json:"insurance_amount,omitempty"`
	MandatoryShippingInsurance  bool                    `json:"mandatory_shipping_insurance,omitempty"`
	SubsidizedShippingInsurance bool                    `json:"subsidized_shipping_insurance,omitempty"`
	InsuranceFeeIncluded        bool                    `json:"insurance_fee_included,omitempty"`
	DeliveryDate                *time.Time              `json:"delivery_date,omitempty"`
	ShippingWeight              float64                 `json:"shipping_weight,omitempty"`
	HasCod                      int                     `json:"has_cod,omitempty"`
	ShippingStartTime           *string                 `json:"shipping_start_time,omitempty"`
	ShippingEndTime             *string                 `json:"shipping_end_time,omitempty"`
	ShippingProvider            entity.ShippingProvider `gorm:"foreignKey:shipping_provider_id;references:id" json:"-,omitempty"`
}

func (o OrderQuoteShipping) TableName() string {
	return "order_quote_shipping"
}
