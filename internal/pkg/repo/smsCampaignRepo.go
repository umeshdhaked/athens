package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/gorm"
)

type ISmsCampaignRepo interface {
}

var smsCampaignRepo ISmsCampaignRepo

type MysqlSmsCampaignRepo struct {
	IRepository
}

func newSmsCampaignRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	smsCampaignRepo = &MysqlSmsCampaignRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetSmsCampaignRepo() ISmsCampaignRepo {
	return smsCampaignRepo
}
