package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuoteItemRepository struct {
	base.BaseRepository
}

type OrderQuoteItemRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteItem, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteItem, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteItem) error
}

func NewOrderQuoteItemRepository(br base.BaseRepository) OrderQuoteItemRepository {
	return &orderQuoteItemRepository{br}
}

func (r *orderQuoteItemRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteItem, error) {
	var orderQuote entity.OrderQuoteItem
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuote.TableName())

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
	}
	if isPreload {
		query = query.Debug().
			Preload("ProductFlat", func(db *gorm.DB) *gorm.DB {
				return db.Select("sku", "is_active", "status", "slug", "name")
			}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("Product.ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			})
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

func (r *orderQuoteItemRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteItem, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuoteItem
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
	}
	if isPreload {
		query = query.Debug().
			Preload("ProductFlat", func(db *gorm.DB) *gorm.DB {
				return db.Select("sku", "is_active", "status", "slug", "name")
			}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("Product.ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			})
		//Preload("ProductCategory", func(db *gorm.DB) *gorm.DB {
		//	return db.Where("store=1")
		//}).
		//Preload("ProductCategory.Category", func(db *gorm.DB) *gorm.DB {
		//	return db.Select("id,name,slug").Where("store=1 and status=1")
		//})
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

func (r *orderQuoteItemRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteItem) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuoteItem{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}
