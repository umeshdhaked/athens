package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ISubscriptionRepo interface {
	BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []models.Subscription) error
}

var subscriptionRepo ISubscriptionRepo

type MysqlSubscriptionRepo struct {
	IRepository
}

func newSubscriptionRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	subscriptionRepo = &MysqlSubscriptionRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetSubscriptionRepo() ISubscriptionRepo {
	return subscriptionRepo
}

func (s *MysqlSubscriptionRepo) BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []models.Subscription) error {
	// todo implement
	//writeRequests := []types.WriteRequest{}
	//for _, us := range userSubsDto {
	//	item, _ := attributevalue.MarshalMap(us)
	//	writeRequests = append(writeRequests,
	//		types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	//}
	//
	//batchWrite := dynamodb.BatchWriteItemInput{
	//	RequestItems: map[string][]types.WriteRequest{
	//		models.TableSubscription: writeRequests,
	//	},
	//}
	//output, er := s.dbClient.BatchWriteItem(ctx, &batchWrite)
	//log.Println(output)
	//if er != nil {
	//	return errors.Join(er, errors.New("unable to add subscription to user"))
	//}
	return nil
}
