package repositories

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
)

type PricingRepo struct {
	svc *dynamodb.DynamoDB
}

func NewPricingRepo(svc *dynamodb.DynamoDB) *PricingRepo {
	return &PricingRepo{svc: svc}
}

func (p *PricingRepo) GetDefaultPricingsForCategoryAndSubCategory(category string, subCategory string) ([]dbo.Pricing, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("category-index"),
		KeyConditionExpression: aws.String("category = :var0"),
		FilterExpression:       aws.String("sub_category= :var1 and pricing_type = :var2 and pricing_state = :var3"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":var0": {S: aws.String(category)},
			":var1": {S: aws.String(subCategory)},
			":var2": {S: aws.String("DEFAULT")},
			":var3": {S: aws.String("ACTIVE")},
		}}
	var resp, err = p.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}
	if *resp.Count > 0 {
		pricing := []dbo.Pricing{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &pricing); err != nil {
			log.Println(err)
			return nil, err
		}
		return pricing, nil
	}

	return nil, nil
}

func (p *PricingRepo) FetchAllActivePricing() ([]dbo.Pricing, error) {

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("pricing"),
		IndexName: aws.String("pricing_state-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"pricing_state": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("ACTIVE"),
					},
				},
			},
		},
	}

	var resp, err = p.svc.Query(queryInput)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if *resp.Count > 0 {
		pricing := []dbo.Pricing{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &pricing); err != nil {
			log.Println(err)
			return nil, err
		}
		return pricing, nil
	}

	return nil, nil
}

func (p *PricingRepo) GetAllDefaultActivePricings() ([]dbo.Pricing, error) {
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state = :var0"),
		FilterExpression:       aws.String("pricing_type= :var1"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":var0": {S: aws.String("ACTIVE")},
			":var1": {S: aws.String("DEFAULT")},
		},
	}

	var pricingResp, err = p.svc.Query(queryInput1)
	if err != nil {
		return nil, err
	}

	if *pricingResp.Count > 0 {
		defaultPricings := []dbo.Pricing{}
		if err := dynamodbattribute.UnmarshalListOfMaps(pricingResp.Items, &defaultPricings); err != nil {
			log.Println(err)
			return nil, err
		}
		return defaultPricings, nil
	}

	return nil, nil
}

func (p *PricingRepo) GetPricingByPricingID(pricingId string) (*dbo.Pricing, error) {
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state= :var1"),
		FilterExpression:       aws.String("id = :var0"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":var0": {S: aws.String(pricingId)},
			":var1": {S: aws.String("ACTIVE")},
		},
	}

	pricingResp, err := p.svc.Query(queryInput1)
	if err != nil {
		return nil, err
	}

	if *pricingResp.Count > 0 {
		pricings := []dbo.Pricing{}
		if err := dynamodbattribute.UnmarshalListOfMaps(pricingResp.Items, &pricings); err != nil {
			log.Println(err)
			return nil, err
		}
		return &pricings[0], nil
	}
	return nil, nil
}

func (p *PricingRepo) CreatePricing(obj *dbo.Pricing) error {
	item, _ := dynamodbattribute.MarshalMap(obj)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("pricing"),
		Item:      item,
	}

	req, output := p.svc.PutItemRequest(params)
	log.Print(output)
	return req.Send()
}
