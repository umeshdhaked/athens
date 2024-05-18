package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type IPendingJobsRepo interface {
}

var pendingJobsRepo IPendingJobsRepo

type MysqlPendingJobsRepo struct {
	IRepository
}

func newPendingJobsRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	pendingJobsRepo = &MysqlPendingJobsRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetPendingJobsRepo() IPendingJobsRepo {
	return pendingJobsRepo
}
