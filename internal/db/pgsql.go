/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 15:16:33
 */

package db

import (
	"fmt"
	"myadmin/internal/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type PgSQLClient struct {
	config *config.DBConfig
	conn   *gorm.DB
}

func NewPgSQLClient(config *config.DBConfig) (PGSqlCli *PgSQLClient, err error) {

	var db *gorm.DB

	var targetSessionAttrs string
	if config.TargetSessionAttrs != "" {
		targetSessionAttrs = fmt.Sprintf("&target_session_attrs=%s", config.TargetSessionAttrs)
	}

	URI := config.URI
	if URI == "" {
		URI = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s&connect_timeout=%d%s",
			config.Username,
			config.Password,
			config.Host,
			config.Database,
			config.Sslmode,
			config.ConnectTimeout,
			targetSessionAttrs,
		)
	}
	// QueryEscape 之后会报错 URI 格式不正确
	// URI = url.QueryEscape(URI)

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	pgc := postgres.Config{
		DSN:                  URI,
		PreferSimpleProtocol: false,
	}

	gc := &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if db, err = gorm.Open(postgres.New(pgc), gc); err != nil {
		return &PgSQLClient{config: config, conn: db}, err
	}
	if sqlDB, err := db.DB(); err != nil {
		return &PgSQLClient{config: config, conn: db}, err
	} else {
		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	return &PgSQLClient{config: config, conn: db}, nil
}

func (db *PgSQLClient) Config() *config.DBConfig {
	return db.config
}

func (db *PgSQLClient) Conn() *gorm.DB {
	return db.conn
}

func (db *PgSQLClient) Raw(sql string, values ...interface{}) *gorm.DB {
	return db.conn.Raw(sql, values...)
}

func (db *PgSQLClient) Exec(sql string, values ...interface{}) *gorm.DB {
	return db.conn.Exec(sql, values...)
}
