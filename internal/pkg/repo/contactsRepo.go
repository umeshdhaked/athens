package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IContactsRepo interface {
}

var contactsRepo IContactsRepo

type MysqlContactsRepo struct {
	IRepository
}

func newContactsRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	contactsRepo = &MysqlContactsRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetContactsRepo() IContactsRepo {
	return contactsRepo
}

// BulkCreateContacts inserts multiple contacts in bulk into the DynamoDB table
func (c *MysqlContactsRepo) BulkCreate(ctx *gin.Context, contacts []models.Contacts) error {
	// todo implement
	//// Prepare the list of PutRequests for batch write operation
	//var putRequests []types.WriteRequest
	//for _, contact := range contacts {
	//	item, err := attributevalue.MarshalMap(contact)
	//	if err != nil {
	//		return fmt.Errorf("error marshaling contact item: %v", err)
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
	//			models.TableContacts: batchWriteRequests,
	//		},
	//	}
	//
	//	_, err := c.dbClient.BatchWriteItem(ctx, params)
	//	if err != nil {
	//		return fmt.Errorf("error performing batch write operation: %v", err)
	//	}
	//}

	return nil
}
