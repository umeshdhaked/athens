package di

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	otcSvc "github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/api/services/promo"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/pkg/services/aws"
	"github.com/fastbiztech/hastinapura/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/pkg/services/dynamo"
	"github.com/fastbiztech/hastinapura/pkg/services/otp"
)

// var dynamoConnection pkg.DynnamoConnection
var regService *register.RegistrationService
var sess *session.Session
var dynamoDb *dynamodb.DynamoDB
var otpSender *otp.OtpSender
var otpService *otcSvc.OtpService
var crp *crypto.Crypto
var promoSvc *promo.PromoService

func InitialiseServices(conf *config.Config) {
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
