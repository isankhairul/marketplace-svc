package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/service"
	"marketplace-svc/helper/message"
)

type QuoteReceiptEndpoint struct {
	Find endpoint.Endpoint
}

func MakeQuoteReceiptEndpoints(s service.QuoteReceiptService) QuoteReceiptEndpoint {
	return QuoteReceiptEndpoint{
		Find: makeQuoteFind(s),
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
