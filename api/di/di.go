package di

import (
	otcSvc "github.com/fastbiztech/hastinapura/api/services/otp"
	"github.com/fastbiztech/hastinapura/api/services/promo"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/api/services/subscription"
	pkgAws "github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/internal/pkg/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/repositories"
	"github.com/fastbiztech/hastinapura/internal/services/group"
)

var regService *register.RegistrationService
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
	//conf := config.GetConfig()
	db.NewDb()
	pkgAws.InitialiseS3Client()
	crp = crypto.NewCrypto()

	//repos
	userRepo = repositories.NewUserRepo(db.GetDb().Client)
	subscriptionRepo = repositories.NewSubscriptionRepo(db.GetDb().Client)
	pricingRepo = repositories.NewPricingRepo(db.GetDb().Client)
	promoRepo = repositories.NewPromotionRepo(db.GetDb().Client)
	otpRepo = repositories.NewOtpRepo(db.GetDb().Client)
	creditRepo = repositories.NewCreditsRepo(db.GetDb().Client)
	creditAuditRepo = repositories.NewCreditsAuditRepo(db.GetDb().Client)
	repo.NewRepository(db.GetDb().Client)

	//services
	otpSender = otp.NewOtpSender()
	otpService = otcSvc.NewOtpService(otpSender, crp, otpRepo)
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
