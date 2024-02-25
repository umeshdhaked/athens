package promo

import (
	"fmt"
	"time"

	"github.com/FastBizTech/hastinapura/pkg/models/dbo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PromoService struct {
	svc *dynamodb.DynamoDB
}

func NewPromoService(svc *dynamodb.DynamoDB) *PromoService {
	return &PromoService{svc: svc}
}

func (s *PromoService) SavePhoneNo(phoneNo string) error {

	//check if already exists ????
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("promo_phones_no"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(phoneNo),
					},
				},
			},
		},
	}
	var resp, er = s.svc.Query(queryInput)
	exPromoPh := []dbo.PromoPhone{}
	if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
		fmt.Println(err)
	}

	obj := dbo.PromoPhone{Mobile: phoneNo, Timestamp: time.Now().Format(time.RFC850)}
	if len(exPromoPh) > 0 && exPromoPh[0].IsAlreadyContacted == "true" {
		return nil
	} else {
		obj.IsAlreadyContacted = "false"
	}

	// then only update with timestamp
	item, er := dynamodbattribute.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	err := req.Send()
	if err != nil {
		return err
	}
	return nil
}

func (s *PromoService) FetchPromoNumbers(isAlreadyConnected string) ([]dbo.PromoPhone, error) {

	//check if already exists ????
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
	var resp, er = s.svc.Query(queryInput)
	if nil != er {
		return nil, er
	}
	exPromoPh := []dbo.PromoPhone{}
	if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &exPromoPh); err != nil {
		fmt.Println(err)
	}
	return exPromoPh, nil
}

func (s *PromoService) MarkContacted(mobile string, comment string) error {
	obj := dbo.PromoPhone{Mobile: mobile,
		Timestamp: time.Now().Format(time.RFC850), IsAlreadyContacted: "true", Comment: comment}

	// then only update with timestamp
	item, er := dynamodbattribute.MarshalMap(obj)
	if er != nil {
		return er
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String("promo_phones_no"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	err := req.Send()
	if err != nil {
		return err
	}
	return nil
}
