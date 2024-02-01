package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/promotion"
	base "marketplace-svc/app/repository"
)

type promotionCatalogProductPriceRepository struct {
	base.BaseRepository
}

type PromotionCatalogProductPriceRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.PromotionCatalogProductPrice, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PromotionCatalogProductPrice, *modelbase.Pagination, error)
}

func NewPromotionCatalogProductPriceRepository(br base.BaseRepository) PromotionCatalogProductPriceRepository {
	return &promotionCatalogProductPriceRepository{br}
}

func (r *promotionCatalogProductPriceRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.PromotionCatalogProductPrice, error) {
	var promotionCatalogProductPrice entity.PromotionCatalogProductPrice
	query := dbc.DB.WithContext(dbc.Context).Table(promotionCatalogProductPrice.TableName())

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "product_id" && v != "" {
			query = query.Where("product_id = ?", v.(uint64))
		}
		if key == "customer_group_id" && v != "" {
			query = query.Where("customer_group_id = ?", v.(int))
		}
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(uint64))
		}
		if key == "store_id" && v != "" {
			query = query.Where("store_id = ?", v.(uint64))
		}
		if key == "latest_start_date" && v != "" {
			query = query.Where("latest_start_date <= ?", v.(string))
		}
		if key == "earliest_end_date" && v != "" {
			query = query.Where("earliest_end_date >= ?", v.(string))
		}
	}

	err := query.First(&promotionCatalogProductPrice).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &promotionCatalogProductPrice, nil
}

func (r *promotionCatalogProductPriceRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PromotionCatalogProductPrice, *modelbase.Pagination, error) {
	var promotionCatalogProductPrice []entity.PromotionCatalogProductPrice
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "product_id" && v != "" {
			query = query.Where("product_id = ?", v.(uint64))
		}
		if key == "customer_group_id" && v != "" {
			query = query.Where("customer_group_id = ?", v.(uint64))
		}
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(uint64))
		}
		if key == "store_id" && v != "" {
			query = query.Where("store_id = ?", v.(uint64))
		}
		if key == "latest_start_date" && v != "" {
			query = query.Where("latest_start_date <= ?", v.(string))
		}
		if key == "earliest_end_date" && v != "" {
			query = query.Where("earliest_end_date >= ?", v.(string))
		}
	}

	err := query.Scopes(r.Paginate(promotionCatalogProductPrice, &pagination, query)).
		Order("id DESC").
		Find(&promotionCatalogProductPrice).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(promotionCatalogProductPrice))

	return &promotionCatalogProductPrice, &pagination, nil
}
