package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ICronProcessingRepo interface {
}

var cronProcessingRepo ICronProcessingRepo

type MysqlCronProcessingRepo struct {
	IRepository
}

func newCronProcessingRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	cronProcessingRepo = &MysqlCronProcessingRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetCronProcessingRepo() ICronProcessingRepo {
	return cronProcessingRepo
}
