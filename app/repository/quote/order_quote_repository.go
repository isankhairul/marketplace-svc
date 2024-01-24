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
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuote, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuote, *modelbase.Pagination, error)
	UpdateByQuoteCode(dbc *base.DBContext, data *entity.OrderQuote) error
}

func NewOrderQuoteRepository(br base.BaseRepository) OrderQuoteRepository {
	return &orderQuoteRepository{br}
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
		query = query.Debug().
			Preload("OrderQuoteAddress").
			Preload("OrderQuotePayment").
			Preload("OrderQuotePayment.PaymentMethod").
			Preload("OrderQuotePayment.PaymentMethod.PaymentMethodType").
			Preload("OrderQuoteMerchant").
			Preload("OrderQuoteMerchant.OrderQuoteItem").
			Preload("OrderQuoteMerchant.OrderQuoteItem.ProductFlat").
			Preload("OrderQuoteMerchant.OrderQuoteItem.Product").
			Preload("OrderQuoteMerchant.OrderQuoteItem.Product.ProductImage", "is_default=1").
			Preload("OrderQuoteMerchant.OrderQuoteShipping").
			Preload("OrderQuoteMerchant.OrderQuoteShipping.ShippingProvider").
			Preload("OrderQuoteMerchant.OrderQuoteShipping.ShippingProvider.ShippingProviderDuration").
			Preload("OrderQuoteMerchant.Merchant")

	}

	err := query.
		Order("id DESC").
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
	}

	if isPreload {
		query = query.Preload("OrderQuoteAddress").
			Preload("OrderQuotePayment").
			Preload("OrderQuoteMerchant").
			Preload("OrderQuoteMerchant.OrderQuoteItem").
			Preload("OrderQuoteMerchant.OrderQuoteShipping")
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

func (r *orderQuoteRepository) UpdateByQuoteCode(dbc *base.DBContext, data *entity.OrderQuote) error {
	// check QuoteCode
	if data.QuoteCode == "" {
		return errors.New("quote_code is required")
	}
	return dbc.DB.WithContext(dbc.Context).Save(&data).Error
}
