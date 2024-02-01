package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/catalog"
	base "marketplace-svc/app/repository"
	"strings"
)

type productRepository struct {
	base.BaseRepository
}

type ProductRepository interface {
	Create(dbc *base.DBContext, oqi *entity.Product) (*entity.Product, error)
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.Product, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.Product, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.Product) error
}

func NewProductRepository(br base.BaseRepository) ProductRepository {
	return &productRepository{br}
}

func (r *productRepository) Create(dbc *base.DBContext, oqi *entity.Product) (*entity.Product, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(oqi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqi, nil
}

func (r *productRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.Product, error) {
	var orderQuote entity.Product
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuote.TableName())

	for key, v := range filter {
		if key == "sku" && v != "" {
			query = query.Where("lower(sku) = ?", strings.ToLower(v.(string)))
		}
		if key == "slug" && v != "" {
			query = query.Where("lower(slug) = ?", strings.ToLower(v.(string)))
		}
		if key == "id" && v != "" {
			query = query.Where("id in ?", v.(uint64))
		}
	}
	if isPreload {
		query = query.
			Preload("ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id,image,image_thumbnail,is_default,status").Where("is_default=1 and status=true")
			})
	}

	err := query.Omit("base_point,base_point_rupiah,created_at,updated_at,deleted_at").
		Order("id DESC").
		First(&orderQuote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuote, nil
}

func (r *productRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.Product, *modelbase.Pagination, error) {
	var product []entity.Product
	var pagination modelbase.Pagination

	query := dbc.DB.WithContext(dbc.Context)
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "sku" && v != "" {
			query = query.Where("lower(sku) = ?", strings.ToLower(v.(string)))
		}
		if key == "slug" && v != "" {
			query = query.Where("lower(slug) = ?", strings.ToLower(v.(string)))
		}
		if key == "id" && v != "" {
			query = query.Where("id in ?", v.(uint64))
		}
	}
	if isPreload {
		query = query.
			Preload("ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id,image,image_thumbnail,is_default,status").Where("is_default=1 and status=true")
			})
	}

	err := query.Omit("base_point,base_point_rupiah,created_at,updated_at,deleted_at").
		Scopes(r.Paginate(product, &pagination, query)).
		Order("id DESC").
		Find(&product).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(product))

	return &product, &pagination, nil
}

func (r *productRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.Product) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.Product{}).
		Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err

}
