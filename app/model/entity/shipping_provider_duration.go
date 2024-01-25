package entity

import "time"

type ShippingProviderDuration struct {
	ID               uint64    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Duration         string    `json:"duration"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	InstanceDelivery int       `json:"instance_delivery"`
	DurationLabel    string    `json:"duration_label"`
}

func (s ShippingProviderDuration) TableName() string {
	return "shipping_provider_duration"
}
