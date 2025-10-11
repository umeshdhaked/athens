package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ISmsAuditRepo interface {
	BulkCreate(ctx *gin.Context, entities []models.SmsAudit) error
}

var smsAuditRepo ISmsAuditRepo

type MysqlSmsAuditRepo struct {
	IRepository
}

func newSmsAuditRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	smsAuditRepo = &MysqlSmsAuditRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetSmsAuditRepo() ISmsAuditRepo {
	return smsAuditRepo
}

func (s *MysqlSmsAuditRepo) BulkCreate(ctx *gin.Context, entities []models.SmsAudit) error {
	// todo implement
	//// Prepare the list of PutRequests for batch write operation
	//var putRequests []types.WriteRequest
	//for _, entity := range entities {
	//	item, err := attributevalue.MarshalMap(entity)
	//	if err != nil {
	//		return fmt.Errorf("error marshaling item: %v", err)
	//	}
	//
	//	putRequests = append(putRequests, types.WriteRequest{
	//		PutRequest: &types.PutRequest{
	//			Item: item,
	//		},
	//	})
	//}
	//
	//// Split the putRequests into batches of 25 (maximum batch size allowed)
	//const batchSize = 25
	//numBatches := (len(putRequests) + batchSize - 1) / batchSize
	//
	//// Perform batch write operation for each batch
	//for i := 0; i < numBatches; i++ {
	//	start := i * batchSize
	//	end := (i + 1) * batchSize
	//	if end > len(putRequests) {
	//		end = len(putRequests)
	//	}
	//
	//	batchWriteRequests := putRequests[start:end]
	//
	//	params := &dynamodb.BatchWriteItemInput{
	//		RequestItems: map[string][]types.WriteRequest{
	//			models.TableSmsAudit: batchWriteRequests,
	//		},
	//	}
	//
	//	_, err := s.dbClient.BatchWriteItem(ctx, params)
	//	if err != nil {
	//		return fmt.Errorf("error performing batch write operation: %v", err)
	//	}
	//}

	return nil
}
