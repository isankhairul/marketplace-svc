package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-kit/kit/auth/jwt"
	"marketplace-svc/app"
	"marketplace-svc/app/api/middleware"
	"marketplace-svc/app/model/base"
	entityquote "marketplace-svc/app/model/entity/quote"
	"marketplace-svc/app/model/request"
	responsequote "marketplace-svc/app/model/response/quote"
	"marketplace-svc/app/repository"
	repomerchant "marketplace-svc/app/repository/merchant"
	repopromo "marketplace-svc/app/repository/promotion"
	repoquote "marketplace-svc/app/repository/quote"
	"marketplace-svc/app/transform"
	helper_kalcare "marketplace-svc/helper/kalcare"
	"marketplace-svc/helper/message"
	"marketplace-svc/pkg/util"
	"strconv"
)

type QuoteReceiptService interface {
	Find(ctx context.Context, quoteCode string, validate bool) (*responsequote.QuoteRs, message.Message, error)
	CheckQuote(ctx context.Context, quoteCode string, quote *entityquote.OrderQuote) (*entityquote.OrderQuote, message.Message, error)
	Save(ctx context.Context, input request.QuoteReceiptRq) (message.Message, error)
}

type QuoteReceiptServiceImpl struct {
	infra               app.Infra
	baseRepo            repository.BaseRepository
	quoteRepo           repoquote.OrderQuoteRepository
	quoteMerchantRepo   repoquote.OrderQuoteMerchantRepository
	quoteItemRepo       repoquote.OrderQuoteItemRepository
	quoteShippingRepo   repoquote.OrderQuoteShippingRepository
	quoteAddressRepo    repoquote.OrderQuoteAddressRepository
	quotePaymentRepo    repoquote.OrderQuotePaymentRepository
	merchantRepo        repomerchant.MerchantRepository
	merchantProductRepo repomerchant.MerchantProductRepository
	OrderTypeID         int
}

func NewQuoteReceiptService(
	infra app.Infra,
	br repository.BaseRepository,
	quoteRepo repoquote.OrderQuoteRepository,
	quoteMerchantRepo repoquote.OrderQuoteMerchantRepository,
	quoteItemRepo repoquote.OrderQuoteItemRepository,
	quoteShippingRepo repoquote.OrderQuoteShippingRepository,
	quoteAddressRepo repoquote.OrderQuoteAddressRepository,
	quotePaymentRepo repoquote.OrderQuotePaymentRepository,
	merchantRepo repomerchant.MerchantRepository,
	merchantProductRepo repomerchant.MerchantProductRepository,
) QuoteReceiptService {
	return &QuoteReceiptServiceImpl{infra, br, quoteRepo,
		quoteMerchantRepo, quoteItemRepo, quoteShippingRepo,
		quoteAddressRepo, quotePaymentRepo, merchantRepo, merchantProductRepo,
		base.ORDER_TYPE_RECEIPT}
}

func (s *QuoteReceiptServiceImpl) CheckQuote(ctx context.Context, quoteCode string, quote *entityquote.OrderQuote) (*entityquote.OrderQuote, message.Message, error) {
	errMsgPrefix := "QUOTE-CHECKQUOTE"
	deviceID := 1
	ctxDeviceID, err := strconv.ParseInt(fmt.Sprint(ctx.Value("device_id")), 10, 8)
	if ctxDeviceID != 0 && err == nil {
		deviceID = int(ctxDeviceID)
	}

	user, isLogged := middleware.IsAuthContext(ctx)
	if !isLogged {
		return quote, message.ErrNoAuth, errors.New(message.ErrNoAuth.Message)
	}

	filter := map[string]interface{}{
		"quote_code":    quoteCode,
		"order_type_id": s.OrderTypeID,
		"customer_id":   user.CustomerID,
	}
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}
	if quote == nil {
		quoteRs, err := s.quoteRepo.FindFirstByParams(&dbc, filter, false)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " " + err.Error()))
			return nil, message.ErrDB, err
		}

		if quoteRs == nil {
			return nil, message.ErrNoData, err
		}

		quote = quoteRs
	}
	// update device_id
	quote.DeviceID = uint8(deviceID)
	dataUpdate := map[string]interface{}{
		"device_id": quote.DeviceID,
	}
	err = s.quoteRepo.UpdateMapByQuoteCode(&dbc, quoteCode, dataUpdate)
	if err != nil {
		fmt.Println("err", err.Error())
		return nil, message.ErrDB, err
	}

	return quote, message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) Find(ctx context.Context, quoteCode string, validate bool) (*responsequote.QuoteRs, message.Message, error) {
	var response *responsequote.QuoteRs
	quote, msg, err := s.CheckQuote(ctx, quoteCode, nil)
	if err != nil {
		return response, msg, err
	}
	response = transform.NewQuoteReceiptTransform(s.infra, s.baseRepo, s.quoteRepo).TransformQuote(ctx, quote)
	return response, message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) Save(ctx context.Context, input request.QuoteReceiptRq) (message.Message, error) {
	msg := message.SuccessMsg
	quoteCode := input.QuoteCode
	errMsgPrefix := "QUOTE-SAVE"
	quote, msg, err := s.CheckQuote(ctx, quoteCode, nil)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quote " + err.Error()))
		return msg, err
	}

	// get user
	user, _ := middleware.IsAuthContext(ctx)

	// check qty valid
	if !s.isValidQty(ctx, input) {
		return message.QuoteQtyInvalid, errors.New(message.QuoteQtyInvalid.Message)
	}

	if len(input.OrderQuoteMerchants) > 0 {
		// map index use order_quote_merchant_id
		quoteItemsMap := map[uint64]map[string]float64{}
		quoteShippingMap := map[uint64]map[string]float64{}
		var arrQuoteItemsRq []request.QuoteReceiptItemsRq
		var arrQuoteShippingRq []request.QuoteReceiptShippingRq
		dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}

		// process quote merchant
		for _, oqm := range input.OrderQuoteMerchants {
			// validation
			if err := oqm.Validate(); err != nil {
				return message.ValidationError, err
			}

			filterQuoteMerchant := map[string]interface{}{"quote_id": quote.ID, "merchant_id": oqm.MerchantID}
			quoteMerchant, err := s.quoteMerchantRepo.FindFirstByParams(&dbc, filterQuoteMerchant, false)
			if err != nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quoteMerchant " + err.Error()))
				return message.ErrDB, errors.New(err.Error())
			}

			// if empty data then create new
			if quoteMerchant == nil {
				quoteMerchant = &entityquote.OrderQuoteMerchant{}
			}
			merchant, err := s.merchantRepo.FindFirstByParams(&dbc, map[string]interface{}{"id": oqm.MerchantID}, false)
			if err != nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " " + err.Error()))
				return message.ErrDB, errors.New(err.Error())
			}

			// merchant not found
			if merchant == nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + "merchant not found, merchant_id:" + fmt.Sprint(oqm.MerchantID)))
				return message.ErrNoData, errors.New("merchant not found, merchant_id:" + fmt.Sprint(oqm.MerchantID))
			}

			// assign to quote_merchant
			quoteMerchant.QuoteID = quote.ID
			quoteMerchant.MerchantID = merchant.ID
			quoteMerchant.Selected = true
			quoteMerchant, err = s.quoteMerchantRepo.Save(&dbc, quoteMerchant)
			if err != nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quoteMerchant " + err.Error()))
				return message.ErrDB, errors.New(err.Error())
			}

			// add to array quote shipping
			if oqm.OrderQuoteShipping != nil {
				oqm.OrderQuoteShipping.QuoteMerchantID = quoteMerchant.ID
				oqm.OrderQuoteShipping.MerchantID = quoteMerchant.MerchantID
				arrQuoteShippingRq = append(arrQuoteShippingRq, *oqm.OrderQuoteShipping)
			}

			// to flatten quoteItemsRq
			for _, item := range oqm.OrderQuoteItems {
				item.QuoteMerchantID = quoteMerchant.ID
				item.MerchantID = quoteMerchant.MerchantID
				arrQuoteItemsRq = append(arrQuoteItemsRq, item)
			}
		}

		// process quote items
		if len(arrQuoteItemsRq) > 0 {
			for _, qi := range arrQuoteItemsRq {
				// validation
				if err := qi.Validate(); err != nil {
					return message.ValidationError, err
				}

				filterQuoteItem := map[string]interface{}{
					"quote_merchant_id": qi.QuoteMerchantID,
					"merchant_sku":      qi.MerchantSku,
				}
				quoteItem, err := s.quoteItemRepo.FindFirstByParams(&dbc, filterQuoteItem, false)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}
				filterProduct := map[string]interface{}{
					"merchant_id":  qi.MerchantID,
					"merchant_sku": qi.MerchantSku,
					"store_id":     1,
				}
				product, err := s.merchantProductRepo.FindFirstDetailMerchantProduct(&dbc, filterProduct)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}

				// product not found
				if product == nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " " + message.QuoteProductNotFound.Message))
					return message.QuoteProductNotFound, errors.New(message.QuoteProductNotFound.Message)
				}

				if quoteItem == nil {
					quoteItem = &entityquote.OrderQuoteItem{}
				}
				// SET QUOTE ITEM
				quoteItem.QuoteID = quote.ID
				quoteItem.QuoteMerchantID = qi.QuoteMerchantID
				quoteItem.ProductID = product.ID
				quoteItem.ProductSku = product.Sku
				quoteItem.MerchantSku = product.MerchantSku
				quoteItem.Name = product.Name
				quoteItem.Price = product.SellingPrice
				quoteItem.OriginalPrice = product.SellingPrice
				quoteItem.BasePrice = product.BasePrice
				quoteItem.ProductType = product.ProductType
				quoteItem.Quantity = qi.Quantity
				quoteItem.BrandName = product.BrandCode
				quoteItem.ItemNotes = qi.ItemNotes
				quoteItem.Quantity = qi.Quantity
				quoteItem.Weight = product.Weight
				quoteItem.Selected = true

				// start get promo price
				pcppRepo := repopromo.NewPromotionCatalogProductPriceRepository(s.baseRepo)

				dateTimeNow := util.TimeNow()
				dateTimeNowStr := dateTimeNow.Format(util.LayoutDefault)
				filterPromoCatalogProductPrice := map[string]interface{}{
					"product_id":        quoteItem.ProductID,
					"customer_group_id": user.GroupID,
					"merchant_id":       qi.MerchantID,
					"store_id":          quote.StoreID,
					"latest_start_date": dateTimeNowStr,
					"earliest_end_date": dateTimeNowStr,
				}
				pcpp, err := pcppRepo.FindFirstByParams(&dbc, filterPromoCatalogProductPrice)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get promo_catalog_product_price " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}
				var specialPrice float64
				var promoDescription interface{}
				if pcpp != nil {
					pcRepo := repopromo.NewPromotionCatalogRepository(s.baseRepo)
					rule, err := pcRepo.FindFirstByParams(&dbc, map[string]interface{}{"id": pcpp.PromotionCatalogID})
					if err != nil {
						s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get promo_catalog " + err.Error()))
						return message.ErrDB, errors.New(err.Error())
					}
					if rule != nil {
						specialPrice = pcpp.RulePrice
						promoDescription = map[string]interface{}{
							"quote_item_id":        quoteItem.ID,
							"type":                 "regular",
							"rule_id":              rule.ID,
							"rule_name":            rule.Name,
							"principal_data":       rule.PrincipalData,
							"discount_type":        rule.SimpleAction,
							"discount_type_amount": rule.DiscountAmount,
							"price":                quoteItem.Price,
							"original_price":       quoteItem.OriginalPrice,
							"base_price":           quoteItem.BasePrice,
						}
					}
				}
				if product.SpecialPrice != 0 &&
					(product.SpecialPriceStartTime != nil && dateTimeNow.After(*product.SpecialPriceStartTime) &&
						(product.SpecialPriceEndTime != nil && product.SpecialPriceEndTime.Before(dateTimeNow))) {
					if specialPrice > 0 && product.SpecialPrice < specialPrice {
						specialPrice = product.SpecialPrice
						promoDescription = map[string]interface{}{}
					}
				}
				// get special price
				if specialPrice > 0 {
					// set quote item price use promo
					jsonPromoDescription, _ := sonic.Marshal(promoDescription)
					quoteItem.Price = specialPrice
					quoteItem.PromoDescription = string(jsonPromoDescription)
				}
				// end get promo price

				quoteItem.RowTotal = float64(quoteItem.Quantity) * quoteItem.Price
				quoteItem.RowWeight = float64(quoteItem.Quantity) * quoteItem.Weight
				quoteItem.RowOriginalPrice = float64(quoteItem.Quantity) * quoteItem.OriginalPrice
				_, err = s.quoteItemRepo.Save(&dbc, quoteItem)

				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quoteItem " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}

				// add to map
				if items, ok := quoteItemsMap[qi.QuoteMerchantID]; ok {
					items["qty"] += float64(quoteItem.Quantity)
					items["price"] += quoteItem.Price
					items["weight"] += quoteItem.Weight
					items["row_weight"] += quoteItem.RowWeight
					items["row_total"] += quoteItem.RowTotal

				} else {
					mapItems := map[string]float64{
						"qty":        float64(quoteItem.Quantity),
						"price":      quoteItem.Price,
						"weight":     quoteItem.Weight,
						"row_weight": quoteItem.RowWeight,
						"row_total":  quoteItem.RowTotal,
					}
					quoteItemsMap[qi.QuoteMerchantID] = mapItems
				}
			}

			// trigger remove quote shipping
			if input.OrderQuoteMerchants[0].OrderQuoteShipping == nil {
				_, _ = s.RemoveQuoteShipping(ctx, quote.ID)
			}
		}

		// process calculate quote shipping after finish process quote item
		if len(arrQuoteShippingRq) > 0 {
			dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}
			for _, qs := range arrQuoteShippingRq {
				quoteMerchant, _ := s.quoteMerchantRepo.FindFirstByParams(&dbc, map[string]interface{}{"id": qs.QuoteMerchantID}, false)
				if quoteMerchant != nil {
					responseOQS, msg, err := s.ProcessQuoteShipping(ctx, &qs, *quoteMerchant)
					if err != nil {
						return msg, err
					}
					// add to quoteShippingMap
					if responseOQS != nil {
						quoteShippingMap = *responseOQS
					}
				}
			}
		}

		// calculate quote merchant total
		for quoteMerchantID, items := range quoteItemsMap {
			if rowTotal, ok := items["row_total"]; ok {
				filterQuoteMerchant := map[string]interface{}{"id": quoteMerchantID, "quote_id": quote.ID}
				quoteMerchant, err := s.quoteMerchantRepo.FindFirstByParams(&dbc, filterQuoteMerchant, false)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quoteMerchant " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}
				quoteMerchant.MerchantGrandTotal = rowTotal
				quoteMerchant.MerchantSubtotal = rowTotal
				quoteMerchant.MerchantTotalQuantity = int(items["qty"])
				quoteMerchant.MerchantTotalWeight = items["row_weight"]
				_, err = s.quoteMerchantRepo.Save(&dbc, quoteMerchant)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quoteMerchant " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}
			}
		}

		// calculate quote shipping amount
		for quoteMerchantID, items := range quoteShippingMap {
			if shippingAmount, ok := items["shipping_amount"]; ok {
				filterQuoteMerchant := map[string]interface{}{"id": quoteMerchantID, "quote_id": quote.ID}
				quoteMerchant, err := s.quoteMerchantRepo.FindFirstByParams(&dbc, filterQuoteMerchant, false)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quoteMerchant " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}
				quoteMerchant.MerchantGrandTotal = quoteMerchant.MerchantGrandTotal + shippingAmount
				_, err = s.quoteMerchantRepo.Save(&dbc, quoteMerchant)
				if err != nil {
					s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quoteMerchant " + err.Error()))
					return message.ErrDB, errors.New(err.Error())
				}

			}
		}
	}

	// process quote address
	if input.OrderQuoteAddress != nil {
		// validation
		if err := input.OrderQuoteAddress.Validate(); err != nil {
			return message.ValidationError, err
		}

		dbc := repository.DBContext{Context: ctx, DB: s.baseRepo.GetDB()}
		caRepo := repository.NewCustomerAddressRepository(s.baseRepo)
		filterCustomerAddress := map[string]interface{}{
			"id":          input.OrderQuoteAddress.CustomerAddressID,
			"customer_id": user.CustomerID,
		}
		ca, err := caRepo.FindFirstByParams(&dbc, filterCustomerAddress)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " customer_address " + err.Error()))
			return msg, err
		}
		if ca == nil {
			return message.QuoteCustomerAddressNotFound, errors.New(message.QuoteCustomerAddressNotFound.Message)
		}

		// delete existing
		arrQuoteAddress, _, _ := s.quoteAddressRepo.FindByParams(&dbc, map[string]interface{}{"quote_id": quote.ID}, 5, 1)
		if arrQuoteAddress != nil {
			for _, qa := range *arrQuoteAddress {
				_ = s.quoteAddressRepo.DeleteByID(&dbc, qa.ID)
			}
		}
		coordinate := ""
		if ca.Latitude != 0 && ca.Longitude != 0 {
			coordinate = fmt.Sprintf("(%v,%v)", ca.Latitude, ca.Longitude)
		}

		quoteAddress := entityquote.OrderQuoteAddress{
			QuoteID:           quote.ID,
			Title:             ca.Title,
			CustomerAddressID: ca.ID,
			AddressTypeID:     1,
			ReceiverName:      ca.ReceiverName,
			Street:            ca.Street,
			Province:          ca.Province,
			City:              ca.City,
			District:          ca.District,
			Subdistrict:       ca.Subdistrict,
			Coordinate:        coordinate,
			PhoneNumber:       ca.ReceiverPhone,
			Zipcode:           ca.Zipcode,
			CustomerNotes:     ca.Notes,
		}
		// save address
		_, err = s.quoteAddressRepo.Create(&dbc, &quoteAddress)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quote_address " + err.Error()))
			return msg, err
		}

		// trigger remove quote shipping
		_, _ = s.RemoveQuoteShipping(ctx, quote.ID)
	}

	// process quote payment
	if input.OrderQuotePayment != nil {
		// validation
		if err := input.OrderQuotePayment.Validate(); err != nil {
			return message.ValidationError, err
		}
		dbc := repository.DBContext{Context: ctx, DB: s.baseRepo.GetDB()}

		oqp, err := s.quotePaymentRepo.FindFirstByParams(&dbc, map[string]interface{}{"quote_id": quote.ID}, false)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quote_payment " + err.Error()))
			return msg, err
		}
		if oqp == nil {
			oqp = &entityquote.OrderQuotePayment{}
		}

		pmRepo := repository.NewPaymentMethodRepository(s.baseRepo)
		pm, err := pmRepo.FindFirstByParams(&dbc, map[string]interface{}{"id": input.OrderQuotePayment.PaymentMethodID})
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " quote_payment payment method " + err.Error()))
			return msg, err
		}
		if pm == nil {
			return message.QuotePaymentMethodNotFound, errors.New(message.QuotePaymentMethodNotFound.Message)
		}
		oqp.PaymentMethodID = pm.ID
		oqp.QuoteID = quote.ID
		_, err = s.quotePaymentRepo.Save(&dbc, oqp)

		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quote_payment " + err.Error()))
			return msg, err
		}
	}

	// process quote receipt
	if input.OrderQuoteReceipt != nil {
		jsonDataReceipt, err := sonic.Marshal(input.OrderQuoteReceipt)
		dbc := repository.DBContext{Context: ctx, DB: s.baseRepo.GetDB()}
		if err == nil {
			quote.DataReceipt = jsonDataReceipt
			err = s.quoteRepo.UpdateByQuoteCode(&dbc, quoteCode, *quote)
			if err != nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save data_receipt " + err.Error()))
				return msg, err
			}
		}
	}

	// recalculate quote
	_, _ = s.Recalculate(ctx, quote)

	return msg, nil
}

func (s QuoteReceiptServiceImpl) Recalculate(ctx context.Context, quote *entityquote.OrderQuote) (message.Message, error) {
	errMsgPrefix := "QUOTE-RECALCULATE"
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}

	var subTotal, totalShippingAmount, totalWeight float64
	var totalQty int

	quoteMerchants, _, err := s.quoteMerchantRepo.FindByParams(&dbc, map[string]interface{}{"quote_id": quote.ID}, false, 100, 1)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
		return message.ErrDB, errors.New(err.Error())
	}

	if quoteMerchants != nil {
		for _, quoteMerchant := range *quoteMerchants {
			subTotal += quoteMerchant.MerchantSubtotal
			totalQty += quoteMerchant.MerchantTotalQuantity
			totalWeight += quoteMerchant.MerchantTotalWeight

			quoteShipping, err := s.quoteShippingRepo.FindFirstByParams(&dbc, map[string]interface{}{"quote_merchant_id": quoteMerchant.ID})
			if err != nil {
				s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
				return message.ErrDB, errors.New(err.Error())
			}
			if quoteShipping != nil {
				totalShippingAmount += quoteShipping.ShippingCostActual
			} else {
				totalShippingAmount = 0
			}
		}
	}

	quote.Subtotal = subTotal
	quote.GrandTotal = quote.Subtotal + totalShippingAmount
	quote.ShippingAmount = totalShippingAmount
	quote.TotalQuantity = totalQty
	quote.Weight = totalWeight
	_, err = s.quoteRepo.Save(&dbc, quote)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
		return message.ErrDB, errors.New(err.Error())
	}

	return message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) ProcessQuoteShipping(ctx context.Context, input *request.QuoteReceiptShippingRq, quoteMerchant entityquote.OrderQuoteMerchant) (*map[uint64]map[string]float64, message.Message, error) {
	response := map[uint64]map[string]float64{}
	errMsgPrefix := "QUOTE-SHIPPING"
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}

	// validation
	if err := input.Validate(); err != nil {
		return nil, message.ValidationError, err
	}

	quoteShipping, err := s.quoteShippingRepo.FindFirstByParams(&dbc, map[string]interface{}{"quote_merchant_id": quoteMerchant.ID})
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
		return nil, message.ErrDB, errors.New(err.Error())
	}
	if quoteShipping == nil {
		quoteShipping = &entityquote.OrderQuoteShipping{}
	}
	helperKalcare := helper_kalcare.NewKalcareHelper(s.infra.Config.KalcareAPI, s.infra.Log)
	token := ctx.Value(jwt.JWTContextKey)
	header := map[string]string{
		"Authorization": "Bearer " + token.(string),
	}
	quoteShipping.QuoteMerchantID = quoteMerchant.ID
	quoteShipping.ShippingProviderDurationID = input.ShippingProviderDurationID

	if input.ShippingProviderID == nil {
		shippingRateDurations, err := helperKalcare.GetShippingRateMerchantDuration(ctx, quoteMerchant.QuoteID, quoteMerchant.MerchantID, quoteMerchant.OrderQuoteShipping.ShippingProviderDurationID, header)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get GetShippingRateMerchantDuration " + err.Error()))
			return nil, message.ErrThirdParty, err
		}
		if shippingRateDurations == nil || len(shippingRateDurations.Data.Records) == 0 {
			return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
		}
		quoteShipping.ShippingProviderID = shippingRateDurations.Data.Records[0].ShippingProviderID
	} else {
		quoteShipping.ShippingProviderID = *input.ShippingProviderID
	}
	// shipping not found
	if quoteShipping.ShippingProviderID == 0 {
		return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
	}

	// get shipping by provider_id
	shippingRateProviders, err := helperKalcare.GetShippingRateMerchantProvider(ctx, quoteMerchant.QuoteID, quoteMerchant.MerchantID, quoteShipping.ShippingProviderID, header)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get GetShippingRateMerchantProvider " + err.Error()))
		return nil, message.ErrThirdParty, err
	}
	if shippingRateProviders == nil || len(shippingRateProviders.Data.Records) == 0 || shippingRateProviders.Data.Records[0].PriceKg < 0 {
		return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
	}

	// assign quote shipping
	quoteShipping.ShippingRate = shippingRateProviders.Data.Records[0].PriceKg
	quoteShipping.ShippingCostActual = shippingRateProviders.Data.Records[0].ShippingCost
	quoteShipping.InsuranceFeeIncluded = shippingRateProviders.Data.Records[0].InsuranceFeeIncluded
	if shippingRateProviders.Data.Records[0].ShippingStartTime != "" {
		quoteShipping.ShippingStartTime = &shippingRateProviders.Data.Records[0].ShippingStartTime
	}
	if shippingRateProviders.Data.Records[0].ShippingEndTime != "" {
		quoteShipping.ShippingEndTime = &shippingRateProviders.Data.Records[0].ShippingEndTime
	}
	quoteShipping.InstanceDelivery = shippingRateProviders.Data.Records[0].InstanceDelivery
	_, err = s.quoteShippingRepo.Save(&dbc, quoteShipping)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quote_shipping " + err.Error()))
		return nil, message.ErrThirdParty, err
	}

	// add to quoteShippingMap
	if items, ok := response[quoteMerchant.ID]; ok {
		items["shipping_amount"] += quoteShipping.ShippingCostActual
	} else {
		mapItems := map[string]float64{
			"shipping_amount": quoteShipping.ShippingCostActual,
		}
		response[quoteMerchant.ID] = mapItems
	}

	return &response, message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) RecalculateQuoteShipping(ctx context.Context, quoteMerchant entityquote.OrderQuoteMerchant) (*map[uint64]map[string]float64, message.Message, error) {
	response := map[uint64]map[string]float64{}
	errMsgPrefix := "QUOTE-SHIPPING-RECALCULATE"
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}

	// if recalculate only get from DB not input
	quoteShipping, err := s.quoteShippingRepo.FindFirstByParams(&dbc, map[string]interface{}{"quote_merchant_id": quoteMerchant.ID})
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
		return nil, message.ErrDB, errors.New(err.Error())
	}
	if quoteShipping == nil {
		return nil, message.SuccessMsg, nil
	}

	helperKalcare := helper_kalcare.NewKalcareHelper(s.infra.Config.KalcareAPI, s.infra.Log)
	token := ctx.Value(jwt.JWTContextKey)
	header := map[string]string{
		"Authorization": "Bearer " + token.(string),
	}
	quoteShipping.QuoteMerchantID = quoteMerchant.ID

	if quoteShipping.ShippingProviderID == 0 {
		shippingRateDurations, err := helperKalcare.GetShippingRateMerchantDuration(ctx, quoteMerchant.QuoteID, quoteMerchant.MerchantID, quoteMerchant.OrderQuoteShipping.ShippingProviderDurationID, header)
		if err != nil {
			s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get GetShippingRateMerchantDuration " + err.Error()))
			return nil, message.ErrThirdParty, err
		}
		if shippingRateDurations == nil || len(shippingRateDurations.Data.Records) == 0 {
			return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
		}
		quoteShipping.ShippingProviderID = shippingRateDurations.Data.Records[0].ShippingProviderID
	}

	// shipping not found
	if quoteShipping.ShippingProviderID == 0 {
		return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
	}

	// get shipping by provider_id
	shippingRateProviders, err := helperKalcare.GetShippingRateMerchantProvider(ctx, quoteMerchant.QuoteID, quoteMerchant.MerchantID, quoteShipping.ShippingProviderID, header)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get GetShippingRateMerchantProvider " + err.Error()))
		return nil, message.ErrThirdParty, err
	}
	if shippingRateProviders == nil || len(shippingRateProviders.Data.Records) == 0 || shippingRateProviders.Data.Records[0].PriceKg < 0 {
		return nil, message.QuoteShippingNotFound, errors.New(message.QuoteShippingNotFound.Message)
	}

	// assign quote shipping
	quoteShipping.ShippingRate = shippingRateProviders.Data.Records[0].PriceKg
	quoteShipping.ShippingCostActual = shippingRateProviders.Data.Records[0].ShippingCost
	quoteShipping.InsuranceFeeIncluded = shippingRateProviders.Data.Records[0].InsuranceFeeIncluded
	if shippingRateProviders.Data.Records[0].ShippingStartTime != "" {
		quoteShipping.ShippingStartTime = &shippingRateProviders.Data.Records[0].ShippingStartTime
	}
	if shippingRateProviders.Data.Records[0].ShippingEndTime != "" {
		quoteShipping.ShippingEndTime = &shippingRateProviders.Data.Records[0].ShippingEndTime
	}
	quoteShipping.InstanceDelivery = shippingRateProviders.Data.Records[0].InstanceDelivery
	_, err = s.quoteShippingRepo.Save(&dbc, quoteShipping)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " save quote_shipping " + err.Error()))
		return nil, message.ErrThirdParty, err
	}

	// add to quoteShippingMap
	if items, ok := response[quoteMerchant.ID]; ok {
		items["shipping_amount"] += quoteShipping.ShippingCostActual
	} else {
		mapItems := map[string]float64{
			"shipping_amount": quoteShipping.ShippingCostActual,
		}
		response[quoteMerchant.ID] = mapItems
	}

	return &response, message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) RemoveQuoteShipping(ctx context.Context, quoteID uint64) (message.Message, error) {
	errMsgPrefix := "QUOTE-SHIPPING-REMOVE"
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}
	quoteMerchants, _, err := s.quoteMerchantRepo.FindByParams(&dbc, map[string]interface{}{"quote_id": quoteID}, false, 100, 1)
	if err != nil {
		s.infra.Log.WithContext(ctx).Error(errors.New(errMsgPrefix + " get quote_shipping " + err.Error()))
		return message.ErrDB, errors.New(err.Error())
	}

	if quoteMerchants != nil {
		for _, quoteMerchant := range *quoteMerchants {
			quoteShipping, _ := s.quoteShippingRepo.FindFirstByParams(&dbc, map[string]interface{}{"quote_merchant_id": quoteMerchant.ID})
			_ = s.quoteShippingRepo.DeleteByID(&dbc, quoteShipping.ID)
		}

	}
	return message.SuccessMsg, nil
}

func (s QuoteReceiptServiceImpl) isValidQty(ctx context.Context, input request.QuoteReceiptRq) bool {
	if len(input.OrderQuoteMerchants) > 0 {
		var quoteItems []request.QuoteReceiptItemsRq
		for _, oqm := range input.OrderQuoteMerchants {
			quoteItems = append(quoteItems, oqm.OrderQuoteItems...)
		}

		for _, qi := range quoteItems {
			if qi.Quantity <= 0 {
				return false
			}
		}
	}

	return true
}
