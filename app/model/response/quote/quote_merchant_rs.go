package responsequote

import (
	"marketplace-svc/app"
	entity "marketplace-svc/app/model/entity/quote"
)

type QuoteMerchantRs struct {
	ID                       uint64           `json:"ID,omitempty"`
	MerchantID               uint64           `json:"merchant_id,omitempty"`
	MerchantTypeID           uint64           `json:"merchant_type_id,omitempty"`
	Image                    string           `json:"image,omitempty"`
	Slug                     string           `json:"slug,omitempty"`
	Status                   int              `json:"status,omitempty"`
	StatusName               string           `json:"status_name,omitempty"`
	Name                     string           `json:"name,omitempty"`
	Code                     string           `json:"code,omitempty"`
	Province                 string           `json:"province,omitempty"`
	City                     string           `json:"city,omitempty"`
	MerchantGrandTotal       float64          `json:"merchant_grand_total,omitempty"`
	MerchantSubtotal         float64          `json:"merchant_subtotal,omitempty"`
	MerchantTotalQuantity    int              `json:"merchant_total_quantity,omitempty"`
	MerchantTotalWeight      float64          `json:"merchant_total_weight,omitempty"`
	MerchantTotalPointEarned float64          `json:"merchant_total_point_earned,omitempty"`
	MerchantTotalPointSpent  float64          `json:"merchant_total_point_spent,omitempty"`
	MerchantNotes            string           `json:"merchant_notes,omitempty"`
	DiscountAmount           float64          `json:"discount_amount,omitempty"`
	ShippingDiscountAmount   float64          `json:"shipping_discount_amount,omitempty"`
	AdminFeeCalculation      float64          `json:"admin_fee_calculation,omitempty"`
	PromoDescription         string           `json:"promo_description,omitempty"`
	AppliedRuleIds           string           `json:"applied_rule_ids,omitempty"`
	Selected                 bool             `json:"selected,omitempty"`
	MerchantWarning          []interface{}    `json:"merchant_warning,omitempty"`
	StockStatus              int              `json:"stock_status,omitempty"`
	OrderQuoteItems          *[]QuoteItemRs   `json:"orderQuoteItems"`
	OrderQuoteShipping       *QuoteShippingRs `json:"orderQuoteShipping"`
}

func (qr QuoteMerchantRs) Transform(qm *entity.OrderQuoteMerchant, infra app.Infra) *QuoteMerchantRs {
	var response QuoteMerchantRs
	if qm == nil {
		return nil
	}
	response = QuoteMerchantRs{
		ID:                     qm.ID,
		MerchantID:             qm.MerchantID,
		MerchantTypeID:         qm.MerchantTypeID,
		Image:                  infra.Config.URL.BaseImageURL + qm.Merchant.Image,
		Slug:                   qm.Merchant.Slug,
		Status:                 qm.Merchant.Status,
		Name:                   qm.Merchant.MerchantName,
		Code:                   qm.Merchant.MerchantCode,
		Province:               qm.Merchant.Province.Name,
		City:                   qm.Merchant.City.Name,
		MerchantGrandTotal:     qm.MerchantGrandTotal,
		MerchantSubtotal:       qm.MerchantSubtotal,
		MerchantTotalQuantity:  qm.MerchantTotalQuantity,
		MerchantTotalWeight:    qm.MerchantTotalWeight,
		MerchantNotes:          qm.MerchantNotes,
		Selected:               qm.Selected,
		PromoDescription:       qm.PromoDescription,
		DiscountAmount:         qm.DiscountAmount,
		ShippingDiscountAmount: qm.ShippingDiscountAmount,
		AdminFeeCalculation:    qm.AdminFeeCalculation,
		OrderQuoteShipping:     QuoteShippingRs{}.Transform(qm.OrderQuoteShipping, infra.Config.URL.BaseImageURL),
	}

	return &response
}
