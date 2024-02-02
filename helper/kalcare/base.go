package helper_kalcare

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
	responsekalcare "marketplace-svc/app/model/response/kalcare"
	"marketplace-svc/helper/config"
	helper "marketplace-svc/helper/http"
)

type IKalcareHelper interface {
	// shipping
	GetShippingRateMerchantProvider(ctx context.Context, quoteID uint64, merchantID uint64, providerID uint64, header map[string]string) (*responsekalcare.ListShippingRateProvider, error)
	GetShippingRateMerchantDuration(ctx context.Context, quoteID uint64, merchantID uint64, durationID uint64, header map[string]string) (*responsekalcare.ListShippingRateDuration, error)
}

type KalcareHelper struct {
	Config config.KalcareAPI
	Log    logger.Logger
	HttpRC helper.IHttpRC
}

func NewKalcareHelper(c config.KalcareAPI, log logger.Logger) IKalcareHelper {
	httpRC := helper.NewHttpRC(resty.New(), log)
	return &KalcareHelper{Config: c, Log: log, HttpRC: httpRC}
}

func getErrorMessage(resMap map[string]interface{}) string {
	msg := "Something wrong, please try again"
	if rawMsg, ok := resMap["message"]; ok {
		msg = fmt.Sprint(rawMsg)
	} else {
		if errorMap, ok := resMap["errors"]; ok {
			errMsg := errorMap.(map[string]interface{})["message"]
			if errMsg != nil {
				msg = fmt.Sprint(errMsg)
			}
		}
	}
	return msg
}

func (s KalcareHelper) GetBaseEndpoint() string {
	return s.Config.Server
}
