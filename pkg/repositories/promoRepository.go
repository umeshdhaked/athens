package repositories

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
)

type PromotionRepo struct {
	svc *dynamodb.DynamoDB
}

func NewPromotionRepo(svc *dynamodb.DynamoDB) *PromotionRepo {
	return &PromotionRepo{svc: svc}
}

func (p *PromotionRepo) GetPromoFromMobile(mobile string) (*dbo.PromoPhone, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("promo_phones_no"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mobile),
					},
				},
			},
		},
	}
	var resp, er = p.svc.Query(queryInput)
	if er != nil {
		return nil, er
	}

	if *resp.Count > 0 {
		exPromoPh := []dbo.PromoPhone{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
			log.Println(err)
		}

		return &exPromoPh[0], nil
	} else {
		return nil, nil
	}

}

func (p *PromotionRepo) AddPromoContact(obj *dbo.PromoPhone) error {
	item, er := dynamodbattribute.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	req, output := p.svc.PutItemRequest(params)
	log.Print(output)
	return req.Send()
}

func (p *PromotionRepo) GetAlreadyContactedPromo(isAlreadyConnected string) ([]dbo.PromoPhone, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("promo_phones_no"),
		IndexName: aws.String("is_already_contacted-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"is_already_contacted": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(isAlreadyConnected),
					},
				},
			},
		},
	}
	var resp, er = p.svc.Query(queryInput)
	if er != nil {
		return nil, er
	}

	if *resp.Count > 0 {
		exPromoPh := []dbo.PromoPhone{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
			log.Println(err)
			return nil, err
		}
		return exPromoPh, nil
	} else {
		return nil, nil
	}
}

func (p *PromotionRepo) MarkContacted(obj *dbo.PromoPhone) error {
	item, er := dynamodbattribute.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	req, output := p.svc.PutItemRequest(params)
	log.Print(output)
	return req.Send()
}
