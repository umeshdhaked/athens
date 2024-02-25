package register

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/FastBizTech/hastinapura/api/services/otp"
	"github.com/FastBizTech/hastinapura/pkg/models/dbo"
	"github.com/FastBizTech/hastinapura/pkg/models/requests"
	"github.com/FastBizTech/hastinapura/pkg/models/responses"
	"github.com/FastBizTech/hastinapura/pkg/services/crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RegistrationService struct {
	svc        *dynamodb.DynamoDB
	otpService *otp.OtpService
	cryp       *crypto.Crypto
}

func NewRegistrationService(svc *dynamodb.DynamoDB, otpService *otp.OtpService, cryp *crypto.Crypto) *RegistrationService {
	return &RegistrationService{svc: svc, otpService: otpService, cryp: cryp}
}

func (s *RegistrationService) SendOtp(user requests.RegisterUserRequest) error {
	err := s.otpService.SendOtp(user.MobileNumber)
	// add OTP send logic here.
	log.Println("otp sent for user ", user.MobileNumber)
	// save hashed otp after sending otp....
	return err
}

func (s *RegistrationService) RegisterUser(user requests.RegisterUserRequest) (*responses.LoginSuccessResponse, error) {
	log.Println("Received register user request for user ", user)
	if err := s.otpService.VerifyOtp(user.MobileNumber, user.Otp); err != nil {
		return nil, err
	}

	obj := dbo.User{Mobile: user.MobileNumber, Hashed_password: s.cryp.HashString(user.Password)}
	item, _ := dynamodbattribute.MarshalMap(obj)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("user_table"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	er := req.Send()
	if er != nil {
		fmt.Errorf("Failed to make Query API call, %v", er)
	}

	//generate some session token here
	return &responses.LoginSuccessResponse{user.MobileNumber, "some-login-token-after-signup"}, nil
}

func (s *RegistrationService) LoginUser(user requests.RegisterUserRequest) (*responses.LoginSuccessResponse, error) {

	mobile := user.MobileNumber
	password := user.Password

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
	var resp1, err1 = s.svc.Query(queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	} else {
		users := []dbo.User{}
		if err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &users); err != nil {
			fmt.Println(err)
		}
		log.Println(users)

		if s.cryp.HashString(password) == users[0].Hashed_password {
			// generate some token here
			return &responses.LoginSuccessResponse{(users[0]).Mobile, "some-login-token"}, nil
		} else {
			return nil, errors.New("pawword did not match")
		}

	}
}

func (s *RegistrationService) SavePhoneNo(phoneNo string) error {

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

func (s *RegistrationService) FetchPromoNumbers(isAlreadyConnected string) ([]dbo.PromoPhone, error) {

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

func (s *RegistrationService) MarkContacted(mobile string, comment string) error {
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
