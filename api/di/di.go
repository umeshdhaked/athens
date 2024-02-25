package di

import (
	otcSvc "github.com/FastBizTech/hastinapura/api/services/otp"
	"github.com/FastBizTech/hastinapura/api/services/register"
	"github.com/FastBizTech/hastinapura/pkg/services/aws"
	"github.com/FastBizTech/hastinapura/pkg/services/crypto"
	"github.com/FastBizTech/hastinapura/pkg/services/dynamo"
	"github.com/FastBizTech/hastinapura/pkg/services/otp"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// var dynamoConnection pkg.DynnamoConnection
var regService *register.RegistrationService
var sess *session.Session
var dynamoVar *dynamodb.DynamoDB
var otpSender *otp.OtpSender
var otpService *otcSvc.OtpService
var crp *crypto.Crypto

func InitialiseServices() {
	sess = aws.ConfigureAwsSdkSession()
	dynamoVar = dynamo.ConfigureDynamoSession(sess)
	otpSender = otp.NewOtpSender(dynamoVar)
	crp = crypto.NewCrypto()
	otpService = otcSvc.NewOtpService(otpSender, crp)
	regService = register.NewRegistrationService(dynamoVar, otpService, crp)
}

func GetRegistrationService() *register.RegistrationService {
	return regService
}
