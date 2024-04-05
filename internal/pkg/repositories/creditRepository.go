package repositories

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
)

type CreditsRepo struct {
	svc *dynamodb.DynamoDB
}

func NewCreditsRepo(svc *dynamodb.DynamoDB) *CreditsRepo {
	return &CreditsRepo{svc: svc}
}

func (c *CreditsRepo) CreateUserCredit(credit *dbo.Credits) error {
	item, _ := dynamodbattribute.MarshalMap(credit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("credit"),
		Item:      item,
	}

	req, output := c.svc.PutItemRequest(params)
	fmt.Print(output)
	return req.Send()
}

func (c *CreditsRepo) FetchUserCredit(userId string) (*dbo.Credits, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("credit"),
		IndexName: aws.String("user_id-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userId),
					},
				},
			},
		},
	}
	var resp, err = c.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	if *resp.Count > 0 {
		credits := []dbo.Credits{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &credits); err != nil {
			log.Println(err)
			return nil, err
		}
		return &credits[0], nil
	}
	return nil, nil
}
