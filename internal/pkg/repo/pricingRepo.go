package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type IPricingRepo interface {
}

var pricingRepo IPricingRepo

type MysqlPricingRepo struct {
	IRepository
}

func newPricingRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	pricingRepo = &MysqlPricingRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetPricingRepo() IPricingRepo {
	return pricingRepo
}
