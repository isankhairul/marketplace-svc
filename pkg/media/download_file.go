package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/model/response"
	helper "marketplace-svc/helper/http"
	helperjwt "marketplace-svc/helper/jwt"

	helperlogger "gitlab.klik.doctor/platform/go-pkg/dapr/logger"

	"github.com/go-resty/resty/v2"
)

func DownloadFile(ctx context.Context, infra app.Infra, http helper.IHttpRC, mediaPath string) (response.DownloadResponse, error) {
	var result response.DownloadResponse
	correlationId := fmt.Sprintf(helperlogger.GetTraceID(ctx))
	var log helperlogger.Logger = infra.Log.WithContext(ctx)
	logTitle := "[Media-Svc] - DownloadFile"
	token, _ := helperjwt.GenerateJWTMedia(infra.Config.MediaSvcConfig)
	urlMediaService := infra.Config.MediaSvcConfig.URLDownloadFile

	mapFormData := map[string]string{
		"path": mediaPath,
	}
	byteFormData, err := json.Marshal(mapFormData)
	if err != nil {
		// handle error
	}
	headers := map[string]string{
		"Authorization":    "Bearer " + token,
		"X-Correlation-ID": correlationId,
	}

	resp, err := http.Execute(resty.MethodPost, urlMediaService, &headers, nil, nil, &byteFormData)
	log.Info(fmt.Sprintf("%s, respStatus: %v, respBody:%v", logTitle, resp.Response.StatusCode(), string(resp.Body)))
	if err != nil {
		log.Error(errors.New(fmt.Sprintf("%s, err : %s", logTitle, err.Error())))
		return result, err
	}

	var ResponseDownloadMedia response.HttpResponseMedia
	err = json.Unmarshal(resp.Body, &ResponseDownloadMedia)
	if err != nil {
		log.Error(errors.New(fmt.Sprintf("%s, err : %s", logTitle, err.Error())))
		return result, err
	}

	result.MediaPath = mediaPath
	result.Url = ResponseDownloadMedia.Data.Record.Url
	return result, nil
}
