package responsequote

import entity "marketplace-svc/app/model/entity/quote"

type QuotePaymentRs struct {
	PaymentMethodID       uint64  `json:"payment_method_id,omitempty"`
	PaymentMethodName     string  `json:"payment_method_name,omitempty"`
	PaymentStatus         int     `json:"payment_status,omitempty"`
	MinimumAmount         float64 `json:"minimum_amount,omitempty"`
	PaymentMethodTypeName string  `json:"payment_method_type_name,omitempty"`
	PaymentMethodTypeID   uint64  `json:"payment_method_type_id,omitempty"`
	PaymentLogo           string  `json:"payment_logo,omitempty"`
}

func (qr QuotePaymentRs) Transform(qp *entity.OrderQuotePayment, baseImageURL string) []QuotePaymentRs {
	var response []QuotePaymentRs //nolint:prealloc
	if qp == nil || qp.ID == 0 {
		return nil
	}

	// set response
	response = append(response, QuotePaymentRs{
		PaymentMethodID:       qp.PaymentMethodID,
		PaymentMethodName:     qp.PaymentMethod.Name,
		PaymentStatus:         qp.PaymentMethod.Status,
		MinimumAmount:         qp.PaymentMethod.MinimumAmount,
		PaymentMethodTypeName: qp.PaymentMethod.PaymentMethodType.Name,
		PaymentMethodTypeID:   qp.PaymentMethod.PaymentMethodTypeID,
		PaymentLogo:           baseImageURL + qp.PaymentMethod.Logo,
	})

	return response
}
