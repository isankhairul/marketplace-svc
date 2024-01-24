package entity

import "marketplace-svc/app/model/entity"

type OrderQuotePayment struct {
	ID              int64                 `json:"id" gorm:"column:id;"`
	PaymentMethodID int64                 `json:"payment_method_id" gorm:"column:payment_method_id;"`
	QuoteID         int64                 `json:"quote_id" gorm:"column:quote_id;"`
	Notes           string                `json:"notes" gorm:"column:notes;type:text"`
	PaymentMethod   *entity.PaymentMethod `gorm:"foreignKey:payment_method_id;references:id" json:"-"`
}

func (o OrderQuotePayment) TableName() string {
	return "order_quote_payment"
}
