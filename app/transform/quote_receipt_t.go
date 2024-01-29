package transform

import (
	"context"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	entityquote "marketplace-svc/app/model/entity/quote"
	responsequote "marketplace-svc/app/model/response/quote"
	"marketplace-svc/app/repository"
	repocatalog "marketplace-svc/app/repository/catalog"
	repomerchant "marketplace-svc/app/repository/merchant"
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
	var arrQuoteMerchantID []uint64

	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quoteMerchantRepo := repoquote.NewOrderQuoteMerchantRepository(s.baseRepo)
	qms, _, err := quoteMerchantRepo.FindByParams(&dbc, map[string]interface{}{"quote_id": quote.ID}, false, 20, 1)
	if qms != nil && err == nil {
		for _, q := range *qms {
			arrQuoteMerchantID = append(arrQuoteMerchantID, q.ID)
		}
	}

	// get quote payment
	chanQuotePayment := make(chan []responsequote.QuotePaymentRs, 1)
	defer close(chanQuotePayment)
	var wg sync.WaitGroup
	// there are 3 process: get payment, merchant, address
	wg.Add(4)

	go func() {
		defer wg.Done()
		s.asyncGetQuotePayment(ctx, quote.ID, chanQuotePayment)
	}()

	// get quote merchant
	chanQuoteMerchant := make(chan []responsequote.QuoteMerchantRs, 1)
	defer close(chanQuoteMerchant)
	go func() {
		defer wg.Done()
		s.asyncGetQuoteMerchant(ctx, quote.ID, chanQuoteMerchant)
	}()

	// get quote item
	chanQuoteItem := make(chan []responsequote.QuoteItemRs, 1)
	defer close(chanQuoteItem)
	go func() {
		defer wg.Done()
		s.asyncGetQuoteItem(ctx, arrQuoteMerchantID, chanQuoteItem)
	}()

	// get quote address
	chanQuoteAddress := make(chan []responsequote.QuoteAddressRs, 1)
	defer close(chanQuoteAddress)
	go func() {
		defer wg.Done()
		s.asyncGetQuoteAddress(ctx, quote.ID, chanQuoteAddress)
	}()

	// wait all and get channel value
	wg.Wait()
	quoteMerchant := <-chanQuoteMerchant
	quoteItems := <-chanQuoteItem
	quotePayment := <-chanQuotePayment
	quoteAddress := <-chanQuoteAddress

	if len(quoteMerchant) > 0 {
		// implement time complexity O(log n)
		// space complexity O(n)
		quoteItemMap := map[uint64][]responsequote.QuoteItemRs{}
		if quoteItems != nil {
			for _, item := range quoteItems {
				if _, ok := quoteItemMap[item.QuoteMerchantID]; !ok {
					quoteItemMap[item.QuoteMerchantID] = append([]responsequote.QuoteItemRs{}, item)
					continue
				}
				quoteItemMap[item.QuoteMerchantID] = append(quoteItemMap[item.QuoteMerchantID], item)
			}
			// assign quote item to quote merchant
			for i := 0; i < len(quoteMerchant); i++ {
				if quoteItem, ok := quoteItemMap[quoteMerchant[i].ID]; ok {
					quoteMerchant[i].OrderQuoteItems = &quoteItem
				}
			}
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
		DiscountAmount:     totalDiscount,
		GrandTotal:         grandTotal,
		Currency:           "IDR",
		TotalQuantity:      int(totalQty),
		CouponCode:         quote.CouponCode,
		OrderTypeID:        quote.OrderTypeID,
		OrderQuotePayment:  &quotePayment,
		OrderQuoteMerchant: &quoteMerchant,
		OrderQuoteAddress:  &quoteAddress,
		OrderQuoteReceipt:  quote.DataReceipt,
	}

	return &quoteRs
}

func (s QuoteReceiptTransform) asyncGetQuotePayment(ctx context.Context, quoteID uint64, chanQ chan<- []responsequote.QuotePaymentRs) {
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

func (s QuoteReceiptTransform) asyncGetQuoteAddress(ctx context.Context, quoteID uint64, chanQ chan<- []responsequote.QuoteAddressRs) {
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

func (s QuoteReceiptTransform) asyncGetQuoteMerchant(ctx context.Context, quoteID uint64, chanQ chan<- []responsequote.QuoteMerchantRs) {
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

	chanQ <- quoteMerchants
	return
}

func (s QuoteReceiptTransform) asyncGetQuoteItem(ctx context.Context, arrQuoteMerchantID []uint64, chanQ chan<- []responsequote.QuoteItemRs) {
	// check count arrQuoteMerchantID
	if len(arrQuoteMerchantID) == 0 {
		chanQ <- nil
		return
	}

	dbc := repository.DBContext{DB: s.baseRepo.GetDB(), Context: ctx}
	quoteItemRepo := repoquote.NewOrderQuoteItemRepository(s.baseRepo)
	quoteItem, err := quoteItemRepo.FindRawByQuoteMerchantID(&dbc, arrQuoteMerchantID)
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
				mpRepo := repomerchant.NewMerchantProductRepository(s.baseRepo)
				pcRepo := repocatalog.NewProductCategoryRepository(s.baseRepo)

				filterMp := map[string]interface{}{
					"merchant_id":  qi.Merchant.ID,
					"product_sku":  qi.ProductSku,
					"merchant_sku": qi.MerchantSku,
				}
				mp, err := mpRepo.FindFirstByParams(&dbc, filterMp, true)
				if err != nil {
					return
				}
				arrCategory, _ := pcRepo.GetCategoryMenu(&dbc, qi.ProductID, 1)
				image := ""
				if qi.Product.ProductImage != nil {
					firstProductImage := *qi.Product.ProductImage
					image = s.infra.Config.URL.BaseImageURL + firstProductImage[0].ImageThumbnail
				}

				chanQI <- *responsequote.QuoteItemRs{}.Transform(&qi, *qi.Merchant, *mp, arrCategory, image)
			}(qi)
		}
		wg.Wait()

		// get channel
		for i := 0; i < count; i++ {
			quoteItems = append(quoteItems, <-chanQI)
		}
	}

	chanQ <- quoteItems
	return
}
