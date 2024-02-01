package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuoteRepository struct {
	base.BaseRepository
}

type OrderQuoteRepository interface {
	Create(dbc *base.DBContext, oq *entity.OrderQuote) (*entity.OrderQuote, error)
	Save(dbc *base.DBContext, oq *entity.OrderQuote) (*entity.OrderQuote, error)
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuote, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuote, *modelbase.Pagination, error)
	UpdateByQuoteCode(dbc *base.DBContext, quoteCode string, data entity.OrderQuote) error
	UpdateMapByQuoteCode(dbc *base.DBContext, quoteCode string, data map[string]interface{}) error
}

func NewOrderQuoteRepository(br base.BaseRepository) OrderQuoteRepository {
	return &orderQuoteRepository{br}
}

func (r *orderQuoteRepository) Save(dbc *base.DBContext, oq *entity.OrderQuote) (*entity.OrderQuote, error) {
	err := dbc.DB.WithContext(dbc.Context).Save(oq).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return oq, nil
}

func (r *orderQuoteRepository) Create(dbc *base.DBContext, oq *entity.OrderQuote) (*entity.OrderQuote, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(oq).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return oq, nil
}

func (r *orderQuoteRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuote, error) {
	var orderQuote entity.OrderQuote
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuote.TableName())

	for key, v := range filter {
		if key == "quote_code" && v != "" {
			query = query.Where("LOWER(quote_code) = ?", v.(string))
		}
	}
	if isPreload {
		query = query.Preload("OrderQuoteAddress").
			Preload("OrderQuotePayment").
			Preload("OrderQuote").
			Preload("OrderQuote.OrderQuoteItem").
			Preload("OrderQuote.OrderQuoteShipping")
	}

	err := query.Omit("created_at,updated_at,converted_at,redeem,event,agent_id,data_source,has_cod,customer_data,data_source_value,qoute_type,total_point_bonus,total_point_discount,total_point_earned,total_point_spent,total_point_spent_conversion,subsidized_amount,scope,admin_fee,admin_fee_calculation,admin_fee_type,admin_fee_type_id").
		Find(&orderQuote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuote, nil
}

func (r *orderQuoteRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuote, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuote
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "quote_code" && v != "" {
			query = query.Where("LOWER(quote_code) = ?", v.(string))
		}
		if key == "quote_id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	if isPreload {
		query = query.Preload("OrderQuoteAddress").
			Preload("OrderQuotePayment").
			Preload("OrderQuote").
			Preload("OrderQuote.OrderQuoteItem").
			Preload("OrderQuote.OrderQuoteShipping")
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

func (r *orderQuoteRepository) UpdateByQuoteCode(dbc *base.DBContext, quoteCode string, data entity.OrderQuote) error {
	if quoteCode == "" {
		return errors.New("quote_code is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuote{}).
		Omit("id", "created_at").
		Where("quote_code = ?", quoteCode).
		Updates(data).
		Error
	return err
}

func (r *orderQuoteRepository) UpdateMapByQuoteCode(dbc *base.DBContext, quoteCode string, data map[string]interface{}) error {
	if quoteCode == "" {
		return errors.New("quote_code is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuote{}).
		Omit("id", "created_at").
		Where("quote_code = ?", quoteCode).
		Updates(data).
		Error
	return err
}
