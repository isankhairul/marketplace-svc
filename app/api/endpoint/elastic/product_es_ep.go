package endpointelastic

import (
	"context"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"

	"github.com/go-kit/kit/endpoint"
)

type EsProductEndpoint struct {
	Search endpoint.Endpoint
}

func MakeEsProductEndpoints(s elasticservice.ElasticProductService) EsProductEndpoint {
	return EsProductEndpoint{
		Search: makeSearchProduct(s),
	}
}

func makeSearchProduct(s elasticservice.ElasticProductService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.ProductRequest)
		v, _, _ := requestEsGroup.Do("SearchProduct_"+request.ToString(), func() (interface{}, error) {
			result, page, msg, err := s.Search(ctx, request)
			response := map[string]interface{}{
				"result": result,
				"page":   page,
				"msg":    msg,
				"err":    err,
			}
			return response, nil
		})
		response := v.(map[string]interface{})
		msg := response["msg"].(message.Message)

		//result, page, msg, err := s.Search(ctx, request)
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}
