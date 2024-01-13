package elasticservice

import (
	"context"
	"marketplace-svc/helper/message"
	"testing"

	"marketplace-svc/app/model/base"
	requestelastic "marketplace-svc/app/model/request/elastic"
	responseelastic "marketplace-svc/app/model/response/elastic"

	"github.com/stretchr/testify/assert"
)

type mockElasticProductService struct{}

func (m *mockElasticProductService) Search(ctx context.Context, input requestelastic.ProductRequest) ([]responseelastic.ProductResponse, base.Pagination, message.Message, error) {
	return nil, base.Pagination{}, message.SuccessMsg, nil
}

func TestSearchProductSuccess(t *testing.T) {
	service := &mockElasticProductService{}

	_, _, msg, err := service.Search(context.Background(), requestelastic.ProductRequest{})

	assert.NoError(t, err)
	assert.Equal(t, message.SuccessMsg, msg)
}
