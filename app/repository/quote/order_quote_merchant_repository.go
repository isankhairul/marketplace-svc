package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuoteMerchantRepository struct {
	base.BaseRepository
}

type OrderQuoteMerchantRepository interface {
	Create(dbc *base.DBContext, merchant *entity.OrderQuoteMerchant) (*entity.OrderQuoteMerchant, error)
	Save(dbc *base.DBContext, merchant *entity.OrderQuoteMerchant) (*entity.OrderQuoteMerchant, error)
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteMerchant, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteMerchant, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteMerchant) error
	DeleteByID(dbc *base.DBContext, id uint64) error
}

func NewOrderQuoteMerchantRepository(br base.BaseRepository) OrderQuoteMerchantRepository {
	return &orderQuoteMerchantRepository{br}
}

func (r *orderQuoteMerchantRepository) Save(dbc *base.DBContext, merchant *entity.OrderQuoteMerchant) (*entity.OrderQuoteMerchant, error) {
	err := dbc.DB.WithContext(dbc.Context).Save(merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return merchant, nil
}

func (r *orderQuoteMerchantRepository) Create(dbc *base.DBContext, merchant *entity.OrderQuoteMerchant) (*entity.OrderQuoteMerchant, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return merchant, nil
}

func (r *orderQuoteMerchantRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteMerchant, error) {
	var orderQuote entity.OrderQuoteMerchant
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuote.TableName())

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(uint64))
		}
	}
	if isPreload {
		query = query.
			//Preload("OrderQuoteItem").
			//Preload("OrderQuoteItem.ProductFlat", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("sku", "is_active", "status", "slug", "name")
			//}).
			//Preload("OrderQuoteItem.Product", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("id", "sku", "slug", "name")
			//}).
			//Preload("OrderQuoteItem.Product.ProductImage", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			//}).
			Preload("Merchant", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "merchant_name", "merchant_code", "image", "province_id")
			}).
			Preload("OrderQuoteShipping").
			Preload("OrderQuoteShipping.ShippingProvider", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "shipping_provider_type_id", "shipping_provider_duration_id", "code", "name")
			}).
			Preload("OrderQuoteShipping.ShippingProvider.ShippingProviderDuration", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "duration", "duration_label")
			})
	}

	err := query.Omit("created_at,updated_at,merchant_total_quantity,merchant_total_weight,merchant_total_point_earned,merchant_total_point_spent,event,redeem").
		First(&orderQuote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderQuote, nil
}

func (r *orderQuoteMerchantRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteMerchant, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuoteMerchant
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(uint64))
		}
	}

	if isPreload {
		query = query.
			//Preload("OrderQuoteItem").
			//Preload("OrderQuoteItem.ProductFlat", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("sku", "is_active", "status", "slug", "name")
			//}).
			//Preload("OrderQuoteItem.Product", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("id", "sku", "slug", "name")
			//}).
			//Preload("OrderQuoteItem.Product.ProductImage", func(db *gorm.DB) *gorm.DB {
			//	return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			//}).
			Preload("Merchant", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,merchant_name,merchant_code,image,province_id,city_id")
			}).
			Preload("Merchant.Province", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,name")
			}).
			Preload("Merchant.City", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,name")
			}).
			Preload("OrderQuoteShipping").
			Preload("OrderQuoteShipping.ShippingProvider", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "shipping_provider_type_id", "shipping_provider_duration_id", "code", "name", "logo")
			}).
			Preload("OrderQuoteShipping.ShippingProvider.ShippingProviderDuration", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "duration", "duration_label")
			})
	}

	err := query.Omit("created_at,updated_at,merchant_total_point_earned,merchant_total_point_spent,event,redeem").
		Scopes(r.Paginate(orderQuotes, &pagination, query)).
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

func (r *orderQuoteMerchantRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteMerchant) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuoteMerchant{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}

func (r *orderQuoteMerchantRepository) DeleteByID(dbc *base.DBContext, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Where("id = ?", id).
		Delete(entity.OrderQuoteMerchant{}).
		Error
	return err
}
