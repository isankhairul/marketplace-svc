package database

import (
	"errors"
	"marketplace-svc/helper/config"
)

const (
	MysqlDBDriver    string = "mysql"
	PostgresDBDriver string = "postgres"
	SqliteDBDriver   string = "sqlite"
)

type Option struct {
	Migrate []interface{}
}

type Database interface {
	Client() interface{}
}

func NewDatabaseConnect(cfg *config.DBConfig, opt *Option) (Database, error) {
	switch cfg.Driver {
	case PostgresDBDriver:
		return NewGormConnectPostgres(cfg, opt)
	}
	return nil, errors.New("the database driver does not support")
}
