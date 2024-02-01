package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace-svc/app/model/entity"
)

type customerAddressRepository struct {
	BaseRepository
}

type CustomerAddressRepository interface {
	FindFirstByParams(dbc *DBContext, filter map[string]interface{}) (*entity.CustomerAddress, error)
}

func NewCustomerAddressRepository(br BaseRepository) CustomerAddressRepository {
	return &customerAddressRepository{br}
}

func (r *customerAddressRepository) FindFirstByParams(dbc *DBContext, filter map[string]interface{}) (*entity.CustomerAddress, error) {
	var customerAddress entity.CustomerAddress
	query := dbc.DB.WithContext(dbc.Context).Table(customerAddress.TableName())

	for key, v := range filter {
		if key == "customer_id" && v != "" {
			query = query.Where("customer_id = ?", v.(int64))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
	}

	err := query.First(&customerAddress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &customerAddress, nil
}
