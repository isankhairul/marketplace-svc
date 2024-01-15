package entity

import "time"

type OrderQuotePayment struct {
	ID              int64     `json:"id"`
	PaymentMethodID int       `json:"payment_method_id"`
	QuoteID         int64     `json:"quote_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Notes           string    `json:"notes"`
}

func (o OrderQuotePayment) TableName() string {
	return "order_quote_payment"
}
