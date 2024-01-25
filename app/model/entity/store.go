package entity

import "time"

type Store struct {
	ID                   uint64    `json:"id"`
	StoreCode            string    `json:"store_code"`
	StoreName            string    `json:"store_name"`
	Status               int       `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            time.Time `json:"deleted_at"`
	UpperLimitPrice      float64   `json:"upper_limit_price"`
	IsMultishop          int       `json:"is_multishop"`
	WebhookPaymentStatus string    `json:"webhook_payment_status"`
	WebhookOrderStatus   string    `json:"webhook_order_status"`
	PushExternalOrder    int       `json:"push_external_order"`
	ValidatePrice        int       `json:"validate_price"`
	IsMochaPickup        int       `json:"is_mocha_pickup"`
	AuthKey              string    `json:"auth_key"`
	HasStore             int       `json:"has_store"`
	ExcludeOosValidation int       `json:"exclude_oos_validation"`
	LockSellingPrice     int       `json:"lock_selling_price"`
	UseAdminFee          int       `json:"use_admin_fee"`
	ExcludePaymentInfo   int       `json:"exclude_payment_info"`
	ExcludeNewQuote      int       `json:"exclude_new_quote"`
}

func (s Store) TableName() string {
	return "store"
}
