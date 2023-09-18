package repository

import (
	"errors"
	"marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/catalog"
	db "marketplace-svc/helper/database"
	"strings"

	"gorm.io/gorm"
)

type productFlatRepository struct {
	BaseRepository
}

type ProductFlatRepository interface {
	FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) ([]entity.ProductFlat, *base.Pagination, error)
}

func NewProductFlatRepository(entityDb *db.Database) ProductFlatRepository {
	return &productFlatRepository{NewBaseRepository(entityDb)}
}

func (r *productFlatRepository) FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) ([]entity.ProductFlat, *base.Pagination, error) {
	var productFlats []entity.ProductFlat
	var pagination base.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "q" && v != "" {
			query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(v.(string))+"%")
		}
	}
	err := query.
		Order("name ASC").
		Scopes(r.Paginate(productFlats, &pagination, query)).
		Find(&productFlats).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(productFlats))

	return productFlats, &pagination, nil
}
