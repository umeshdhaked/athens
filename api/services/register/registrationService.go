package register

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/pkg/models/responses"
	"github.com/fastbiztech/hastinapura/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/pkg/services/jwt"
	"github.com/google/uuid"
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

	obj := dbo.User{Id: uuid.New().String(), Mobile: user.MobileNumber, Hashed_password: s.cryp.HashString(user.Password)}
	item, _ := dynamodbattribute.MarshalMap(obj)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("user_table"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	er := req.Send()
	if er != nil {
		return nil, errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO"))
	}

	token, err := jwt.CreateToken(obj.Id, obj.Mobile)
	if err != nil {
		return nil, err
	}

	return &responses.LoginSuccessResponse{MobileNumber: user.MobileNumber, LoginToken: token}, nil
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
			token, err := jwt.CreateToken(users[0].Id, users[0].Mobile)
			if err != nil {
				return nil, err
			}

			return &responses.LoginSuccessResponse{MobileNumber: (users[0]).Mobile, LoginToken: token}, nil
		} else {
			return nil, errors.New("password did not match")
		}

	}
}
