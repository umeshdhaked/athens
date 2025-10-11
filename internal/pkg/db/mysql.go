package db

import (
	"time"

	"github.com/umeshdhaked/athens/internal/config"
	"github.com/umeshdhaked/athens/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlDbInit() *gorm.DB {

	db, err := gorm.Open(getGormDialect(), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.GetLogger().WithField("error", err).Panic("mysql db connection failed")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.GetLogger().WithField("error", err).Panic("mysql db connection failed")
	}

	/*
		gorm v2 by default don't allow global delete without where clause.
		TODO: need to test once
	*/

	sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifeTime * time.Second))
	sqlDB.SetMaxOpenConns(MaxOpenConnections)
	sqlDB.SetMaxIdleConns(MaxIdleConnections)

	addCallbacks(db)

	if txDB := db.Exec("SELECT 1"); txDB.Error != nil {
		logger.GetLogger().WithField("error", err).Panic("ping failed")
	}

	logger.GetLogger().Info("mysql db connection successful")

	return db
}

func addCallbacks(gormDB *gorm.DB) {

	gormDB.Callback().Create().Before("gorm:create").Register("setCreatedUpdatedAt", func(db *gorm.DB) {
		if field := db.Statement.Schema.LookUpField("CreatedAt"); field != nil {
			db.Statement.SetColumn("CreatedAt", time.Now().Unix())
		}
		if field := db.Statement.Schema.LookUpField("UpdatedAt"); field != nil {
			db.Statement.SetColumn("UpdatedAt", time.Now().Unix())
		}
	})
}

func getGormDialect() gorm.Dialector {
	return mysql.New(mysql.Config{
		DSN:                       config.GetConfig().Db.URL(), // data source name
		DefaultStringSize:         256,                         // default size for string fields
		DisableDatetimePrecision:  true,                        // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                        // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                        // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                       // auto configure based on currently MySQL version
	})
}
