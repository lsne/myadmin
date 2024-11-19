/*
 * @Author: Liu Sainan
 * @Date: 2023-12-09 17:14:14
 */

package db

import (
	"fmt"
	"myadmin/internal/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type MysqlClient struct {
	config *config.DBConfig
	conn   *gorm.DB
}

func NewMysqlClient(config *config.DBConfig) (mysqlCli *MysqlClient, err error) {

	var db *gorm.DB

	loc := config.Loc
	if loc == "" {
		loc = "Local"
	}

	URI := config.URI
	if URI == "" {
		URI = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
			config.Username,
			config.Password,
			config.Host,
			config.Database,
			config.Charset,
			config.ParseTime,
			config.Loc)
	}
	// QueryEscape 之后会报错 URI 格式不正确
	// URI = url.QueryEscape(URI)

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	mc := mysql.Config{
		DSN:                       URI,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}

	gc := &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if db, err = gorm.Open(mysql.New(mc), gc); err != nil {
		return &MysqlClient{config: config, conn: db}, err
	}
	if sqlDB, err := db.DB(); err != nil {
		return &MysqlClient{config: config, conn: db}, err
	} else {
		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	return &MysqlClient{config: config, conn: db}, nil
}

func (db *MysqlClient) Config() *config.DBConfig {
	return db.config
}

func (db *MysqlClient) Conn() *gorm.DB {
	return db.conn
}

func (db *MysqlClient) Raw(sql string, values ...interface{}) *gorm.DB {
	return db.conn.Raw(sql, values...)
}

func (db *MysqlClient) Exec(sql string, values ...interface{}) *gorm.DB {
	return db.conn.Exec(sql, values...)
}

// ShowTables 执行 show tables
func (db *MysqlClient) ShowTables() ([]string, error) {
	return db.conn.Migrator().GetTables()
}
