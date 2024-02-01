package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/entity"
)

type paymentMethodRepository struct {
	BaseRepository
}

type PaymentMethodRepository interface {
	FindFirstByParams(dbc *DBContext, filter map[string]interface{}) (*entity.PaymentMethod, error)
	FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PaymentMethod, *base.Pagination, error)
	UpdateByID(dbc *DBContext, id uint64, data entity.PaymentMethod) error
	DeleteByID(dbc *DBContext, id uint64) error
	Create(dbc *DBContext, oqi *entity.PaymentMethod) (*entity.PaymentMethod, error)
	Save(dbc *DBContext, oqi *entity.PaymentMethod) (*entity.PaymentMethod, error)
}

func NewPaymentMethodRepository(br BaseRepository) PaymentMethodRepository {
	return &paymentMethodRepository{br}
}

func (r *paymentMethodRepository) Create(dbc *DBContext, oqa *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(oqa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqa, nil
}

func (r *paymentMethodRepository) Save(dbc *DBContext, oqa *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	err := dbc.DB.WithContext(dbc.Context).Save(oqa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqa, nil
}

func (r *paymentMethodRepository) FindFirstByParams(dbc *DBContext, filter map[string]interface{}) (*entity.PaymentMethod, error) {
	var paymentMethod entity.PaymentMethod
	query := dbc.DB.WithContext(dbc.Context).Table(paymentMethod.TableName())

	for key, v := range filter {
		if key == "code" && v != "" {
			query = query.Where("code = ?", v.(string))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "payment_method_type_id" && v != "" {
			query = query.Where("payment_method_type_id = ?", v.(uint64))
		}
	}

	err := query.First(&paymentMethod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &paymentMethod, nil
}

func (r *paymentMethodRepository) FindByParams(dbc *DBContext, filter map[string]interface{}, limit int, page int) (*[]entity.PaymentMethod, *base.Pagination, error) {
	var paymentMethod []entity.PaymentMethod
	var pagination base.Pagination

	query := dbc.DB
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "code" && v != "" {
			query = query.Where("code = ?", v.(string))
		}
		if key == "id" && v != "" {
			query = query.Where("id = ?", v.(uint64))
		}
		if key == "payment_method_type_id" && v != "" {
			query = query.Where("payment_method_type_id = ?", v.(uint64))
		}
	}

	err := query.Scopes(r.Paginate(paymentMethod, &pagination, query)).
		Order("id DESC").
		Find(&paymentMethod).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(paymentMethod))

	return &paymentMethod, &pagination, nil
}

func (r *paymentMethodRepository) UpdateByID(dbc *DBContext, id uint64, data entity.PaymentMethod) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.PaymentMethod{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}

func (r *paymentMethodRepository) DeleteByID(dbc *DBContext, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Where("id = ?", id).
		Delete(entity.PaymentMethod{}).
		Error
	return err
}
