package elasticservice

import (
	"context"
	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"
	"marketplace-svc/helper/config"
	"marketplace-svc/helper/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockElasticMerchantService struct{}

func (m *mockElasticMerchantService) SearchMerchantProduct(ctx context.Context, cfg *config.KalcareAPI, input requestelastic.MerchantProductRequest) ([]responseelastic.MerchantProductResponse, base.Pagination, message.Message, error) {
	return nil, base.Pagination{}, message.SuccessMsg, nil
}

func TestSearchMerchantProductSuccess(t *testing.T) {
	service := &mockElasticMerchantService{}

	_, _, msg, err := service.SearchMerchantProduct(context.Background(), &config.KalcareAPI{}, requestelastic.MerchantProductRequest{})

	assert.NoError(t, err)
	assert.Equal(t, message.SuccessMsg, msg)
}
