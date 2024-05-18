package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ISmsTemplateRepo interface {
}

var smsTemplateRepo ISmsTemplateRepo

type MysqlSmsTemplateRepo struct {
	IRepository
}

func newSmsTemplateRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	smsTemplateRepo = &MysqlSmsTemplateRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetSmsTemplateRepo() ISmsTemplateRepo {
	return smsTemplateRepo
}
