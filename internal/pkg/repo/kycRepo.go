package repo

import (
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IKycRepo interface {
}

var kycRepo *MysqlKycRepo

type MysqlKycRepo struct {
	IRepository
}

func newKycRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	kycRepo = &MysqlKycRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetKycRepo() IUserRepo {
	return kycRepo
}
