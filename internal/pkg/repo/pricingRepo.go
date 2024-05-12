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

var pricingRepo *PricingRepo

type PricingRepo struct {
	Repository
}

func newPricingRepo(client *dynamodb.Client) {
	pricingRepo = &PricingRepo{Repository: Repository{dbClient: client}}
}

func GetPricingRepo() *PricingRepo {
	return pricingRepo
}

func (s *PricingRepo) FetchByID(ctx *gin.Context, id string) (*models.Pricing, error) {
	queryInput := dtos.DbQueryInputConditions{
		PKey: map[string]interface{}{
			models.ColumnPricingID: id,
		},
	}

	items, err := s.QueryItems(ctx, models.TablePricing, queryInput)
	if err != nil {
		log.Printf("error fetching item: %v\n", err)
		return nil, err
	}

	if len(items) != 1 {
		log.Println("fetch by id returned more than 1 items")
		return nil, errors.New("fetch by id returned more than 1 items")
	}

	var entity models.Pricing
	if err := attributevalue.UnmarshalMap(items[0], &entity); err != nil {
		log.Println(err)
		return nil, err
	}

	return &entity, nil
}

func (p *PricingRepo) GetDefaultPricingsForCategoryAndSubCategory(ctx *gin.Context, category string, subCategory string) ([]models.Pricing, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName:              aws.String(models.TablePricing),
		IndexName:              aws.String(models.IndexTablePricingIndexCategory),
		KeyConditionExpression: aws.String("Category = :var0"),
		FilterExpression:       aws.String("SubCategory= :var1 and PricingType = :var2 and PricingState = :var3"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: category},
			":var1": &types.AttributeValueMemberS{Value: subCategory},
			":var2": &types.AttributeValueMemberS{Value: "DEFAULT"},
			":var3": &types.AttributeValueMemberS{Value: "ACTIVE"},
		}}
	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		pricing := []models.Pricing{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &pricing); err != nil {
			log.Println(err)
			return nil, err
		}
		return pricing, nil
	}

	return nil, nil
}

func (p *PricingRepo) FetchAllActivePricing(ctx *gin.Context) ([]models.Pricing, error) {

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TablePricing),
		IndexName: aws.String(models.IndexTablePricingIndexPricingState),
		KeyConditions: map[string]types.Condition{
			models.ColumnSubscriptionsPricingStatus: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "ACTIVE"},
				},
			},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.Count > 0 {
		pricing := []models.Pricing{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &pricing); err != nil {
			log.Println(err)
			return nil, err
		}
		return pricing, nil
	}

	return nil, nil
}

func (p *PricingRepo) GetAllDefaultActivePricings(ctx *gin.Context) ([]models.Pricing, error) {
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String(models.TablePricing),
		IndexName:              aws.String(models.IndexTablePricingIndexPricingState),
		KeyConditionExpression: aws.String("PricingState = :var0"),
		FilterExpression:       aws.String("PricingType= :var1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: "ACTIVE"},
			":var1": &types.AttributeValueMemberS{Value: "DEFAULT"},
		},
	}

	var pricingResp, err = p.dbClient.Query(ctx, queryInput1)
	if err != nil {
		return nil, err
	}

	if pricingResp.Count > 0 {
		defaultPricings := []models.Pricing{}
		if err := attributevalue.UnmarshalListOfMaps(pricingResp.Items, &defaultPricings); err != nil {
			log.Println(err)
			return nil, err
		}
		return defaultPricings, nil
	}

	return nil, nil
}

func (p *PricingRepo) GetPricingByPricingID(ctx *gin.Context, pricingId string) (*models.Pricing, error) {
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String(models.TablePricing),
		IndexName:              aws.String(models.IndexTablePricingIndexPricingState),
		KeyConditionExpression: aws.String("PricingState= :var1"),
		FilterExpression:       aws.String("Id = :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberS{Value: pricingId},
			":var1": &types.AttributeValueMemberS{Value: "ACTIVE"},
		},
	}

	pricingResp, err := p.dbClient.Query(ctx, queryInput1)
	if err != nil {
		return nil, err
	}

	if pricingResp.Count > 0 {
		pricings := []models.Pricing{}
		if err := attributevalue.UnmarshalListOfMaps(pricingResp.Items, &pricings); err != nil {
			log.Println(err)
			return nil, err
		}
		return &pricings[0], nil
	}
	return nil, nil
}

func (p *PricingRepo) CreatePricing(ctx *gin.Context, obj *models.Pricing) error {
	item, _ := attributevalue.MarshalMap(obj)
	params := &dynamodb.PutItemInput{
		TableName: aws.String(models.TablePricing),
		Item:      item,
	}

	output, err := p.dbClient.PutItem(ctx, params)
	log.Println(output)
	return err
}
