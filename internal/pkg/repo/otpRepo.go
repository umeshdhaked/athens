package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type IOtpRepo interface {
}

var otpRepo IOtpRepo

type MysqlOtpRepo struct {
	IRepository
}

func newOtpRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	otpRepo = &MysqlOtpRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetOtpRepo() IOtpRepo {
	return otpRepo
}
