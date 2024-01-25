package entity

import (
	"marketplace-svc/app/model/entity"
	"time"
)

type MerchantShipping struct {
	ID                 uint64                   `json:"id"`
	ShippingProviderID uint64                   `json:"shipping_provider_id"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
	Status             int                      `json:"status"`
	MerchantID         int                      `json:"merchant_id"`
	MinimumAmount      float64                  `json:"minimum_amount"`
	ShippingProvider   *entity.ShippingProvider `gorm:"foreignKey:shipping_provider_id;references:id" json:"-"`
}

func (m MerchantShipping) TableName() string {
	return "merchant_shipping"
}
