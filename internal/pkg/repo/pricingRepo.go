package repo

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
)

type PricingRepo struct {
	client *dynamodb.Client
}

func NewPricingRepo(client *dynamodb.Client) *PricingRepo {
	return &PricingRepo{client: client}
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
	var resp, err = p.client.Query(ctx, queryInput)
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
			models.ColumnPricingState: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "ACTIVE"},
				},
			},
		},
	}

	var resp, err = p.client.Query(ctx, queryInput)
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

	var pricingResp, err = p.client.Query(ctx, queryInput1)
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

	pricingResp, err := p.client.Query(ctx, queryInput1)
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

	output, err := p.client.PutItem(ctx, params)
	log.Println(output)
	return err
}
