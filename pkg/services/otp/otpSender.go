package otp

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/google/uuid"
)

type OtpSender struct {
	svc *dynamodb.DynamoDB
}

func NewOtpSender(svc *dynamodb.DynamoDB) *OtpSender {
	return &OtpSender{svc: svc}
}

func (o *OtpSender) GenerateOtp() string {

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (o *OtpSender) SendOtp(otp string) error {
	log.Printf("otp sent: %s", otp)
	//send otp here
	return nil
}

func (o *OtpSender) SaveOtp(mobile string, hashedOtp string) error {
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
	err := req.Send()
	if err != nil {
		return err
	}
	return nil
}

func (o *OtpSender) FetchOtp(mobileNo string) *dbo.Otp {

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("otp"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mobileNo),
					},
				},
			},
		},
	}

	var resp1, err1 = o.svc.Query(queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil
	} else {
		otp := []dbo.Otp{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &otp); err != nil {
			fmt.Println(err)
		}
		return &otp[0]
	}
}
