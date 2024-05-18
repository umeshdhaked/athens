package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ICreditsAuditRepo interface {
}

var creditsAuditRepo ICreditsAuditRepo

type MysqlCreditsAuditRepo struct {
	IRepository
}

func newCreditsAuditRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	creditsAuditRepo = &MysqlCreditsAuditRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetCreditsAuditRepo() ICreditsAuditRepo {
	return creditsAuditRepo
}
