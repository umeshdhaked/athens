package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"

	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/gin-gonic/gin"
)

type SubscriptionRepo struct {
	client *dynamodb.Client
}

func NewSubscriptionRepo(client *dynamodb.Client) *SubscriptionRepo {
	return &SubscriptionRepo{client: client}
}

func (s *SubscriptionRepo) FetchAllSubscriptionForAUser(ctx *gin.Context, userId string) ([]dbo.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("user_subscriptions"),
		IndexName: aws.String(models.IndexTableUserSubscriptionIndexUserID),
		KeyConditions: map[string]types.Condition{
			"user_id": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: userId},
				},
			},
		},
	}

	var resp, err = s.client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		subscriptions := []dbo.UserSubscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) GetSubscriptionFromId(ctx *gin.Context, subId string) (*dbo.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("user_subscriptions"),
		IndexName: aws.String("id-index"),
		KeyConditions: map[string]types.Condition{
			"id": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: subId},
				},
			},
		},
	}

	var resp, err = s.client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		subscriptions := []dbo.UserSubscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return &subscriptions[0], nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []dbo.UserSubscription) error {
	writeRequests := []types.WriteRequest{}
	for _, us := range userSubsDto {
		item, _ := attributevalue.MarshalMap(us)
		writeRequests = append(writeRequests,
			types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	batchWrite := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"user_subscriptions": writeRequests,
		},
	}
	output, er := s.client.BatchWriteItem(ctx, &batchWrite)
	log.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("unable to add subscription to user"))
	}
	return nil
}

func (s *SubscriptionRepo) CreateUserSubscription(ctx *gin.Context, userSubsDto *dbo.UserSubscription) error {
	item, _ := attributevalue.MarshalMap(userSubsDto)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("user_subscriptions"),
		Item:      item,
	}

	output, err := s.client.PutItem(ctx, params)
	log.Println(output)
	return err
}
