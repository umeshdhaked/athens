package repo

import (
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IUserRepo interface {
}

var userRepo *MysqlUserRepo

type MysqlUserRepo struct {
	IRepository
}

func newUserRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	userRepo = &MysqlUserRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetUserRepo() IUserRepo {
	return userRepo
}
