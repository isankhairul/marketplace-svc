package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"marketplace-svc/app"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"gitlab.klik.doctor/platform/go-pkg/dapr/logger"
)

const hiddenChar string = "**********"
const hiddenRequestVariable = "password,phone"

func LoggingMiddleware(infra *app.Infra) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := infra.Log.WithContext(ctx)
			traceId := r.Header.Get("X-Correlation-ID")
			if traceId == "" {
				traceId = logger.GetTraceID(ctx)
				r.Header.Set("X-Correlation-ID", traceId)
			}

			reqBodyBytes := []byte{}
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				var err error
				reqBodyBytes, err = ioutil.ReadAll(r.Body)

				if err != nil {
					log.Error(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

				if len(reqBodyBytes) > 0 {
					// change request hidden
					var HiddenRequest map[string]any
					json.Unmarshal(reqBodyBytes, &HiddenRequest)

					events := strings.Split(hiddenRequestVariable, ",")
					for _, event := range events {
						if HiddenRequest[event] != "" {
							HiddenRequest[event] = hiddenChar
							modifiedReqBodyBytes, _ := json.Marshal(HiddenRequest)
							reqBodyBytes = modifiedReqBodyBytes
						}
					}

					minifiedBody := new(bytes.Buffer)
					err = json.Compact(minifiedBody, reqBodyBytes)
					if err != nil {
						log.Error(err)
					}

					reqBodyBytes = minifiedBody.Bytes()
				}
			}

			headers := make(map[string][]string)
			for key, values := range r.Header {
				headers[key] = values
			}
			jsonHeaders, _ := json.Marshal(headers)

			// get base url
			// detect url forwarder KrakenD
			var fullurl string
			if r.Header.Get("X-Forwarded-Host") != "" {
				scheme := r.Header.Get("X-Forwarded-Proto")
				host := r.Header.Get("X-Forwarded-Host")
				baseUrl := &url.URL{
					Scheme: scheme,
					Host:   host,
				}
				fullurl = baseUrl.String() + r.URL.String()
			} else {
				scheme := "http"
				host := r.Host
				if r.TLS != nil {
					scheme = "https"
				}
				baseUrl := &url.URL{
					Scheme: scheme,
					Host:   host,
				}
				fullurl = baseUrl.String() + r.URL.String()
			}

			// correlationId := logger.GetTraceIdentifier(ctx)
			correlationId := traceId
			ipAddress := r.RemoteAddr
			method := r.Method
			request := string(reqBodyBytes)
			header := string(jsonHeaders)

			log.Info("api-info: " + correlationId + " | input | IP: " + ipAddress + " | URL: " + fullurl + " | Method: " + method + " | Request: " + request + " | Header: " + header)

			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, r)

			respBody := recorder.Body.Bytes()
			code := strconv.Itoa(recorder.Code)
			response := string(respBody)

			log.Info("api-info: " + correlationId + " | output | Code: " + code + " | Response: " + response)

			for key, values := range recorder.Header() {
				for _, value := range values {
					w.Header().Add(key, value)
					w.WriteHeader(recorder.Code)
				}
			}
			_, _ = w.Write(respBody)
		})
	}
}
