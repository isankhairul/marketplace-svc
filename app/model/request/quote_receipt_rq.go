package request

import (
	validation "github.com/itgelo/ozzo-validation/v4"
	"time"
)

type CheckQuoteRq struct {
	DeviceID int
}

type QuoteReceiptRq struct {
	QuoteCode           string                    `json:"-"`
	OrderQuotePayment   *QuoteReceiptPaymentRq    `json:"orderQuotePayment"`
	OrderQuoteAddress   *QuoteReceiptAddressRq    `json:"orderQuoteAddress"`
	OrderQuoteReceipt   *interface{}              `json:"orderQuoteReceipt"`
	OrderQuoteMerchants []QuoteReceiptMerchantsRq `json:"orderQuoteMerchants"`
}

type QuoteReceiptPaymentRq struct {
	PaymentMethodID uint64 `json:"payment_method_id"`
}

func (req QuoteReceiptPaymentRq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.PaymentMethodID, validation.Required.Error("payment_method_id is required")))
}

type QuoteReceiptAddressRq struct {
	CustomerAddressID uint64 `json:"customer_address_id"`
}

func (req QuoteReceiptAddressRq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CustomerAddressID, validation.Required.Error("customer_address_id is required")),
	)
}

type QuoteReceiptMerchantsRq struct {
	MerchantID         uint64                 `json:"merchant_id"`
	OrderQuoteItems    []QuoteReceiptItemsRq  `json:"orderQuoteItems"`
	OrderQuoteShipping QuoteReceiptShippingRq `json:"orderQuoteShipping"`
}

func (req QuoteReceiptMerchantsRq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.MerchantID, validation.Required.Error("merchant_id is required")),
	)
}

type QuoteReceiptItemsRq struct {
	QuoteMerchantID uint64 `json:"-"`
	MerchantID      uint64 `json:"-"`
	MerchantSku     string `json:"merchant_sku"`
	ItemNotes       string `json:"item_notes"`
	Quantity        int32  `json:"quantity"`
}

func (req QuoteReceiptItemsRq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.MerchantSku, validation.Required.Error("merchant_sku is required")),
		validation.Field(&req.Quantity, validation.Required.Error("quantity is required"),
			validation.Min(1).Error("quantity tidak boleh 0 atau minus")),
	)
}

type QuoteReceiptShippingRq struct {
	QuoteMerchantID            uint64 `json:"-"`
	ShippingProviderDurationID uint64 `json:"shipping_provider_duration_id"`
	ShippingProviderID         uint64 `json:"shipping_provider_id"`
}

func (req QuoteReceiptShippingRq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(int64(req.ShippingProviderDurationID), validation.Required.Error("shipping_provider_duration_id is required")),
	)
}

type DataQuoteReceipt struct {
	UID       string                 `json:"uid"`
	Items     []DataItemQuoteReceipt `json:"items"`
	Status    string                 `json:"status"`
	ExpiredAt time.Time              `json:"expired_at"`
}

type DataItemQuoteReceipt struct {
	Qty         int    `json:"qty"`
	Sku         string `json:"sku"`
	Uom         string `json:"uom"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Image       string `json:"image"`
	UomName     string `json:"uom_name"`
	PriceMax    int    `json:"price_max"`
	PriceMin    int    `json:"price_min"`
	IsOrdered   bool   `json:"is_ordered"`
	AturanPakai string `json:"aturan_pakai"`
}
