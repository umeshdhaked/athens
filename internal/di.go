package internal

import (
	pkgAws "github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/internal/pkg/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/services/group"
	otcSvc "github.com/fastbiztech/hastinapura/internal/services/otp"
	"github.com/fastbiztech/hastinapura/internal/services/promo"
	"github.com/fastbiztech/hastinapura/internal/services/register"
	"github.com/fastbiztech/hastinapura/internal/services/sms"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
)

var regService *register.RegistrationService
var otpSender *otp.OtpSender
var otpService *otcSvc.OtpService
var crp *crypto.Crypto
var promoSvc *promo.PromoService
var subService *subscription.SubscriptionService
var userRepo *repo.UserRepo
var subscriptionRepo *repo.SubscriptionRepo
var pricingRepo *repo.PricingRepo
var promoRepo *repo.PromotionRepo
var otpRepo *repo.OtpRepo
var creditRepo *repo.CreditsRepo
var creditAuditRepo *repo.CreditsAuditRepo

// InitialiseDeps *Make sure service are in correct order based on their dependency on each other* //
func InitialiseDeps() {

	// db initialisation
	db.NewDb()

	// pkg initialisation
	pkgAws.InitialiseS3Client()
	crp = crypto.NewCrypto()

	// repos
	repo.NewRepository(db.GetDb().Client)
	userRepo = repo.NewUserRepo(db.GetDb().Client)
	subscriptionRepo = repo.NewSubscriptionRepo(db.GetDb().Client)
	pricingRepo = repo.NewPricingRepo(db.GetDb().Client)
	promoRepo = repo.NewPromotionRepo(db.GetDb().Client)
	otpRepo = repo.NewOtpRepo(db.GetDb().Client)
	creditRepo = repo.NewCreditsRepo(db.GetDb().Client)
	creditAuditRepo = repo.NewCreditsAuditRepo(db.GetDb().Client)

	// services
	otpSender = otp.NewOtpSender()
	otpService = otcSvc.NewOtpService(otpSender, crp, otpRepo)
	regService = register.NewRegistrationService(userRepo, otpService, crp)
	promoSvc = promo.NewPromoService(promoRepo)
	subService = subscription.NewSubscriptionService(pricingRepo, subscriptionRepo, userRepo, creditRepo, creditAuditRepo)

	group.InitialiseService()
	sms.InitialiseService()

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
