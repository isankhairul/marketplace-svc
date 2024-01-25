package entity

import "marketplace-svc/app/model/entity"

type OrderQuotePayment struct {
	ID              uint64                `json:"id" gorm:"column:id;"`
	PaymentMethodID uint64                `json:"payment_method_id" gorm:"column:payment_method_id;"`
	QuoteID         uint64                `json:"quote_id" gorm:"column:quote_id;"`
	Notes           string                `json:"notes" gorm:"column:notes;type:text"`
	PaymentMethod   *entity.PaymentMethod `gorm:"foreignKey:payment_method_id;references:id" json:"-"`
}

func (o OrderQuotePayment) TableName() string {
	return "order_quote_payment"
}
