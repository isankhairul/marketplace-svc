package responsequote

import (
	entity "marketplace-svc/app/model/entity/quote"
	"time"
)

type QuoteShippingRs struct { //nolint:maligned
	ShippingProviderID            int64      `json:"shipping_provider_id"`
	ShippingProviderName          string     `json:"shipping_provider_name"`
	InsuranceFeeIncluded          bool       `json:"insurance_fee_included"`
	ShippingProviderDurationID    int        `json:"shipping_provider_duration_id"`
	ShippingLogo                  string     `json:"shipping_logo"`
	ShippingProviderDurationName  string     `json:"shipping_provider_duration_name"`
	ShippingProviderDuration      string     `json:"shipping_provider_duration"`
	ShippingProviderDurationLabel string     `json:"shipping_provider_duration_label"`
	InstanceDelivery              int        `json:"instance_delivery"`
	InsuranceCost                 float64    `json:"insurance_cost"`
	InsuranceDiscount             float64    `json:"insurance_discount"`
	InsuranceAmount               float64    `json:"insurance_amount"`
	MandatoryShippingInsurance    bool       `json:"mandatory_shipping_insurance,omitempty"`
	SubsidizedShippingInsurance   bool       `json:"subsidized_shipping_insurance,omitempty"`
	ButtonInsuranceDisabled       bool       `json:"button_insurance_disabled"`
	UseInsurance                  int        `json:"use_insurance"`
	ShippingRate                  float64    `json:"shipping_rate"`
	ShippingWeight                int        `json:"shipping_weight"`
	DeliveryDate                  *time.Time `json:"delivery_date"`
	ShippingStartTime             *time.Time `json:"shipping_start_time,omitempty"`
	ShippingEndTime               *time.Time `json:"shipping_end_time,omitempty"`
	ShippingNotice                string     `json:"shipping_notice"`
	HasCod                        int        `json:"has_cod,omitempty"`
}

func (qr QuoteShippingRs) Transform(qss *entity.OrderQuoteShipping, baseImageURL string) *QuoteShippingRs {
	var response QuoteShippingRs //nolint:prealloc
	if qss == nil {
		return nil
	}

	response = QuoteShippingRs{
		ShippingProviderID:            qss.ShippingProviderID,
		ShippingProviderName:          qss.ShippingProvider.Name,
		ShippingProviderDurationID:    qss.ShippingProviderDurationID,
		InsuranceFeeIncluded:          qss.InsuranceFeeIncluded,
		ShippingLogo:                  baseImageURL + qss.ShippingProvider.Logo,
		ShippingProviderDurationName:  qss.ShippingProvider.ShippingProviderDuration.Name,
		ShippingProviderDuration:      qss.ShippingProvider.ShippingProviderDuration.Duration,
		ShippingProviderDurationLabel: qss.ShippingProvider.ShippingProviderDuration.DurationLabel,
		InstanceDelivery:              qss.InstanceDelivery,
		InsuranceCost:                 qss.InsuranceCost,
		InsuranceDiscount:             qss.InsuranceDiscount,
		InsuranceAmount:               qss.InsuranceAmount,
		DeliveryDate:                  qss.DeliveryDate,
		ShippingStartTime:             qss.ShippingStartTime,
		ShippingEndTime:               qss.ShippingEndTime,
	}

	return &response
}
