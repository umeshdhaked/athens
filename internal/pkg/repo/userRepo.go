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

var userRepo *UserRepo

type UserRepo struct {
	Repository
}

func newUserRepo(client *dynamodb.Client) {
	userRepo = &UserRepo{Repository: Repository{dbClient: client}}
}

func GetUserRepo() *UserRepo {
	return userRepo
}

func (u *UserRepo) GetUserFromMobile(ctx context.Context, mobile string) (*models.User, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableUser),
		IndexName: aws.String(models.IndexTableUserIndexMobile),
		KeyConditions: map[string]types.Condition{
			models.ColumnMobile: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: mobile},
				},
			},
		},
	}
	var resp1, err1 = u.dbClient.Query(ctx, queryInput)
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

func (u *UserRepo) UpdateUser(ctx *gin.Context, user *models.User) error {
	item, _ := attributevalue.MarshalMap(user)
	putItem := &dynamodb.PutItemInput{
		TableName: aws.String(models.TableUser),
		Item:      item,
	}

	output, er := u.dbClient.PutItem(ctx, putItem)
	log.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO"))
	}
	return nil
}
