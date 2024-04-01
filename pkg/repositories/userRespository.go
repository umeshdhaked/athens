package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
)

type UserRepo struct {
	svc *dynamodb.DynamoDB
}

func NewUserRepo(svc *dynamodb.DynamoDB) *UserRepo {
	return &UserRepo{svc: svc}
}

func (u *UserRepo) GetUserFromMobile(mobile string) (*dbo.User, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_table"),
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
	var resp1, err1 = u.svc.Query(queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}
	if *resp1.Count > 0 {
		users := []dbo.User{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &users); err != nil {
			fmt.Println(err)
		}
		return &users[0], nil
	} else {
		return nil, nil
	}
}

func (u *UserRepo) CreateUser(user *dbo.User) error {
	item, _ := dynamodbattribute.MarshalMap(user)
	putItem := &dynamodb.PutItemInput{
		TableName: aws.String("user_table"),
		Item:      item,
	}

	req, output := u.svc.PutItemRequest(putItem)
	log.Println(output)
	er := req.Send()
	if er != nil {
		return errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO"))
	}
	return nil
}
