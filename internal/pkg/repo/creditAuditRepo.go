package repo

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
)

type CreditsAuditRepo struct {
	client *dynamodb.Client
}

func NewCreditsAuditRepo(client *dynamodb.Client) *CreditsAuditRepo {
	return &CreditsAuditRepo{client: client}
}

func (c *CreditsAuditRepo) CreateUserCreditAudit(ctx *gin.Context, creditAudit *models.CreditAudits) error {
	item, _ := attributevalue.MarshalMap(creditAudit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableCreditAudits),
		Item:      item,
	}
	output, er := c.client.PutItem(ctx, params)
	log.Print(output)
	return er
}
