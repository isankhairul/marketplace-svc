package transform

import (
	"context"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	entityquote "marketplace-svc/app/model/entity/quote"
	responsequote "marketplace-svc/app/model/response/quote"
	"marketplace-svc/app/repository"
	repoquote "marketplace-svc/app/repository/quote"
	"sync"
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
	return &QuoteReceiptTransform{infra, br, quoteRepo, base.ORDER_TYPE_RECEIPT}
}

func (s QuoteReceiptTransform) TransformQuote(ctx context.Context, quote *entityquote.OrderQuote) *responsequote.QuoteRs {
	if nil == quote {
		return nil
	}
	var subTotal, totalQty, totalPointEarned, totalPointSpent, totalWeight, totalDiscount, grandTotal float64

	// get quote payment
	chanQuotePayment := make(chan *[]responsequote.QuotePaymentRs, 1)
	defer close(chanQuotePayment)
	var wg sync.WaitGroup
	// there are 3 process: get payment, merchant, address
	wg.Add(3)

	go func(ctx context.Context, quoteID uint64) {
		defer wg.Done()
		s.asyncGetQuotePayment(ctx, quoteID, chanQuotePayment)
	}(ctx, quote.ID)

	// get quote merchant
	chanQuoteMerchant := make(chan *[]responsequote.QuoteMerchantRs, 1)
	defer close(chanQuoteMerchant)
	go func(ctx context.Context, quoteID uint64) {
		defer wg.Done()
		s.asyncGetQuoteMerchant(ctx, quoteID, chanQuoteMerchant)
	}(ctx, quote.ID)

	// get quote address
	chanQuoteAddress := make(chan *[]responsequote.QuoteAddressRs, 1)
	defer close(chanQuoteAddress)
	go func(ctx context.Context, quoteID uint64) {
		defer wg.Done()
		s.asyncGetQuoteAddress(ctx, quoteID, chanQuoteAddress)
	}(ctx, quote.ID)

	// wait all and get channel value
	wg.Wait()
	quoteMerchant := <-chanQuoteMerchant //*[]responsequote.QuoteMerchantRs
	quotePayment := <-chanQuotePayment
	quoteAddress := <-chanQuoteAddress

	var quoteItems []responsequote.QuoteItemRs
	if quoteMerchant != nil && len(*quoteMerchant) > 0 {
		for _, quoteMerchant := range *quoteMerchant {
			quoteItems = append(quoteItems, *quoteMerchant.OrderQuoteItems...)
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
		QuoteCode:          quote.QuoteCode,
		DeviceID:           quote.DeviceID,
		StoreID:            quote.StoreID,
		CustomerID:         quote.CustomerID,
		CustomerEmail:      quote.CustomerEmail,
		CustomerFirstname:  quote.CustomerFirstname,
		CustomerLastname:   quote.CustomerLastname,
		Weight:             totalWeight,
		Subtotal:           subTotal,
		ShippingAmount:     quote.ShippingAmount,
		GrandTotal:         grandTotal,
		Currency:           "IDR",
		Status:             quote.Status,
		OrderTypeID:        quote.OrderTypeID,
		TotalQuantity:      int(totalQty),
		CustomerGroupID:    quote.CustomerGroupID,
		CreatedAt:          quote.CreatedAt,
		UpdatedAt:          quote.UpdatedAt,
		OrderQuotePayment:  quotePayment,
		OrderQuoteMerchant: quoteMerchant,
		OrderQuoteAddress:  quoteAddress,
		OrderQuoteReceipt:  responsequote.QuoteReceiptRs{}.Transform(quote.DataReceipt),
	}

	return &quoteRs
}

func (s QuoteReceiptTransform) asyncGetQuotePayment(ctx context.Context, quoteID uint64, chanQ chan<- *[]responsequote.QuotePaymentRs) {
	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quotePaymentRepo := repoquote.NewOrderQuotePaymentRepository(s.baseRepo)
	filter := map[string]interface{}{
		"quote_id": quoteID,
	}
	quotePayment, err := quotePaymentRepo.FindFirstByParams(&dbc, filter, true)
	if err != nil {
		chanQ <- nil
		return
	}
	chanQ <- responsequote.QuotePaymentRs{}.Transform(quotePayment, s.infra.Config.URL.BaseImageURL)
	return
}

func (s QuoteReceiptTransform) asyncGetQuoteAddress(ctx context.Context, quoteID uint64, chanQ chan<- *[]responsequote.QuoteAddressRs) {
	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quoteAddressRepo := repoquote.NewOrderQuoteAddressRepository(s.baseRepo)
	filter := map[string]interface{}{
		"quote_id": quoteID,
	}
	quoteAddress, err := quoteAddressRepo.FindFirstByParams(&dbc, filter)
	if err != nil || quoteAddress == nil {
		chanQ <- nil
		return
	}
	chanQ <- responsequote.QuoteAddressRs{}.Transform(quoteAddress)
	return
}

func (s QuoteReceiptTransform) asyncGetQuoteMerchant(ctx context.Context, quoteID uint64, chanQ chan<- *[]responsequote.QuoteMerchantRs) {
	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quoteMerchantRepo := repoquote.NewOrderQuoteMerchantRepository(s.baseRepo)
	filter := map[string]interface{}{
		"quote_id": quoteID,
	}
	quoteMerchant, _, err := quoteMerchantRepo.FindByParams(&dbc, filter, true, 100, 1)
	if err != nil || quoteMerchant == nil {
		chanQ <- nil
		return
	}
	var quoteMerchants []responsequote.QuoteMerchantRs
	if quoteMerchant != nil && len(*quoteMerchant) > 0 {
		lenMerchant := len(*quoteMerchant)
		chanQM := make(chan responsequote.QuoteMerchantRs, lenMerchant)
		defer close(chanQM)
		var wg sync.WaitGroup
		wg.Add(lenMerchant)
		for _, qm := range *quoteMerchant {
			go func(qm entityquote.OrderQuoteMerchant) {
				defer wg.Done()
				quoteMerchantRs := *responsequote.QuoteMerchantRs{}.Transform(&qm, s.infra)

				// get quote item
				chanQI := make(chan *[]responsequote.QuoteItemRs, 1)
				defer close(chanQI)
				go func() {
					s.asyncGetQuoteItem(ctx, quoteMerchantRs.ID, qm.Merchant, chanQI)
				}()
				quoteMerchantRs.OrderQuoteItems = <-chanQI

				chanQM <- quoteMerchantRs
				return
			}(qm)
		}
		wg.Wait()

		// get channel
		for i := 0; i < lenMerchant; i++ {
			quoteMerchants = append(quoteMerchants, <-chanQM)
		}
	}

	chanQ <- &quoteMerchants
	return
}

func (s QuoteReceiptTransform) asyncGetQuoteItem(ctx context.Context, quoteMerchantID uint64, merchant entitymerchant.Merchant, chanQ chan<- *[]responsequote.QuoteItemRs) {
	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quoteItemRepo := repoquote.NewOrderQuoteItemRepository(s.baseRepo)
	filter := map[string]interface{}{
		"quote_merchant_id": quoteMerchantID,
	}
	quoteItem, _, err := quoteItemRepo.FindByParams(&dbc, filter, true, 100, 1)
	if err != nil || quoteItem == nil {
		chanQ <- nil
		return
	}
	var quoteItems []responsequote.QuoteItemRs
	if quoteItem != nil && len(*quoteItem) > 0 {
		count := len(*quoteItem)
		chanQI := make(chan responsequote.QuoteItemRs, count)
		defer close(chanQI)
		var wg sync.WaitGroup
		wg.Add(count)
		for _, qi := range *quoteItem {
			go func(qi entityquote.OrderQuoteItem) {
				defer wg.Done()
				chanQI <- *responsequote.QuoteItemRs{}.Transform(&qi, merchant, s.infra)
			}(qi)
		}
		wg.Wait()

		// get channel
		for i := 0; i < count; i++ {
			quoteItems = append(quoteItems, <-chanQI)
		}
	}

	chanQ <- &quoteItems
	return
}
