package entity

import "time"

type MerchantType struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Status          int       `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Description     string    `json:"description"`
	Slug            string    `json:"slug"`
	AutoIncludeItem int       `json:"auto_include_item"`
	ValidateZipcode int       `json:"validate_zipcode"`
}

func (m MerchantType) TableName() string {
	return "merchant_type"
}
