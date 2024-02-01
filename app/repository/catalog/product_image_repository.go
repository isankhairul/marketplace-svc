package repository

import (
	"errors"
	"gorm.io/gorm"
	entity "marketplace-svc/app/model/entity/merchant"
	base "marketplace-svc/app/repository"
)

type merchantProductRepository struct {
	base.BaseRepository
}

type MerchantProductRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.MerchantProduct, error)
}

func NewMerchantProductRepository(br base.BaseRepository) MerchantProductRepository {
	return &merchantProductRepository{br}
}

func (r *merchantProductRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.MerchantProduct, error) {
	var merchantProduct entity.MerchantProduct
	query := dbc.DB.WithContext(dbc.Context).Table(merchantProduct.TableName())

	for key, v := range filter {
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(int64))
		}
		if key == "product_sku" && v != "" {
			query = query.Where("product_sku = ?", v.(string))
		}
	}
	if isPreload {
		query = query.Preload("MerchantProductPrice")
	}

	err := query.
		Order("id DESC").
		Find(&merchantProduct).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &merchantProduct, nil
}
