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
	"github.com/fastbiztech/hastinapura/pkg/mutex"
)

// InitialiseDeps *Make sure service are in correct order based on their dependency on each other* //
func InitialiseDeps() {

	// db initialisation
	db.NewDb()
	pkgAws.InitialiseS3Client()

	// mutex connection check
	mutex.ConnectCheck()

	// pkg initialisation

	// repos
	repo.InitialiseRepositories(db.GetDb().Client)

	// services
	crypto.NewCrypto()

	otp.NewOtpSender()
	otcSvc.NewOtpService(otp.GetOtpSender(), crypto.GetCrypto(), repo.GetOtpRepo())
	register.NewRegistrationService(repo.GetUserRepo(), otcSvc.GetOtpService(), crypto.GetCrypto())
	promo.NewPromoService(repo.GetPromotionRepo())
	subscription.NewSubscriptionService(repo.GetPricingRepo(), repo.GetSubscriptionRepo(), repo.GetUserRepo(), repo.GetCreditsRepo(), repo.GetCreditsAuditRepo())

	group.InitialiseService()
	sms.InitialiseService()

}
