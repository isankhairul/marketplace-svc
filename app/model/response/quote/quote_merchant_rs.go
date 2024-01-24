package responsequote

import (
	"fmt"
	"marketplace-svc/app"
	entity "marketplace-svc/app/model/entity/quote"
)

type QuoteMerchantRs struct {
	MerchantID               int64            `json:"merchant_id,omitempty"`
	MerchantTypeID           int              `json:"merchant_type_id,omitempty"`
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
	OrderQuoteItems          *[]QuoteItemRs   `json:"orderQuoteItems,omitempty"`
	OrderQuoteShipping       *QuoteShippingRs `json:"orderQuoteShipping,omitempty"`
}

func (qr QuoteMerchantRs) Transform(qms *[]entity.OrderQuoteMerchant, infra app.Infra) *[]QuoteMerchantRs {
	var response []QuoteMerchantRs
	if qms == nil {
		return nil
	}
	for _, qm := range *qms {
		response = append(response, QuoteMerchantRs{
			MerchantID:             qm.MerchantID,
			MerchantTypeID:         qm.MerchantTypeID,
			Image:                  infra.Config.URL.BaseImageURL + qm.Merchant.Image,
			Slug:                   qm.Merchant.Slug,
			Status:                 qm.Merchant.Status,
			Name:                   qm.Merchant.MerchantName,
			Code:                   qm.Merchant.MerchantCode,
			Province:               fmt.Sprint(qm.Merchant.ProvinceID),
			City:                   fmt.Sprint(qm.Merchant.CityID),
			MerchantGrandTotal:     qm.MerchantGrandTotal,
			MerchantTotalQuantity:  qm.MerchantTotalQuantity,
			MerchantTotalWeight:    qm.MerchantTotalWeight,
			MerchantNotes:          qm.MerchantNotes,
			Selected:               qm.Selected,
			PromoDescription:       qm.PromoDescription,
			DiscountAmount:         qm.DiscountAmount,
			ShippingDiscountAmount: qm.ShippingDiscountAmount,
			AdminFeeCalculation:    qm.AdminFeeCalculation,
			OrderQuoteItems:        QuoteItemRs{}.Transform(qm.OrderQuoteItem, qm.Merchant, infra),
			OrderQuoteShipping:     QuoteShippingRs{}.Transform(qm.OrderQuoteShipping, infra.Config.URL.BaseImageURL),
		})
	}

	return &response
}
