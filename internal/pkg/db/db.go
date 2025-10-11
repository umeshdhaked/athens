package db

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/umeshdhaked/athens/internal/config"
	"gorm.io/gorm"
)

const (
	ConnMaxLifeTime    = 300
	MaxOpenConnections = 2
	MaxIdleConnections = 1
)

var (
	once sync.Once
	db   *Db
)

type Db struct {
	Mysql  *gorm.DB
	Dynamo *dynamodb.Client
}

func NewDb() {
	once.Do(func() {
		db = &Db{}

		if config.GetConfig().Db.Dynamo.Enabled {
			db.Dynamo = dynomoDbInit()
		}

		if config.GetConfig().Db.Mysql.Enabled {
			db.Mysql = mysqlDbInit()
		}
	})
}

func GetDb() *Db {
	return db
}
