package repository

import (
	"errors"
	"gorm.io/gorm"
	entity "marketplace-svc/app/model/entity/catalog"
	base "marketplace-svc/app/repository"
)

type productImageRepository struct {
	base.BaseRepository
}

type ProductImageRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.ProductImage, error)
}

func NewProductImageRepository(br base.BaseRepository) ProductImageRepository {
	return &productImageRepository{br}
}

func (r *productImageRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.ProductImage, error) {
	var pi entity.ProductImage
	query := dbc.DB.WithContext(dbc.Context).Table(pi.TableName())

	for key, v := range filter {
		if key == "product_id" && v != "" {
			query = query.Where("product_id = ?", v.(uint64))
		}
		if key == "status" && v != "" {
			query = query.Where("status = ?", v.(bool))
		}
		if key == "is_default" && v != "" {
			query = query.Where("is_default = ?", v.(int))
		}
	}
	err := query.First(&pi).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &pi, nil
}
