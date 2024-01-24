package registry

import (
	"marketplace-svc/app"
	base "marketplace-svc/app/repository"
	repository "marketplace-svc/app/repository/quote"
	"marketplace-svc/app/service"
)

func RegisterQuoteReceiptService(app *app.Infra) service.QuoteReceiptService {
	baseRepo := base.NewBaseRepository(app.DB)
	return service.NewQuoteReceiptService(
		*app,
		baseRepo,
		repository.NewOrderQuoteRepository(baseRepo),
	)
}
