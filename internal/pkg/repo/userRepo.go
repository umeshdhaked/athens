package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepo struct {
	client *dynamodb.Client
}

func NewUserRepo(client *dynamodb.Client) *UserRepo {
	return &UserRepo{client: client}
}

func (u *UserRepo) GetUserFromMobile(ctx context.Context, mobile string) (*models.User, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_table"),
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
	var resp1, err1 = u.client.Query(ctx, queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}
	if resp1.Count > 0 {
		users := []models.User{}
		if err := attributevalue.UnmarshalListOfMaps(resp1.Items, &users); err != nil {
			fmt.Println(err)
		}
		return &users[0], nil
	} else {
		return nil, nil
	}
}

func (u *UserRepo) CreateUser(ctx *gin.Context, user *models.User) error {
	item, _ := attributevalue.MarshalMap(user)
	putItem := &dynamodb.PutItemInput{
		TableName: aws.String("user_table"),
		Item:      item,
	}

	output, er := u.client.PutItem(ctx, putItem)
	log.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO"))
	}
	return nil
}
