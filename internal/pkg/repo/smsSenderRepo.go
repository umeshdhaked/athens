package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ISmsSenderRepo interface {
}

var smsSenderRepo ISmsSenderRepo

type MysqlSmsSenderRepo struct {
	IRepository
}

func newSmsSenderRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	smsSenderRepo = &MysqlSmsSenderRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetSmsSenderRepo() ISmsSenderRepo {
	return smsSenderRepo
}
