package service

import (
	"context"
	"errors"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/api/middleware"
	entityquote "marketplace-svc/app/model/entity/quote"
	responsequote "marketplace-svc/app/model/response/quote"
	"marketplace-svc/app/repository"
	repoquote "marketplace-svc/app/repository/quote"
	"marketplace-svc/app/transform"
	helperconst "marketplace-svc/helper/const"
	"marketplace-svc/helper/message"
	"strconv"
)

type QuoteReceiptService interface {
	Find(ctx context.Context, quoteCode string, validate bool) (*responsequote.QuoteRs, message.Message, error)
	CheckQuote(ctx context.Context, quoteCode string, quote *entityquote.OrderQuote) (*entityquote.OrderQuote, message.Message, error)
}

type QuoteReceiptServiceImpl struct {
	infra       app.Infra
	baseRepo    repository.BaseRepository
	quoteRepo   repoquote.OrderQuoteRepository
	OrderTypeID int
}

func NewQuoteReceiptService(
	infra app.Infra,
	br repository.BaseRepository,
	quoteRepo repoquote.OrderQuoteRepository,
) QuoteReceiptService {
	return &QuoteReceiptServiceImpl{infra, br, quoteRepo, helperconst.ORDER_TYPE_RECEIPT}
}

func (s *QuoteReceiptServiceImpl) CheckQuote(ctx context.Context, quoteCode string, quote *entityquote.OrderQuote) (*entityquote.OrderQuote, message.Message, error) {
	deviceID := 1
	ctxDeviceID, err := strconv.ParseInt(fmt.Sprint(ctx.Value("device_id")), 10, 8)
	if ctxDeviceID != 0 && err == nil {
		deviceID = int(ctxDeviceID)
	}

	user, isLogged := middleware.IsAuthContext(ctx)
	if !isLogged {
		return quote, message.ErrNoAuth, errors.New(message.ErrNoAuth.Message)
	}
	dbc := repository.DBContext{Context: context.Background(), DB: s.baseRepo.GetDB()}
	if quote == nil {
		filter := map[string]interface{}{
			"quote_code":    quoteCode,
			"order_type_id": s.OrderTypeID,
			"customer_id":   user.CustomerID,
		}
		quoteRs, err := s.quoteRepo.FindFirstByParams(&dbc, filter, true)
		if err != nil {
			return nil, message.ErrNoData, err
		}
		quote = quoteRs
	}
	// update device_id
	quote.DeviceID = deviceID
	//err = s.quoteRepo.UpdateByQuoteCode(&dbc, quote)
	//if err != nil {
	//	return nil, message.ErrDB, err
	//}

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
