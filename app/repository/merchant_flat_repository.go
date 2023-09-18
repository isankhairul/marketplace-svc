package repository

import (
	"errors"
	"marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/merchant"
	db "marketplace-svc/helper/database"
	"strings"

	"gorm.io/gorm"
)

type merchantFlatRepository struct {
	BaseRepository
}

type MerchantFlatRepository interface {
	FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) ([]entity.MerchantFlat, *base.Pagination, error)
}

func NewMerchantFlatRepository(entityDb *db.Database) MerchantFlatRepository {
	return &merchantFlatRepository{NewBaseRepository(entityDb)}
}

func (r *merchantFlatRepository) FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) ([]entity.MerchantFlat, *base.Pagination, error) {
	var merchantFlats []entity.MerchantFlat
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
		Scopes(r.Paginate(merchantFlats, &pagination, query)).
		Find(&merchantFlats).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(merchantFlats))

	return merchantFlats, &pagination, nil
}
