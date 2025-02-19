package endpointelastic

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/sync/singleflight"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"
)

var requestEsGroup singleflight.Group

type EsBannerEndpoint struct {
	Search endpoint.Endpoint
}

func MakeEsBannerEndpoints(s elasticservice.ElasticBannerService) EsBannerEndpoint {
	return EsBannerEndpoint{
		Search: makeSearchBanner(s),
	}
}

func makeSearchBanner(s elasticservice.ElasticBannerService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.BannerRequest)
		v, _, _ := requestEsGroup.Do("SearchBanner_"+request.ToString(), func() (interface{}, error) {
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
		//pagination := page
		//return base.SetHttpResponse(ctx, msg, result, &pagination), nil
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}
