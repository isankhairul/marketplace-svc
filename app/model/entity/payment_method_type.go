package entity

import "time"

type PaymentMethodType struct {
	ID                    uint64    `json:"id"`
	PaymentMethodTypeCode string    `json:"payment_method_type_code"`
	Name                  string    `json:"name"`
	SortOrder             int       `json:"sort_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (m PaymentMethodType) TableName() string {
	return "payment_method_type"
}
