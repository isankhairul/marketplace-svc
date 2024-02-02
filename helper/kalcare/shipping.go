package helper_kalcare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	responsekalcare "marketplace-svc/app/model/response/kalcare"
	"marketplace-svc/helper/message"
	"strings"
)

func (s KalcareHelper) GetShippingRateMerchantProvider(ctx context.Context, quoteID uint64, merchantID uint64, providerID uint64, header map[string]string) (*responsekalcare.ListShippingRateProvider, error) {
	response := responsekalcare.ListShippingRateProvider{}
	urlReplacer := strings.NewReplacer("{id}", fmt.Sprint(quoteID), "{merchant_id}", fmt.Sprint(merchantID), "{provider_id}", fmt.Sprint(providerID))
	urlShippingProvider := urlReplacer.Replace(s.Config.EndpointShippingRateProvider)
	url := s.GetBaseEndpoint() + urlShippingProvider

	rcResponse, err := s.HttpRC.Execute(ctx, resty.MethodGet, url, &header, nil, nil, nil)
	if err != nil {
		s.Log.WithContext(ctx).Error(err)
		return nil, err
	}

	err = json.Unmarshal(rcResponse.Body, &response)
	if err != nil {
		s.Log.WithContext(ctx).Error(err)
		return nil, err
	}

	if rcResponse.Response.StatusCode() < 200 || rcResponse.Response.StatusCode() >= 300 {
		msg := fmt.Sprintf("error call endpoint (%v) status=(%v) kalcare, err=(%v)", url, rcResponse.Response.StatusCode(), string(rcResponse.Body))
		s.Log.WithContext(ctx).Error(errors.New(msg))
		return nil, errors.New(message.ErrThirdParty.Message)
	}

	return &response, nil
}

func (s KalcareHelper) GetShippingRateMerchantDuration(ctx context.Context, quoteID uint64, merchantID uint64, durationID uint64, header map[string]string) (*responsekalcare.ListShippingRateDuration, error) {
	response := responsekalcare.ListShippingRateDuration{}
	urlReplacer := strings.NewReplacer("{id}", fmt.Sprint(quoteID), "{merchant_id}", fmt.Sprint(merchantID), "{duration_id}", fmt.Sprint(durationID))
	urlShippingProvider := urlReplacer.Replace(s.Config.EndpointShippingRateDuration)
	url := s.GetBaseEndpoint() + urlShippingProvider

	rcResponse, err := s.HttpRC.Execute(ctx, resty.MethodGet, url, &header, nil, nil, nil)
	if err != nil {
		s.Log.WithContext(ctx).Error(err)
		return nil, err
	}

	err = json.Unmarshal(rcResponse.Body, &response)
	if err != nil {
		s.Log.WithContext(ctx).Error(err)
		return nil, err
	}

	if rcResponse.Response.StatusCode() < 200 || rcResponse.Response.StatusCode() >= 300 {
		msg := fmt.Sprintf("error call endpoint (%v) status=(%v) kalcare, err=(%v)", url, rcResponse.Response.StatusCode(), string(rcResponse.Body))
		s.Log.WithContext(ctx).Error(errors.New(msg))
		return nil, errors.New(message.ErrThirdParty.Message)
	}

	return &response, nil
}
