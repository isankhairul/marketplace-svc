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
	UpdateByQuoteCode(dbc *base.DBContext, quoteCode string, data entity.OrderQuote) error
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
			Preload("OrderQuotePayment.PaymentMethod", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "code", "payment_method_type_id", "name", "logo")
			}).
			Preload("OrderQuotePayment.PaymentMethod.PaymentMethodType", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "payment_method_type_code", "name")
			}).
			Preload("OrderQuoteMerchant").
			Preload("OrderQuoteMerchant.OrderQuoteItem").
			Preload("OrderQuoteMerchant.OrderQuoteItem.ProductFlat", func(db *gorm.DB) *gorm.DB {
				return db.Select("sku", "is_active", "status", "slug", "name")
			}).
			Preload("OrderQuoteMerchant.OrderQuoteItem.Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("OrderQuoteMerchant.OrderQuoteItem.Product.ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			}).
			Preload("OrderQuoteMerchant.OrderQuoteShipping").
			Preload("OrderQuoteMerchant.OrderQuoteShipping.ShippingProvider", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "shipping_provider_type_id", "shipping_provider_duration_id", "code", "name")
			}).
			Preload("OrderQuoteMerchant.OrderQuoteShipping.ShippingProvider.ShippingProviderDuration", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "duration", "duration_label")
			}).
			Preload("OrderQuoteMerchant.Merchant", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "merchant_name", "merchant_code", "image", "province_id")
			})
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

func (r *orderQuoteRepository) UpdateByQuoteCode(dbc *base.DBContext, quoteCode string, data entity.OrderQuote) error {
	if data.QuoteCode == "" {
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
