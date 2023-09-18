package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/entity/booking"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/logger"
	"net/http"
)

type OTPApi interface {
	SendOTP(ctx context.Context, otp *booking.OTP) ([]byte, error)
}

type otpApiImp struct {
	microServiceCfg *config.MicroServiceConfig
	otpCfg          *config.OtpConfig
}

func NewOTPApi(microServiceCfg *config.MicroServiceConfig, otpCfg *config.OtpConfig) OTPApi {
	return &otpApiImp{microServiceCfg: microServiceCfg, otpCfg: otpCfg}
}

func (api *otpApiImp) SendOTP(ctx context.Context, otp *booking.OTP) ([]byte, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"job_body": map[string]interface{}{
			"msisdn":   otp.User,
			"uploadby": "KlikMedis",
			"channel":  2,
			"template": api.otpCfg.Template,
			"variables": map[string]string{
				"otp": otp.OTP,
			},
			"ip_address": ctx.Value(base.RequestIPAddressContextKey),
		},
		"job_delay": 0,
		"tube_name": "kd_send_sms",
	})
	apiReqParams := ApiRequestParams{
		Url: api.microServiceCfg.BaseUrl + api.microServiceCfg.SMSPathUrl,
		Headers: map[string]string{
			"X-Correlation-ID": logger.GetTraceIdentifier(ctx),
			"X-API-Auth":       api.microServiceCfg.ApiKey,
		},
		Method: http.MethodPost,
		Body:   reqBody,
	}
	status, data, err := performRequest(&apiReqParams)
	if err != nil {
		return nil, err
	}

	if status >= http.StatusInternalServerError {
		return nil, errors.New(fmt.Sprint(status, string(data[:])))
	}

	return data, nil
}
