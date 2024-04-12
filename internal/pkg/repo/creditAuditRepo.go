package repo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"log"
)

var creditsAuditRepo *CreditsAuditRepo

type CreditsAuditRepo struct {
	Repository
}

func newCreditsAuditRepo(client *dynamodb.Client) {
	creditsAuditRepo = &CreditsAuditRepo{Repository: Repository{dbClient: client}}
}

func GetCreditsAuditRepo() *CreditsAuditRepo {
	return creditsAuditRepo
}

func (c *CreditsAuditRepo) CreateUserCreditAudit(ctx *gin.Context, creditAudit *models.CreditAudits) error {
	item, _ := attributevalue.MarshalMap(creditAudit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableCreditAudits),
		Item:      item,
	}
	output, er := c.dbClient.PutItem(ctx, params)
	log.Print(output)
	return er
}
