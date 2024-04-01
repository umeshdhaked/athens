package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/google/uuid"
)

type OtpRepo struct {
	svc *dynamodb.DynamoDB
}

func NewOtpRepo(svc *dynamodb.DynamoDB) *OtpRepo {
	return &OtpRepo{svc: svc}
}

func (o *OtpRepo) SaveOtp(mobile string, hashedOtp string) error {
	item, er := dynamodbattribute.MarshalMap(
		dbo.Otp{
			Id:     uuid.New().String(),
			Mobile: mobile,
			Otp:    hashedOtp,
			Exp:    time.Now().Add(2 * time.Minute).Unix(),
		})
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("otp"),
		Item:      item,
	}

	req, output := o.svc.PutItemRequest(params)
	fmt.Print(output)
	return req.Send()
}

func (o *OtpRepo) GetOtp(mobile string) (*dbo.Otp, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("otp"),
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

	var resp, err = o.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}
	if *resp.Count > 0 {
		otp := []dbo.Otp{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &otp); err != nil {
			log.Println(err)
			return nil, err
		}
		return &otp[0], nil
	}
	return nil, nil

}
