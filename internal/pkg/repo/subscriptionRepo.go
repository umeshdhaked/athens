package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"

	"github.com/gin-gonic/gin"
)

var subscriptionRepo *SubscriptionRepo

type SubscriptionRepo struct {
	Repository
}

func newSubscriptionRepo(client *dynamodb.Client) {
	subscriptionRepo = &SubscriptionRepo{Repository: Repository{dbClient: client}}
}

func GetSubscriptionRepo() *SubscriptionRepo {
	return subscriptionRepo
}

func (s *SubscriptionRepo) FetchAllSubscriptionByStatus(ctx *gin.Context, userId string, status string) ([]models.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:        aws.String(models.TableUserSubscription),
		IndexName:        aws.String(models.IndexTableUserSubscriptionIndexUserID),
		FilterExpression: aws.String("SubStatus= :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: status},
		},
		KeyConditions: map[string]types.Condition{
			models.ColumnUserId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: userId},
				},
			},
		},
	}

	var resp, err = s.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		subscriptions := []models.UserSubscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) FetchSubscriptionByTypeSubType(ctx *gin.Context, userId string, Type string, SubType string) ([]models.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:        aws.String(models.TableUserSubscription),
		IndexName:        aws.String(models.IndexTableUserSubscriptionIndexUserID),
		FilterExpression: aws.String("SubType = :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: SubType},
		},
		KeyConditions: map[string]types.Condition{
			models.ColumnUserId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: userId},
				},
			},
			models.ColumnType: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: Type},
				},
			},
		},
	}

	var resp, err = s.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		subscriptions := []models.UserSubscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) GetSubscriptionFromId(ctx *gin.Context, subId string) (*models.UserSubscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(models.TableUserSubscription),
		//IndexName: aws.String("id-index"),
		KeyConditions: map[string]types.Condition{
			models.ColumnId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: subId},
				},
			},
		},
	}

	var resp, err = s.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		subscriptions := []models.UserSubscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return &subscriptions[0], nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []models.UserSubscription) error {
	writeRequests := []types.WriteRequest{}
	for _, us := range userSubsDto {
		item, _ := attributevalue.MarshalMap(us)
		writeRequests = append(writeRequests,
			types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	batchWrite := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			models.TableUserSubscription: writeRequests,
		},
	}
	output, er := s.dbClient.BatchWriteItem(ctx, &batchWrite)
	log.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("unable to add subscription to user"))
	}
	return nil
}

func (s *SubscriptionRepo) CreateUserSubscription(ctx *gin.Context, userSubsDto *models.UserSubscription) error {
	item, _ := attributevalue.MarshalMap(userSubsDto)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableUserSubscription),
		Item:      item,
	}

	output, err := s.dbClient.PutItem(ctx, params)
	log.Println(output)
	return err
}
