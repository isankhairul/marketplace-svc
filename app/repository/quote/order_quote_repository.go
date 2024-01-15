package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
	"strings"
)

type orderQuoteRepository struct {
	base.BaseRepository
}

type OrderQuoteRepository interface {
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) ([]entity.OrderQuote, *modelbase.Pagination, error)
}

func NewOrderQuoteRepository(br base.BaseRepository) OrderQuoteRepository {
	return &orderQuoteRepository{br}
}

func (r *orderQuoteRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) ([]entity.OrderQuote, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuote
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "q" && v != "" {
			query = query.Where("LOWER(quote_code) LIKE ?", "%"+strings.ToLower(v.(string))+"%")
		}
	}
	err := query.Preload("OrderQuoteAddress").
		Preload("OrderQuotePayment").
		Preload("OrderQuoteMerchant").
		Preload("OrderQuoteMerchant.OrderQuoteItem").
		//Scopes(r.Paginate(orderQuotes, &pagination, query)).
		Order("id DESC").
		Limit(10).
		Find(&orderQuotes).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(orderQuotes))

	return orderQuotes, &pagination, nil
}
