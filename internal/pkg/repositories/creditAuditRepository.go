package repositories

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
)

type CreditsAuditRepo struct {
	svc *dynamodb.DynamoDB
}

func NewCreditsAuditRepo(svc *dynamodb.DynamoDB) *CreditsAuditRepo {
	return &CreditsAuditRepo{svc: svc}
}

func (c *CreditsAuditRepo) CreateUserCreditAudit(creditAudit *dbo.CreditAudits) error {
	item, _ := dynamodbattribute.MarshalMap(creditAudit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("credits_audit"),
		Item:      item,
	}

	req, output := c.svc.PutItemRequest(params)
	fmt.Print(output)
	return req.Send()
}
