package di

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	otcSvc "github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/api/services/promo"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/api/services/subscription"
	"github.com/fastbiztech/hastinapura/internal/config"
	pkgAws "github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
	servAws "github.com/fastbiztech/hastinapura/internal/pkg/services/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/dynamo"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/otp"
	"github.com/fastbiztech/hastinapura/internal/services/group"
)

var regService *register.RegistrationService
var awsConf aws.Config
var dynamoClient *dynamodb.Client
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
var creditRepo *repositories.CreditsRepo
var creditAuditRepo *repositories.CreditsAuditRepo

// InitialiseServices *Make sure service are in correct order based on their dependency on each other* //
func InitialiseDeps() {
	conf := config.GetConfig()

	crp = crypto.NewCrypto()
	awsConf = servAws.ConfigureAwsSdkConfig(conf)
	dynamoClient = dynamo.ConfigureDynamoClient(awsConf)
	// Init services/dependencies
	pkgAws.InitialiseS3Client()

	//repos
	userRepo = repositories.NewUserRepo(dynamoClient)
	subscriptionRepo = repositories.NewSubscriptionRepo(dynamoClient)
	pricingRepo = repositories.NewPricingRepo(dynamoClient)
	promoRepo = repositories.NewPromotionRepo(dynamoClient)
	otpRepo = repositories.NewOtpRepo(dynamoClient)
	creditRepo = repositories.NewCreditsRepo(dynamoClient)
	creditAuditRepo = repositories.NewCreditsAuditRepo(dynamoClient)
	repo.NewRepository(db.GetDb().Client)

	//services
	otpSender = otp.NewOtpSender(otpRepo)
	otpService = otcSvc.NewOtpService(otpSender, crp)
	regService = register.NewRegistrationService(userRepo, otpService, crp)
	promoSvc = promo.NewPromoService(promoRepo)
	subService = subscription.NewSubscriptionService(pricingRepo, subscriptionRepo, userRepo, creditRepo, creditAuditRepo)
	// Group Service
	group.InitialiseService()

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
