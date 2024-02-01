package repository

import (
	"errors"
	"gorm.io/gorm"
	entity "marketplace-svc/app/model/entity/merchant"
	base "marketplace-svc/app/repository"
)

type merchantRepository struct {
	base.BaseRepository
}

type MerchantRepository interface {
	Create(dbc *base.DBContext, merchant *entity.Merchant) (*entity.Merchant, error)
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.Merchant, error)
}

func NewMerchantRepository(br base.BaseRepository) MerchantRepository {
	return &merchantRepository{br}
}

func (r *merchantRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.Merchant, error) {
	var merchant entity.Merchant
	query := dbc.DB.WithContext(dbc.Context).Table(merchant.TableName())

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.
		First(&merchant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) Create(dbc *base.DBContext, merchant *entity.Merchant) (*entity.Merchant, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return merchant, nil
}
