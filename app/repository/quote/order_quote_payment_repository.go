package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuotePaymentRepository struct {
	base.BaseRepository
}

type OrderQuotePaymentRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuotePayment, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuotePayment, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuotePayment) error
}

func NewOrderQuotePaymentRepository(br base.BaseRepository) OrderQuotePaymentRepository {
	return &orderQuotePaymentRepository{br}
}

func (r *orderQuotePaymentRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuotePayment, error) {
	var orderQuotePayment entity.OrderQuotePayment
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuotePayment.TableName())

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	if isPreload {
		query = query.Debug().
			Preload("PaymentMethod", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "code", "payment_method_type_id", "name", "logo")
			}).
			Preload("PaymentMethod.PaymentMethodType", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "payment_method_type_code", "name")
			})
	}

	err := query.
		Order("id DESC").
		Find(&orderQuotePayment).Error

	if err != nil {
		fmt.Println("ERR", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuotePayment, nil
}

func (r *orderQuotePaymentRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuotePayment, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuotePayment
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	if isPreload {
		query = query.Debug().
			Preload("PaymentMethod", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "code", "payment_method_type_id", "name", "logo")
			}).
			Preload("PaymentMethod.PaymentMethodType", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "payment_method_type_code", "name")
			})
	}

	err := query.Scopes(r.Paginate(orderQuotes, &pagination, query)).
		Order("id DESC").
		Find(&orderQuotes).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(orderQuotes))

	return &orderQuotes, &pagination, nil
}

func (r *orderQuotePaymentRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuotePayment) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuotePayment{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}
