package thumbor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"marketplace-svc/helper/config"
	"net/url"
	"strings"
)

func GetThumborUrl(baseUrl string, secret string, input string) string {
	keyForSign := []byte(secret)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	replacer := strings.NewReplacer("/", "_", "+", "-")
	signature := replacer.Replace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	result := fmt.Sprintf("%s/%s/%s", baseUrl, signature, input)
	return result
}

func GetNewThumborImagesOriginal(cfg config.Config, mediaPath string) string {
	newMediaPath, _ := url.Parse(mediaPath)
	var formatImage string = cfg.ThumborConfig.FormatImage
	var sizeOriginal string = cfg.ThumborConfig.SizeArOriginal
	mediaPathThumbor := fmt.Sprintf("%s/%s/%s", sizeOriginal, formatImage, newMediaPath)
	var baseUrl string = cfg.ThumborConfig.BaseURL
	var secret string = cfg.ThumborConfig.Secret
	result := GetThumborUrl(baseUrl, secret, mediaPathThumbor)

	return result
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
