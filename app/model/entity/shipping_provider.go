package entity

import "time"

type ShippingProvider struct {
	ID                           int64                     `json:"id"`
	Code                         string                    `json:"code"`
	ShippingProviderTypeID       int                       `json:"shipping_provider_type_id"`
	Name                         string                    `json:"name"`
	CreatedAt                    time.Time                 `json:"created_at"`
	UpdatedAt                    time.Time                 `json:"updated_at"`
	ShippingProviderDurationID   int64                     `json:"shipping_provider_duration_id"`
	Status                       int                       `json:"status"`
	FinanceCoa                   string                    `json:"finance_coa"`
	StartTime                    string                    `json:"start_time"`
	EndTime                      string                    `json:"end_time"`
	InsuranceFee                 int                       `json:"insurance_fee"`
	InsuranceFeeIncluded         bool                      `json:"insurance_fee_included"`
	MdrTypeID                    int                       `json:"mdr_type_id"`
	Logo                         string                    `json:"logo"`
	RepickupTime                 int                       `json:"repickup_time"`
	ShippingMethodIntegration    int                       `json:"shipping_method_integration"`
	MaxWeight                    int                       `json:"max_weight"`
	MaxDistance                  int                       `json:"max_distance"`
	MaxBaggage                   int                       `json:"max_baggage"`
	MinimumAmount                float64                   `json:"minimum_amount"`
	HasCod                       int                       `json:"has_cod"`
	RateTypeID                   int                       `json:"rate_type_id"`
	Priority                     int                       `json:"priority"`
	Price                        int                       `json:"price"`
	InternalShippingCoverageCode string                    `json:"internal_shipping_coverage_code"`
	ShippingProviderDuration     *ShippingProviderDuration `gorm:"foreignKey:shipping_provider_duration_id;references:id" json:"-"`
}

func (s ShippingProvider) TableName() string {
	return "shipping_provider"
}
