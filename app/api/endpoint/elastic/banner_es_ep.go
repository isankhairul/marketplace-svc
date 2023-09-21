package endpointelastic

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"
)

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
		result, page, msg, err := s.Search(ctx, req.(requestelastic.BannerRequest))
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, result, &page), nil
	}
}
