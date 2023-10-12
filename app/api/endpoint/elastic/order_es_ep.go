package endpointelastic

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"
)

type EsOrderEndpoint struct {
	Search endpoint.Endpoint
}

func MakeEsOrderEndpoints(s elasticservice.ElasticOrderService) EsOrderEndpoint {
	return EsOrderEndpoint{
		Search: makeSearchOrder(s),
	}
}

func makeSearchOrder(s elasticservice.ElasticOrderService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.OrderRequest)
		result, pagination, msg, err := s.Search(ctx, request)

		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, result, &pagination), nil
	}
}
