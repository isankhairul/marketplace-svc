package responsequote

import (
	"time"
)

type QuoteRs struct {
	ID                   int                `json:"id,omitempty"`
	QuoteCode            string             `json:"quote_code,omitempty"`
	StoreID              int                `json:"store_id,omitempty"`
	ContactID            string             `json:"contact_id,omitempty"`
	CustomerEmail        string             `json:"customer_email,omitempty"`
	CustomerFirstname    string             `json:"customer_firstname,omitempty"`
	CustomerLastname     string             `json:"customer_lastname,omitempty"`
	Weight               float64            `json:"weight,omitempty"`
	Subtotal             float64            `json:"subtotal,omitempty"`
	ShippingAmount       float64            `json:"shipping_amount,omitempty"`
	GrandTotal           float64            `json:"grand_total,omitempty"`
	Currency             string             `json:"currency,omitempty"`
	Status               int                `json:"status,omitempty"`
	OrderTypeID          int                `json:"order_type_id,omitempty"`
	TotalQuantity        int                `json:"total_quantity,omitempty"`
	DeviceID             int                `json:"device_id,omitempty"`
	CreatedAt            *time.Time         `json:"created_at,omitempty,omitempty"`
	UpdatedAt            *time.Time         `json:"updated_at,omitempty,omitempty"`
	ConvertedAt          *time.Time         `json:"converted_at,omitempty,omitempty"`
	CustomerID           int64              `json:"customer_id,omitempty"`
	CustomerGroupID      int                `json:"customer_group_id,omitempty"`
	SubtotalWithDiscount float64            `json:"subtotal_with_discount,omitempty"`
	BaseGrandTotal       float64            `json:"base_grand_total,omitempty"`
	OrderQuoteAddress    *QuoteAddressRs    `json:"OrderQuoteAddress,omitempty"`
	OrderQuotePayment    *QuotePaymentRs    `json:"OrderQuotePayment,omitempty"`
	OrderQuoteMerchant   *[]QuoteMerchantRs `json:"OrderQuoteMerchant,omitempty"`
}
