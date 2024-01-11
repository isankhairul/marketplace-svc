package communicates

import (
	"encoding/json"
	"io"
	"marketplace-svc/helper/config"
	"net/http"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

func GetCustomerInfo(log logger.Logger, cfg *config.KalcareAPI, token string) (*KalstoreCustomerInfoDetail, error) {
	client := &http.Client{}

	log.Info("[KalstoreAPIRequest] URL:" + cfg.Server + cfg.EndpointCustomerInfo + " REQ: token " + token)

	token = "Bearer " + token

	req, err := http.NewRequest(http.MethodGet, cfg.Server+cfg.EndpointCustomerInfo, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	if err != nil {
		log.Error(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}

	log.Info("[KalstoreAPIResponse] URL:" + cfg.Server + cfg.EndpointCustomerInfo + " RESP:" + string(resp))

	var result KalstoreCustomerInfoDetail
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error(err)
		return nil, err
	}

	return &result, nil
}
