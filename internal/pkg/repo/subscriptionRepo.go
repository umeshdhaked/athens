package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/pkg/dtos"

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

func (s *SubscriptionRepo) FetchByUserIDAndConditions(ctx *gin.Context, userID string, conditions map[string]interface{}) ([]models.Subscription, error) {
	queryInput := dtos.DbQueryInputConditions{
		Index: models.IndexTableSubscriptionIndexUserID,
		PKey: map[string]interface{}{
			models.ColumnSubscriptionsUserId: userID,
		},
	}

	queryInput.NonPKey = make(map[string]interface{})

	for k, v := range conditions {
		queryInput.NonPKey[k] = v
	}

	subscriptionItems, err := s.QueryItems(ctx, models.TableSubscription, queryInput)
	if err != nil {
		log.Printf("error fetching pending jobs item: %v\n", err)
		return nil, err
	}

	var subscriptionEntity []models.Subscription
	if err := attributevalue.UnmarshalListOfMaps(subscriptionItems, &subscriptionEntity); err != nil {
		log.Println(err)
		return nil, err
	}

	return subscriptionEntity, nil
}

func (s *SubscriptionRepo) FetchAllSubscriptionByStatus(ctx *gin.Context, userId string, status string) ([]models.Subscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:        aws.String(models.TableSubscription),
		IndexName:        aws.String(models.IndexTableSubscriptionIndexUserID),
		FilterExpression: aws.String("SubStatus= :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: status},
		},
		KeyConditions: map[string]types.Condition{
			models.ColumnSubscriptionsUserId: {
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
		subscriptions := []models.Subscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) FetchSubscriptionByTypeSubType(ctx *gin.Context, userId string, Type string, SubType string) ([]models.Subscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:        aws.String(models.TableSubscription),
		IndexName:        aws.String(models.IndexTableSubscriptionIndexUserID),
		FilterExpression: aws.String("SubType = :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: SubType},
		},
		KeyConditions: map[string]types.Condition{
			models.ColumnSubscriptionsUserId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: userId},
				},
			},
			models.ColumnSubscriptionsType: {
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
		subscriptions := []models.Subscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return subscriptions, nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) GetSubscriptionFromId(ctx *gin.Context, subId string) (*models.Subscription, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(models.TableSubscription),
		//IndexName: aws.String("id-index"),
		KeyConditions: map[string]types.Condition{
			models.ColumnSubscriptionsID: {
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
		subscriptions := []models.Subscription{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &subscriptions); err != nil {
			log.Println(err)
			return nil, err
		}
		return &subscriptions[0], nil
	}
	return nil, nil
}

func (s *SubscriptionRepo) BatchCreateUserSubscription(ctx *gin.Context, userSubsDto []models.Subscription) error {
	writeRequests := []types.WriteRequest{}
	for _, us := range userSubsDto {
		item, _ := attributevalue.MarshalMap(us)
		writeRequests = append(writeRequests,
			types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	batchWrite := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			models.TableSubscription: writeRequests,
		},
	}
	output, er := s.dbClient.BatchWriteItem(ctx, &batchWrite)
	log.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("unable to add subscription to user"))
	}
	return nil
}

func (s *SubscriptionRepo) Create(ctx *gin.Context, model *models.Subscription) error {
	item, _ := attributevalue.MarshalMap(model)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableSubscription),
		Item:      item,
	}

	_, er := s.dbClient.PutItem(ctx, params)

	return er
}
