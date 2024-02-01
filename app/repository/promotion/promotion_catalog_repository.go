package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/promotion"
	base "marketplace-svc/app/repository"
)

type promotionCatalogRepository struct {
	base.BaseRepository
}

type PromotionCatalogRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.PromotionCatalog, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PromotionCatalog, *modelbase.Pagination, error)
}

func NewPromotionCatalogRepository(br base.BaseRepository) PromotionCatalogRepository {
	return &promotionCatalogRepository{br}
}

func (r *promotionCatalogRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.PromotionCatalog, error) {
	var promotionCatalog entity.PromotionCatalog
	query := dbc.DB.WithContext(dbc.Context).Table(promotionCatalog.TableName())

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.First(&promotionCatalog).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &promotionCatalog, nil
}

func (r *promotionCatalogRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PromotionCatalog, *modelbase.Pagination, error) {
	var promotionCatalog []entity.PromotionCatalog
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.Scopes(r.Paginate(promotionCatalog, &pagination, query)).
		Order("id DESC").
		Find(&promotionCatalog).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(promotionCatalog))

	return &promotionCatalog, &pagination, nil
}
