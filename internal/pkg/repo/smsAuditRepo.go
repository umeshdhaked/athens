package repo

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
