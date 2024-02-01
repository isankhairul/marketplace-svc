package repository

import (
	"errors"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
)

type orderQuoteAddressRepository struct {
	base.BaseRepository
}

type OrderQuoteAddressRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.OrderQuoteAddress, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.OrderQuoteAddress, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteAddress) error
	DeleteByID(dbc *base.DBContext, id uint64) error
	Create(dbc *base.DBContext, oqi *entity.OrderQuoteAddress) (*entity.OrderQuoteAddress, error)
	Save(dbc *base.DBContext, oqi *entity.OrderQuoteAddress) (*entity.OrderQuoteAddress, error)
}

func NewOrderQuoteAddressRepository(br base.BaseRepository) OrderQuoteAddressRepository {
	return &orderQuoteAddressRepository{br}
}

func (r *orderQuoteAddressRepository) Create(dbc *base.DBContext, oqa *entity.OrderQuoteAddress) (*entity.OrderQuoteAddress, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(oqa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqa, nil
}

func (r *orderQuoteAddressRepository) Save(dbc *base.DBContext, oqa *entity.OrderQuoteAddress) (*entity.OrderQuoteAddress, error) {
	err := dbc.DB.WithContext(dbc.Context).Save(oqa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqa, nil
}

func (r *orderQuoteAddressRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}) (*entity.OrderQuoteAddress, error) {
	var orderQuoteAddress entity.OrderQuoteAddress
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuoteAddress.TableName())

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.Select("id,quote_id,customer_address_id,title,email,receiver_name,street,province_id,city_id,district_id,zipcode,phone_number,customer_notes,postcode_id " +
		" coordinate,province,city,district,subdistrict").
		First(&orderQuoteAddress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuoteAddress, nil
}

func (r *orderQuoteAddressRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.OrderQuoteAddress, *modelbase.Pagination, error) {
	var orderQuoteAddress []entity.OrderQuoteAddress
	var pagination modelbase.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.Scopes(r.Paginate(orderQuoteAddress, &pagination, query)).
		Order("id DESC").
		Find(&orderQuoteAddress).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(orderQuoteAddress))

	return &orderQuoteAddress, &pagination, nil
}

func (r *orderQuoteAddressRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteAddress) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuoteAddress{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}

func (r *orderQuoteAddressRepository) DeleteByID(dbc *base.DBContext, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Where("id = ?", id).
		Delete(entity.OrderQuoteAddress{}).
		Error
	return err
}
