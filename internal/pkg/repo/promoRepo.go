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

type PromotionRepo struct {
	client *dynamodb.Client
}

func NewPromotionRepo(client *dynamodb.Client) *PromotionRepo {
	return &PromotionRepo{client: client}
}

func (p *PromotionRepo) GetPromoFromMobile(ctx *gin.Context, mobile string) (*models.PromoPhone, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("promo_phones_no"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]types.Condition{
			"mobile": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: mobile},
				},
			},
		},
	}
	var resp, er = p.client.Query(ctx, queryInput)
	if er != nil {
		return nil, er
	}

	if resp.Count > 0 {
		exPromoPh := []models.PromoPhone{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
			log.Println(err)
		}

		return &exPromoPh[0], nil
	} else {
		return nil, nil
	}

}

func (p *PromotionRepo) AddPromoContact(ctx *gin.Context, obj *models.PromoPhone) error {
	item, er := attributevalue.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	output, err := p.client.PutItem(ctx, params)
	log.Print(output)
	return err
}

func (p *PromotionRepo) GetAlreadyContactedPromo(ctx *gin.Context, isAlreadyConnected string) ([]models.PromoPhone, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("promo_phones_no"),
		IndexName: aws.String("is_already_contacted-index"),
		KeyConditions: map[string]types.Condition{
			"is_already_contacted": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: isAlreadyConnected},
				},
			},
		},
	}
	var resp, er = p.client.Query(ctx, queryInput)
	if er != nil {
		return nil, er
	}

	if resp.Count > 0 {
		exPromoPh := []models.PromoPhone{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
			log.Println(err)
			return nil, err
		}
		return exPromoPh, nil
	} else {
		return nil, nil
	}
}

func (p *PromotionRepo) MarkContacted(ctx *gin.Context, obj *models.PromoPhone) error {
	item, er := attributevalue.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	output, err := p.client.PutItem(ctx, params)
	log.Print(output)
	return err
}
