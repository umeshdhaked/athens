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
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("category-index"),
		KeyConditionExpression: aws.String("category = :var0"),
		FilterExpression:       aws.String("sub_category= :var1 and pricing_type = :var2 and pricing_state = :var3"),
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
		TableName: aws.String("pricing"),
		IndexName: aws.String("pricing_state-index"),
		KeyConditions: map[string]types.Condition{
			"pricing_state": {
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
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state = :var0"),
		FilterExpression:       aws.String("pricing_type= :var1"),
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
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state= :var1"),
		FilterExpression:       aws.String("id = :var0"),
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
		TableName: aws.String("pricing"),
		Item:      item,
	}

	output, err := p.client.PutItem(ctx, params)
	log.Println(output)
	return err
}
