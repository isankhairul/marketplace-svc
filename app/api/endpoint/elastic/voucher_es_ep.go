package endpointelastic

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"
)

type EsVoucherEndpoint struct {
	Search endpoint.Endpoint
}

func MakeEsVoucherEndpoints(s elasticservice.ElasticVoucherService) EsVoucherEndpoint {
	return EsVoucherEndpoint{
		Search: makeSearchVoucher(s),
	}
}

func makeSearchVoucher(s elasticservice.ElasticVoucherService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.VoucherRequest)
		v, _, _ := requestEsGroup.Do("SearchVoucher_"+request.ToString(), func() (interface{}, error) {
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

		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}
