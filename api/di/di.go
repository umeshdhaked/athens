package di

import (
	otcSvc "github.com/FastBizTech/hastinapura/api/services/otp"
	"github.com/FastBizTech/hastinapura/api/services/promo"
	"github.com/FastBizTech/hastinapura/api/services/register"
	"github.com/FastBizTech/hastinapura/pkg/models"
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
var dynamoDb *dynamodb.DynamoDB
var otpSender *otp.OtpSender
var otpService *otcSvc.OtpService
var crp *crypto.Crypto
var promoSvc *promo.PromoService

func InitialiseServices(conf *models.ApplicationConfig) {
	sess = aws.ConfigureAwsSdkSession(conf)
	dynamoDb = dynamo.ConfigureDynamoSession(sess)
	otpSender = otp.NewOtpSender(dynamoDb)
	crp = crypto.NewCrypto()
	otpService = otcSvc.NewOtpService(otpSender, crp)
	regService = register.NewRegistrationService(dynamoDb, otpService, crp)
	promoSvc = promo.NewPromoService(dynamoDb)
}

func GetRegistrationService() *register.RegistrationService {
	return regService
}

func GetPromoService() *promo.PromoService {
	return promoSvc
}
