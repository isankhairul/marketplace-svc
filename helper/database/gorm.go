package database

import (
	"database/sql"
	"fmt"
	"log"
	"marketplace-svc/helper/config"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type gormdb struct {
	db *gorm.DB
}

func NewGormConnectPostgres(cfg *config.DBConfig, opt *Option) (Database, error) {
	g := &gormdb{}

	sqlDB, err := sql.Open("pgx", g.postgresDsn(cfg))
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), g.options(cfg))

	if err != nil {
		return nil, err
	}

	if len(opt.Migrate) > 0 {
		e := db.AutoMigrate(opt.Migrate...)
		if e != nil {
			return nil, e
		}
	}

	//assign the db to struct
	g.db = db

	//connection Pool
	_, err = g.connPool(cfg)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (m *gormdb) Client() interface{} {
	return m.db
}

func (m *gormdb) connPool(cfg *config.DBConfig) (*sql.DB, error) {
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
	sqlDB.SetConnMaxIdleTime(time.Second * cfg.ConnectionMaxIdleTime)
	sqlDB.SetConnMaxLifetime(time.Second * cfg.ConnectionMaxLifeTime)

	return sqlDB, nil
}

func (m *gormdb) options(cfg *config.DBConfig) *gorm.Config {
	schemaName := ""
	if cfg.SchemaName != "" {
		schemaName = cfg.SchemaName + "."
	}

	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   schemaName,
		},
		Logger: m.logger(cfg.LogConfig),
	}
}

func (m *gormdb) postgresDsn(cfg *config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		"disable",
		cfg.SchemaName,
	)
}

func (m *gormdb) logger(cfg config.DBLogConfig) logger.Interface {
	var logLevel logger.LogLevel
	switch cfg.Level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: cfg.IgnoreNotFound,
		},
	)
}
