package endpointelastic

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/message"
)

type EsCategoryEndpoint struct {
	Search     endpoint.Endpoint
	SearchTree endpoint.Endpoint
}

func MakeEsCategoryEndpoints(s elasticservice.ElasticCategoryService) EsCategoryEndpoint {
	return EsCategoryEndpoint{
		Search:     makeSearchCategory(s),
		SearchTree: makeSearchCategoryTree(s),
	}
}

func makeSearchCategory(s elasticservice.ElasticCategoryService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.CategoryRequest)
		v, _, _ := requestEsGroup.Do("SearchCategory_"+request.ToString(), func() (interface{}, error) {
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

		//result, page, msg, err := s.Search(ctx, req.(requestelastic.CategoryRequest))
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}

func makeSearchCategoryTree(s elasticservice.ElasticCategoryService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.CategoryTreeRequest)
		v, _, _ := requestEsGroup.Do("SearchCategoryTree_"+request.ToString(), func() (interface{}, error) {
			result, msg, err := s.SearchTree(ctx, request)
			response := map[string]interface{}{
				"result": result,
				"msg":    msg,
				"err":    err,
			}
			return response, nil
		})
		response := v.(map[string]interface{})
		msg := response["msg"].(message.Message)

		//result, page, msg, err := s.Search(ctx, req.(requestelastic.CategoryRequest))
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		return base.SetHttpResponse(ctx, msg, response["result"], nil), nil
	}
}
