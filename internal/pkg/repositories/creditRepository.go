package repositories

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/dbo"
	"github.com/gin-gonic/gin"
	"log"
)

type CreditsRepo struct {
	client *dynamodb.Client
}

func NewCreditsRepo(client *dynamodb.Client) *CreditsRepo {
	return &CreditsRepo{client: client}
}

func (c *CreditsRepo) CreateUserCredit(ctx *gin.Context, credit *dbo.Credits) error {
	item, _ := attributevalue.MarshalMap(credit)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("credit"),
		Item:      item,
	}

	output, er := c.client.PutItem(ctx, params)
	log.Println(output)
	return er
}

func (c *CreditsRepo) FetchUserCredit(ctx *gin.Context, userId string) (*dbo.Credits, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("credit"),
		IndexName: aws.String("user_id-index"),
		KeyConditions: map[string]types.Condition{
			"user_id": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: userId},
				},
			},
		},
	}
	var resp, err = c.client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		credits := []dbo.Credits{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &credits); err != nil {
			log.Println(err)
			return nil, err
		}
		return &credits[0], nil
	}
	return nil, nil
}
