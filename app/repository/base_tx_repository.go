package repository

import (
	"context"
	"gorm.io/gorm"
)

func NewDBContext(db *gorm.DB, ctx context.Context) *DBContext {
	return &DBContext{db, ctx}
}

type DBContext struct {
	DB      *gorm.DB
	Context context.Context
}

func (dbc *DBContext) TxBegin() {
	if dbc != nil && dbc.DB != nil {
		dbc.DB = dbc.DB.Begin()
	}
}

func (dbc *DBContext) TxCommit() {
	if dbc != nil && dbc.DB != nil {
		dbc.DB.Commit()
	}
}

func (dbc *DBContext) TxRollback() {
	if dbc != nil && dbc.DB != nil {
		dbc.DB.Rollback()
	}
}
