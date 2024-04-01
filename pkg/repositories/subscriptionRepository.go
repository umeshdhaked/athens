package repositories

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/gin-gonic/gin"
)

type SubscriptionRepo struct {
	svc *dynamodb.DynamoDB
}

func NewSubscriptionRepo(svc *dynamodb.DynamoDB) *SubscriptionRepo {
	return &SubscriptionRepo{svc: svc}
}

func (s *SubscriptionRepo) FetchAllSubscriptionForAUser(userId string) ([]dbo.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("user_subscriptions"),
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

	var resp, err = s.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	if *resp.Count > 0 {
		subscriptions := []dbo.UserSubscription{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			fmt.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []dbo.UserSubscription) error {
	writeRequests := []*dynamodb.WriteRequest{}
	for _, us := range userSubsDto {
		item, _ := dynamodbattribute.MarshalMap(us)
		writeRequests = append(writeRequests,
			&dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}})
	}

	batchWrite := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"user_subscriptions": writeRequests,
		},
	}
	output, er := s.svc.BatchWriteItemWithContext(ctx, &batchWrite)
	fmt.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("unable to add subscription to user"))
	}
	return nil
}

func (s *SubscriptionRepo) CreateUserSubscription(userSubsDto *dbo.UserSubscription) error {
	item, _ := dynamodbattribute.MarshalMap(userSubsDto)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("user_subscriptions"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	return req.Send()
}
