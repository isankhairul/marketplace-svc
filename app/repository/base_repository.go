package repository

import (
	"context"
	"math"

	"gorm.io/gorm"

	"marketplace-svc/app/model/base"
	db "marketplace-svc/helper/database"
)

type baseRepository struct {
	db *gorm.DB
}

type BaseRepository interface {
	GetDB() *gorm.DB
	Transaction(fc func(tx *gorm.DB) error) error
	BeginTx() *gorm.DB
	Create(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}, uid string, params map[string]interface{}) error
	Delete(ctx context.Context, entity interface{}, uid string) error
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB
}

func NewBaseRepository(db *db.Database) BaseRepository {
	g, ok := (*db).Client().(*gorm.DB)
	if !ok {
		return &baseRepository{}
	}
	return &baseRepository{
		db: g,
	}
}

func (br *baseRepository) GetDB() *gorm.DB {
	return br.db
}

func (br *baseRepository) Create(ctx context.Context, entity interface{}) error {
	return br.GetDB().WithContext(ctx).Create(entity).Error
}

func (br *baseRepository) Update(ctx context.Context, entity interface{}, uid string, params map[string]interface{}) error {
	return br.GetDB().WithContext(ctx).Model(entity).Where("uid=?", uid).Updates(params).Error
}

func (br *baseRepository) Delete(ctx context.Context, entity interface{}, uid string) error {
	return br.GetDB().WithContext(ctx).Where("uid=?", uid).Delete(entity).Error
}

// Transaction help you to do transaction implicitly
// read more: https://gitlab.klik.doctor/platform/backend/boilerplate/-/blob/doc/db-transaction/docs/handle-db-transaction.md#short-way
func (br *baseRepository) Transaction(fc func(tx *gorm.DB) error) error {
	return br.db.Transaction(fc)
}

// BeginTx return a Tx instance which allow you to control transaction manually
// read more: https://gitlab.klik.doctor/platform/backend/boilerplate/-/blob/doc/db-transaction/docs/handle-db-transaction.md#manual-way
func (br *baseRepository) BeginTx() *gorm.DB {
	return br.db.Begin()
}

func (br *baseRepository) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Debug().Model(value).Select(1).Count(&totalRecords)
	//db.Debug().Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))

	var records int64
	records = (int64(pagination.Limit * pagination.Page)) / int64(pagination.Page)

	if pagination.Page == pagination.TotalPage {
		total := (pagination.TotalPage - 1) * pagination.Limit
		records = totalRecords - int64(total)
	}

	if records >= totalRecords {
		records = totalRecords
	}
	pagination.Records = records
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
