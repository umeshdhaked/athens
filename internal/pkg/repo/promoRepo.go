package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IPromotionRepo interface {
	MarkContacted(ctx *gin.Context, obj *models.PromoPhone) error
}

var promotionRepo IPromotionRepo

type MysqlPromotionRepo struct {
	IRepository
}

func newPromotionRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	promotionRepo = &MysqlPromotionRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetPromotionRepo() IPromotionRepo {
	return promotionRepo
}

func (p *MysqlPromotionRepo) MarkContacted(ctx *gin.Context, obj *models.PromoPhone) error {
	// todo implement
	//item, er := attributevalue.MarshalMap(obj)
	//if er != nil {
	//	return er
	//}
	//params := &dynamodb.PutItemInput{
	//	TableName: aws.String(models.TablePromoPhone),
	//	Item:      item,
	//}
	//
	//output, err := p.dbClient.PutItem(ctx, params)
	//log.Print(output)
	//return err
	return nil
}
