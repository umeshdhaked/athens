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

type OtpRepo struct {
	client *dynamodb.Client
}

func NewOtpRepo(client *dynamodb.Client) *OtpRepo {
	return &OtpRepo{client: client}
}

func (o *OtpRepo) SaveOtp(ctx *gin.Context, otp dbo.Otp) error {
	item, er := attributevalue.MarshalMap(otp)
	if er != nil {
		return er
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("otp"),
		Item:      item,
	}

	output, err := o.client.PutItem(ctx, params)
	log.Print(output)
	return err
}

func (o *OtpRepo) GetOtp(ctx *gin.Context, mobile string) (*dbo.Otp, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("otp"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]types.Condition{
			"mobile": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{&types.AttributeValueMemberS{Value: mobile}},
			},
		},
	}

	var resp, err = o.client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		otp := []dbo.Otp{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &otp); err != nil {
			log.Println(err)
			return nil, err
		}
		return &otp[0], nil
	}
	return nil, nil

}
