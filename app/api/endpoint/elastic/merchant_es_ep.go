package endpointelastic

import (
	"context"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/global"
	"marketplace-svc/helper/message"

	"github.com/go-kit/kit/endpoint"
)

type EsMerchantEndpoint struct {
	Search                endpoint.Endpoint
	Detail                endpoint.Endpoint
	SearchByZipcode       endpoint.Endpoint
	SearchMerchantProduct endpoint.Endpoint
}

func MakeEsMerchantEndpoints(s elasticservice.ElasticMerchantService, cfg *config.Config) EsMerchantEndpoint {
	return EsMerchantEndpoint{
		Search:                makeSearchMerchant(s),
		Detail:                makeDetailMerchant(s),
		SearchByZipcode:       makeSearchMerchantByZipcode(s),
		SearchMerchantProduct: makeSearchMerchantProduct(s, cfg),
	}
}

func makeSearchMerchant(s elasticservice.ElasticMerchantService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.MerchantRequest)
		v, _, _ := requestEsGroup.Do("SearchMerchant_"+request.ToString(), func() (interface{}, error) {
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

		//result, page, msg, err := s.Search(ctx, req.(requestelastic.MerchantRequest))
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}

func makeDetailMerchant(s elasticservice.ElasticMerchantService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.MerchantDetailRequest)
		v, _, _ := requestEsGroup.Do("DetailMerchant_"+request.ToString(), func() (interface{}, error) {
			result, msg, err := s.Detail(ctx, request)
			response := map[string]interface{}{
				"result": result,
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
		return base.SetHttpResponse(ctx, msg, response["result"], nil), nil
	}
}

func makeSearchMerchantByZipcode(s elasticservice.ElasticMerchantService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.MerchantZipcodeRequest)
		v, _, _ := requestEsGroup.Do("SearchMerchantByZipcode_"+request.ToString(), func() (interface{}, error) {
			result, page, msg, err := s.SearchByZipcode(ctx, request)
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

func makeSearchMerchantProduct(s elasticservice.ElasticMerchantService, cfg *config.Config) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(requestelastic.MerchantProductRequest)

		// get JWT payload
		jwtInfo, err := global.GetJWTInfoFromContext(ctx, &cfg.Security.JwtConfig)
		if err != nil {
			msg := message.Message{
				Code:    message.ErrNoAuth.Code,
				Message: err.Error(),
			}
			return base.SetHttpResponse(ctx, msg, nil, nil), nil

		}
		request.JwtPayload = *jwtInfo

		v, _, _ := requestEsGroup.Do("SearchMerchantProduct_"+request.ToString(), func() (interface{}, error) {
			result, page, msg, err := s.SearchMerchantProduct(ctx, &cfg.KalcareAPI, request)
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

		//result, page, msg, err := s.Search(ctx, req.(requestelastic.MerchantRequest))
		if msg != message.SuccessMsg {
			return base.SetErrorResponse(ctx, msg, err), nil
		}
		pagination := response["page"].(base.Pagination)
		return base.SetHttpResponse(ctx, msg, response["result"], &pagination), nil
	}
}
