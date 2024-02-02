package response_kalcare

type ListShippingRateDuration struct {
	Meta meta `json:"meta"`
	Data struct {
		Records []ShippingRateDuration `json:"records"`
	} `json:"data"`
	Errors struct {
	} `json:"errors"`
}

type ListShippingRateProvider struct {
	Meta meta `json:"meta"`
	Data struct {
		Records []ShippingRateDuration `json:"records"`
	} `json:"data"`
	Errors struct {
	} `json:"errors"`
}

type ShippingRateDuration struct {
	ShippingProviderID           uint64  `json:"shipping_provider_id"`
	Code                         string  `json:"code"`
	PriceKg                      float64 `json:"price_kg"`
	ShippingCost                 float64 `json:"shipping_cost"`
	Message                      []any   `json:"message"`
	ShippingProviderName         string  `json:"shipping_provider_name"`
	ShippingProviderDurationID   uint64  `json:"shipping_provider_duration_id"`
	ShippingProviderDurationName string  `json:"shipping_provider_duration_name"`
	ShippingProviderDuration     string  `json:"shipping_provider_duration"`
	InstanceDelivery             int     `json:"instance_delivery"`
	ShippingInsurance            int     `json:"shipping_insurance"`
	InsuranceFeeIncluded         bool    `json:"insurance_fee_included"`
	ShippingLogo                 string  `json:"shipping_logo"`
	ShippingStartTime            string  `json:"shipping_start_time"`
	ShippingEndTime              string  `json:"shipping_end_time"`
	HasCod                       int     `json:"has_cod"`
}
