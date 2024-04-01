package di

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	otcSvc "github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/api/services/promo"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/api/services/subscription"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/dynamo"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/otp"
)

// var dynamoConnection pkg.DynnamoConnection
var regService *register.RegistrationService
var sess *session.Session
var dynamoDb *dynamodb.DynamoDB
var otpSender *otp.OtpSender
var otpService *otcSvc.OtpService
var crp *crypto.Crypto
var promoSvc *promo.PromoService
var subService *subscription.SubscriptionService
var userRepo *repositories.UserRepo
var subscriptionRepo *repositories.SubscriptionRepo
var pricingRepo *repositories.PricingRepo
var promoRepo *repositories.PromotionRepo
var otpRepo *repositories.OtpRepo

func InitialiseServices(conf *config.Config) {
	crp = crypto.NewCrypto()
	sess = aws.ConfigureAwsSdkSession(conf)
	dynamoDb = dynamo.ConfigureDynamoSession(sess)
	//repos
	userRepo = repositories.NewUserRepo(dynamoDb)
	subscriptionRepo = repositories.NewSubscriptionRepo(dynamoDb)
	pricingRepo = repositories.NewPricingRepo(dynamoDb)
	promoRepo = repositories.NewPromotionRepo(dynamoDb)
	otpRepo = repositories.NewOtpRepo(dynamoDb)
	//services
	regService = register.NewRegistrationService(userRepo, otpService, crp)
	otpSender = otp.NewOtpSender(otpRepo)
	otpService = otcSvc.NewOtpService(otpSender, crp)
	promoSvc = promo.NewPromoService(promoRepo)
	subService = subscription.NewSubscriptionService(pricingRepo, subscriptionRepo, userRepo)
}

func GetRegistrationService() *register.RegistrationService {
	return regService
}

func GetPromoService() *promo.PromoService {
	return promoSvc
}

func GetSubscriptionService() *subscription.SubscriptionService {
	return subService
}
