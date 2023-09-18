package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-svc/app"
	"marketplace-svc/app/model/request"
	"marketplace-svc/app/model/response"
	helper "marketplace-svc/helper/http"
	helperjwt "marketplace-svc/helper/jwt"
	"marketplace-svc/helper/message"
	"strings"

	helperlogger "gitlab.klik.doctor/platform/go-pkg/dapr/logger"

	"github.com/go-resty/resty/v2"
)

func UploadFile(ctx context.Context, infra app.Infra, http helper.IHttpRC, input request.UploadImageRequest) (response.UploadResponse, message.Message, error) {
	var result response.UploadResponse
	var messages message.Message

	correlationId := fmt.Sprint(helperlogger.GetTraceID(ctx))
	var log helperlogger.Logger = infra.Log.WithContext(ctx)
	var logTitle string = "[Media-Svc] - UploadFile"
	token, _ := helperjwt.GenerateJWTMedia(infra.Config.MediaSvcConfig)

	urlMediaService := infra.Config.MediaSvcConfig.URLUploadFile
	if strings.TrimSpace(input.CategoryUID) == "" {
		return result, message.ErrUploadMedia, errors.New("category_uid is required")
	}

	mapFormData := map[string]string{
		"name":               input.FileName,
		"media_category_uid": input.CategoryUID,
		"description":        input.Description,
	}
	headers := map[string]string{
		"Authorization":    "Bearer " + token,
		"X-Correlation-ID": correlationId,
	}
	fileReaders := []helper.FileReader{
		{
			Param:    "document",
			Filename: input.FileName,
			Reader:   bytes.NewReader(input.Image),
		},
	}

	// init wrap http resty
	resp, err := http.Execute(resty.MethodPost, urlMediaService, &headers, &mapFormData, &fileReaders, nil)

	if err != nil {
		log.Error(errors.New(fmt.Sprintf("%s, err: %s", logTitle, err.Error())))
		return result, message.ErrUploadMedia, err
	}
	log.Info(fmt.Sprintf("%s, respStatus: %v, respBody:%v", logTitle, resp.Response.StatusCode(), string(resp.Body)))
	rspStatusCode := resp.Response.StatusCode()

	var ResponseMedia response.ResponseHttpMedia
	_ = json.Unmarshal(resp.Body, &ResponseMedia)

	result.SourceType = input.SourceType
	if rspStatusCode != 200 {
		return result, message.ErrUploadMedia, errors.New(fmt.Sprint(ResponseMedia.Errors))
	}

	result.UID = &ResponseMedia.Data.Record.UID
	result.MediaPath = &ResponseMedia.Data.Record.DocumentPath
	mediaImage := infra.Config.MediaSvcConfig.S3BaseURL + "/" + ResponseMedia.Data.Record.DocumentPath
	result.MediaImage = &mediaImage
	messages = message.SuccessMsg

	return result, messages, nil
}
