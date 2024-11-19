/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 17:44:01
 */

package db

// type DBClient struct {
// 	config *config.DBConfig
// 	Conn   *gorm.DB
// }

// func NewDBClient(config *config.DBConfig) (DBCli *DBClient, err error) {

// 	var db *gorm.DB

// 	switch config.Dialect {
// 	case "mysql":
// 		db, err = NewMysqlGormDB(config)
// 	case "postgres":
// 		db, err = NewPgSQLGormDB(config)
// 	default:
// 		db, err = NewMysqlGormDB(config)
// 	}

// 	return &DBClient{config: config, Conn: db}, err
// }

// func NewMysqlGormDB(config *config.DBConfig) (db *gorm.DB, err error) {

// 	loc := config.Loc
// 	if loc == "" {
// 		loc = "Local"
// 	}

// 	URI := config.URI
// 	if URI == "" {
// 		URI = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
// 			config.Username,
// 			config.Password,
// 			config.Host,
// 			config.Database,
// 			config.Charset,
// 			config.ParseTime,
// 			config.Loc)
// 	}
// 	URI = url.QueryEscape(URI)

// 	logger := zapgorm2.New(zap.L())
// 	logger.SetAsDefault()

// 	mc := mysql.Config{
// 		DSN:                       URI,
// 		DefaultStringSize:         256,
// 		SkipInitializeWithVersion: false,
// 	}

// 	gc := &gorm.Config{
// 		Logger:                                   logger,
// 		DisableForeignKeyConstraintWhenMigrating: true,
// 	}

// 	if db, err = gorm.Open(mysql.New(mc), gc); err != nil {
// 		return nil, err
// 	}
// 	if sqlDB, err := db.DB(); err != nil {
// 		return nil, err
// 	} else {
// 		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
// 		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
// 	}

// 	return db, nil
// }

// func NewPgSQLGormDB(config *config.DBConfig) (db *gorm.DB, err error) {

// 	var targetSessionAttrs string
// 	if config.TargetSessionAttrs != "" {
// 		targetSessionAttrs = fmt.Sprintf("&target_session_attrs=%s", config.TargetSessionAttrs)
// 	}

// 	URI := config.URI
// 	if URI == "" {
// 		URI = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s&connect_timeout=%d%s",
// 			config.Username,
// 			config.Password,
// 			config.Host,
// 			config.Database,
// 			config.Sslmode,
// 			config.ConnectTimeout,
// 			targetSessionAttrs,
// 		)
// 	}
// 	URI = url.QueryEscape(URI)

// 	logger := zapgorm2.New(zap.L())
// 	logger.SetAsDefault()

// 	pgc := postgres.Config{
// 		DSN:                  URI,
// 		PreferSimpleProtocol: false,
// 	}

// 	gc := &gorm.Config{
// 		Logger:                                   logger,
// 		DisableForeignKeyConstraintWhenMigrating: true,
// 	}

// 	if db, err = gorm.Open(postgres.New(pgc), gc); err != nil {
// 		return nil, err
// 	}
// 	if sqlDB, err := db.DB(); err != nil {
// 		return nil, err
// 	} else {
// 		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
// 		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
// 	}
// 	return db, nil
// }
