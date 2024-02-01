package registry

import (
	"marketplace-svc/app"
	base "marketplace-svc/app/repository"
	repomerchant "marketplace-svc/app/repository/merchant"
	repository "marketplace-svc/app/repository/quote"
	"marketplace-svc/app/service"
)

func RegisterQuoteReceiptService(app *app.Infra) service.QuoteReceiptService {
	baseRepo := base.NewBaseRepository(app.DB)
	return service.NewQuoteReceiptService(
		*app,
		baseRepo,
		repository.NewOrderQuoteRepository(baseRepo),
		repository.NewOrderQuoteMerchantRepository(baseRepo),
		repository.NewOrderQuoteItemRepository(baseRepo),
		repository.NewOrderQuoteShippingRepository(baseRepo),
		repository.NewOrderQuoteAddressRepository(baseRepo),
		repository.NewOrderQuotePaymentRepository(baseRepo),
		repomerchant.NewMerchantRepository(baseRepo),
		repomerchant.NewMerchantProductRepository(baseRepo),
	)
}
