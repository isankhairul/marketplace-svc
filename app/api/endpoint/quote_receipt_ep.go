package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/request"
	"marketplace-svc/app/service"
	"marketplace-svc/helper/message"
)

type QuoteReceiptEndpoint struct {
	Find     endpoint.Endpoint
	Save     endpoint.Endpoint
	Create   endpoint.Endpoint
	Validate endpoint.Endpoint
}

func MakeQuoteReceiptEndpoints(s service.QuoteReceiptService) QuoteReceiptEndpoint {
	return QuoteReceiptEndpoint{
		Find:     makeQuoteFind(s),
		Save:     makeQuoteSave(s),
		Create:   makeQuoteCreate(s),
		Validate: makeQuoteValidate(s),
	}
}

func makeQuoteFind(s service.QuoteReceiptService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		quoteCode := req.(string)
		result, msg, err := s.Find(ctx, quoteCode, false)
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, result, nil), nil
	}
}

func makeQuoteSave(s service.QuoteReceiptService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rq := req.(request.QuoteReceiptRq)
		msg, err := s.Save(ctx, rq.QuoteCode, rq)
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, nil, nil), nil
	}
}

func makeQuoteCreate(s service.QuoteReceiptService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rq := req.(request.QuoteReceiptRq)
		result, msg, err := s.Create(ctx, rq)
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, result, nil), nil
	}
}

func makeQuoteValidate(s service.QuoteReceiptService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		quoteCode := req.(string)
		msg, err := s.ItemsAvailability(ctx, quoteCode)
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, nil, nil), nil
	}
}
