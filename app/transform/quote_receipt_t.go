package transform

import (
	"context"
	"marketplace-svc/app"
	entityquote "marketplace-svc/app/model/entity/quote"
	responsequote "marketplace-svc/app/model/response/quote"
	"marketplace-svc/app/repository"
	repoquote "marketplace-svc/app/repository/quote"
	helperconst "marketplace-svc/helper/const"
)

type QuoteReceiptTransform struct {
	infra       app.Infra
	baseRepo    repository.BaseRepository
	quoteRepo   repoquote.OrderQuoteRepository
	OrderTypeID int
}

func NewQuoteReceiptTransform(
	infra app.Infra,
	br repository.BaseRepository,
	quoteRepo repoquote.OrderQuoteRepository,
) *QuoteReceiptTransform {
	return &QuoteReceiptTransform{infra, br, quoteRepo, helperconst.ORDER_TYPE_RECEIPT}
}

func (s QuoteReceiptTransform) TransformQuote(ctx context.Context, quote *entityquote.OrderQuote) *responsequote.QuoteRs {
	if nil == quote {
		return nil
	}
	var subTotal, totalQty, totalPointEarned, totalPointSpent, totalWeight, totalDiscount, grandTotal float64
	var quoteItems []entityquote.OrderQuoteItem
	if quote.OrderQuoteMerchant != nil && len(*quote.OrderQuoteMerchant) > 0 {
		for _, quoteMerchant := range *quote.OrderQuoteMerchant {
			quoteItems = append(quoteItems, *quoteMerchant.OrderQuoteItem...)
		}
	}

	// calculate total on quote item
	for _, quoteItem := range quoteItems {
		if quoteItem.Selected {
			subTotal += float64(float64(quoteItem.Quantity) * quoteItem.Price)
			totalQty += float64(quoteItem.Quantity)
			totalPointEarned += float64(quoteItem.PointEarned)
			totalPointSpent += float64(quoteItem.PointSpent)
			totalWeight += float64(float64(quoteItem.Quantity) * quoteItem.Weight)
			totalDiscount += float64(quoteItem.DiscountAmount)
		}
	}
	grandTotal = (subTotal + quote.ShippingAmount + float64(quote.InsuranceAmount)) - totalDiscount - quote.ShippingDiscountAmount

	quoteRs := responsequote.QuoteRs{
		QuoteCode:         quote.QuoteCode,
		DeviceID:          quote.DeviceID,
		StoreID:           quote.StoreID,
		CustomerID:        quote.CustomerID,
		CustomerEmail:     quote.CustomerEmail,
		CustomerFirstname: quote.CustomerFirstname,
		CustomerLastname:  quote.CustomerLastname,
		Weight:            totalWeight,
		Subtotal:          subTotal,
		ShippingAmount:    quote.ShippingAmount,
		GrandTotal:        grandTotal,
		Currency:          "IDR",
		Status:            quote.Status,
		OrderTypeID:       quote.OrderTypeID,
		TotalQuantity:     int(totalQty),
		CustomerGroupID:   quote.CustomerGroupID,
		CreatedAt:         quote.CreatedAt,
		UpdatedAt:         quote.UpdatedAt,
		OrderQuotePayment: responsequote.QuotePaymentRs{}.Transform(quote.OrderQuotePayment, s.infra.Config.URL.BaseImageURL),
		OrderQuoteAddress: responsequote.QuoteAddressRs{}.Transform(quote.OrderQuoteAddress),
	}
	if quote.OrderQuoteMerchant != nil {
		quoteRs.OrderQuoteMerchant = responsequote.QuoteMerchantRs{}.Transform(quote.OrderQuoteMerchant, s.infra)
	}

	return &quoteRs
}
