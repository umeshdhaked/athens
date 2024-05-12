package repo

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
)

var smsAuditRepo *SmsAuditRepo

type SmsAuditRepo struct {
	Repository
}

func newSmsAuditRepo(client *dynamodb.Client) {
	smsAuditRepo = &SmsAuditRepo{Repository: Repository{dbClient: client}}
}

func GetSmsAuditRepo() *SmsAuditRepo {
	return smsAuditRepo
}

func (s *SmsAuditRepo) CreateSmsAudit(ctx *gin.Context, smsAudit *models.SmsAudit) error {
	item, err := attributevalue.MarshalMap(smsAudit)
	if err != nil {
		log.Printf("error marhsalling item: %v", err)
		return nil
	}

	// Insert item into the database
	err = s.CreateItem(ctx, models.TableSmsAudit, item)
	if err != nil {
		log.Printf("error inserting item: %v", err)
		return nil
	}

	return err
}

func (s *SmsAuditRepo) BulkCreate(ctx *gin.Context, entities []models.SmsAudit) error {
	// Prepare the list of PutRequests for batch write operation
	var putRequests []types.WriteRequest
	for _, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			return fmt.Errorf("error marshaling item: %v", err)
		}

		putRequests = append(putRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	// Split the putRequests into batches of 25 (maximum batch size allowed)
	const batchSize = 25
	numBatches := (len(putRequests) + batchSize - 1) / batchSize

	// Perform batch write operation for each batch
	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(putRequests) {
			end = len(putRequests)
		}

		batchWriteRequests := putRequests[start:end]

		params := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				models.TableSmsAudit: batchWriteRequests,
			},
		}

		_, err := s.dbClient.BatchWriteItem(ctx, params)
		if err != nil {
			return fmt.Errorf("error performing batch write operation: %v", err)
		}
	}

	return nil
}
