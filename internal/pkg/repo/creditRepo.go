package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ICreditsRepo interface {
}

var creditsRepo ICreditsRepo

type MysqlCreditsRepo struct {
	IRepository
}

func newCreditsRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	creditsRepo = &MysqlCreditsRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetCreditsRepo() ICreditsRepo {
	return creditsRepo
}
