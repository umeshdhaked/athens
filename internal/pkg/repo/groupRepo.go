package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type IGroupRepo interface {
}

var groupRepo IGroupRepo

type MysqlGroupRepo struct {
	IRepository
}

func newGroupRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	groupRepo = &MysqlGroupRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetGroupRepo() IGroupRepo {
	return groupRepo
}
