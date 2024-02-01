package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuoteShippingRepository struct {
	base.BaseRepository
}

type OrderQuoteShippingRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.OrderQuoteShipping, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.OrderQuoteShipping, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteShipping) error
}

func NewOrderQuoteShippingRepository(br base.BaseRepository) OrderQuoteShippingRepository {
	return &orderQuoteShippingRepository{br}
}

func (r *orderQuoteShippingRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.OrderQuoteShipping, error) {
	var orderQuoteShipping entity.OrderQuoteShipping
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuoteShipping.TableName())

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
	}

	err := query.
		First(&orderQuoteShipping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuoteShipping, nil
}

func (r *orderQuoteShippingRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.OrderQuoteShipping, *modelbase.Pagination, error) {
	var orderQuoteShipping []entity.OrderQuoteShipping
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
	}

	err := query.Scopes(r.Paginate(orderQuoteShipping, &pagination, query)).
		Order("id DESC").
		Find(&orderQuoteShipping).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(orderQuoteShipping))

	return &orderQuoteShipping, &pagination, nil
}

func (r *orderQuoteShippingRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteShipping) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuoteShipping{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}
