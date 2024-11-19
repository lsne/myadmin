/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 20:03:24
 */

package db

import (
	"myadmin/internal/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 01:48:46
 */

// 没必要用 sync.Map , 直接使用 map[string]
var (
	DBMap      = make(map[string]DBClient)
	RedisMap   = make(map[string]*RedisClient)
	MongoDBMap = make(map[string]*MongoDBClient)
	S3CephMap  = make(map[string]*S3Ceph)
)

type DBClient interface {
	Config() *config.DBConfig
	Conn() *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
}

func DB(names ...string) DBClient {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if _, ok := DBMap[name]; !ok {
		CreateDBClient(name, config.GlobalConfig.DB[name])
	}

	return DBMap[name]
}

func Redis(names ...string) *RedisClient {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if _, ok := RedisMap[name]; !ok {
		CreateRedisClient(name, config.GlobalConfig.Redis[name])
	}

	return RedisMap[name]
}

func MongoDB(names ...string) *MongoDBClient {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if _, ok := MongoDBMap[name]; !ok {
		CreateMongoDBClient(name, config.GlobalConfig.Mongo[name])
	}

	return MongoDBMap[name]
}

func S3(names ...string) *S3Ceph {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if _, ok := S3CephMap[name]; !ok {
		CreateS3Client(name, config.GlobalConfig.S3[name])
	}

	return S3CephMap[name]
}

func Setup() {
	for key, config := range config.GlobalConfig.DB {
		CreateDBClient(key, config)
	}
	for key, config := range config.GlobalConfig.Redis {
		CreateRedisClient(key, config)
	}
	for key, config := range config.GlobalConfig.Mongo {
		CreateMongoDBClient(key, config)
	}
	for key, config := range config.GlobalConfig.S3 {
		CreateS3Client(key, config)
	}
}

func CreateDBClient(name string, config *config.DBConfig) {

	var client DBClient
	var err error

	switch config.Dialect {
	case "mysql":
		client, err = NewMysqlClient(config)
	case "postgres":
		client, err = NewPgSQLClient(config)
	}

	if err != nil {
		zap.L().Error("connect mysql error", zap.String("error", err.Error()))
	}

	if client != nil {
		DBMap[name] = client
	}
}

func CreateMongoDBClient(name string, config *config.MongoConfig) {

	client, err := NewMongoDBClient(config)
	if err != nil {
		zap.L().Error("connect mongoclient error", zap.String("error", err.Error()))
	}
	if client != nil {
		MongoDBMap[name] = client
	}
}

func CreateRedisClient(name string, config *config.RedisConfig) {
	client, err := NewRedisClient(config)
	if err != nil {
		zap.L().Error("connect redis error", zap.String("error", err.Error()))
	}
	if client != nil {
		RedisMap[name] = client
	}
}

func CreateS3Client(name string, config *config.S3Config) {
	client, err := NewS3Ceph(config)
	if err != nil {
		zap.L().Error("connect s3Ceph error", zap.String("error", err.Error()))
	}
	if client != nil {
		S3CephMap[name] = client
	}
}
