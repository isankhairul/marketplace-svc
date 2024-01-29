package responsequote

import (
	"encoding/json"
)

type QuoteRs struct {
	ID                 uint64             `json:"id,omitempty"`
	QuoteCode          string             `json:"quote_code"`
	StoreID            uint64             `json:"store_id"`
	CustomerEmail      string             `json:"customer_email"`
	CustomerFirstname  string             `json:"customer_firstname"`
	CustomerLastname   string             `json:"customer_lastname"`
	Weight             float64            `json:"weight"`
	Subtotal           float64            `json:"subtotal"`
	ShippingAmount     float64            `json:"shipping_amount"`
	DiscountAmount     float64            `json:"discount_amount"`
	GrandTotal         float64            `json:"grand_total"`
	Currency           string             `json:"currency,omitempty"`
	TotalQuantity      int                `json:"total_quantity"`
	DeviceID           uint8              `json:"device_id"`
	CustomerID         uint64             `json:"customer_id,omitempty"`
	CouponCode         string             `json:"coupon_code"`
	PromoDescription   *string            `json:"promo_description"`
	OrderTypeID        uint8              `json:"order_type_id"`
	OrderQuoteAddress  *[]QuoteAddressRs  `json:"OrderQuoteAddress"`
	OrderQuotePayment  *[]QuotePaymentRs  `json:"OrderQuotePayment"`
	OrderQuoteMerchant *[]QuoteMerchantRs `json:"OrderQuoteMerchant"`
	OrderQuoteReceipt  *json.RawMessage   `json:"OrderQuoteReceipt"`
}
