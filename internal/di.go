package internal

import (
	pkgAws "github.com/umeshdhaked/athens/internal/pkg/aws"
	"github.com/umeshdhaked/athens/internal/pkg/crypto"
	"github.com/umeshdhaked/athens/internal/pkg/db"
	"github.com/umeshdhaked/athens/internal/pkg/otp"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/internal/pkg/rzp"
	"github.com/umeshdhaked/athens/internal/services/contacts"
	"github.com/umeshdhaked/athens/internal/services/cronProcessing"
	"github.com/umeshdhaked/athens/internal/services/group"
	"github.com/umeshdhaked/athens/internal/services/invoices"
	"github.com/umeshdhaked/athens/internal/services/kyc"
	otcSvc "github.com/umeshdhaked/athens/internal/services/otp"
	"github.com/umeshdhaked/athens/internal/services/payments"
	"github.com/umeshdhaked/athens/internal/services/pendingJobs"
	"github.com/umeshdhaked/athens/internal/services/promo"
	"github.com/umeshdhaked/athens/internal/services/register"
	"github.com/umeshdhaked/athens/internal/services/sms"
	"github.com/umeshdhaked/athens/internal/services/subscription"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/umeshdhaked/athens/pkg/mutex"
)

// InitialiseDeps *Make sure service are in correct order based on their dependency on each other* //
func InitialiseDeps() {

	// logger initialisation
	logger.Build()

	// db initialisation
	db.NewDb()

	// repos
	repo.InitialiseRepositories()

	// pkg initialisation
	pkgAws.InitialiseS3Client()

	initServices()

	initCrons()

	// mutex initialisations
	mutex.Initialise()
}

func initServices() {
	// services
	crypto.NewCrypto()

	otp.NewOtpSender()
	kyc.NewKycService()
	otcSvc.NewOtpService(otp.GetOtpSender(), crypto.GetCrypto())
	promo.NewPromoService(repo.GetPromotionRepo())
	subscription.NewSubscriptionService(repo.GetPricingRepo(), repo.GetSubscriptionRepo(), repo.GetUserRepo(), repo.GetCreditsRepo(), repo.GetCreditsAuditRepo())
	register.NewRegistrationService(otcSvc.GetOtpService(), crypto.GetCrypto(), subscription.GetSubscriptionService())
	rzp.NewRzpService()
	invoices.NewInvoiceService(repo.GetInvoiceRepo())
	payments.NewPaymentService(rzp.GetRzpService(), repo.GetPaymentsRepo(), repo.GetInvoiceRepo(), subscription.GetSubscriptionService())
	//payments.NewPaymentCronService(rzp.GetRzpService(), repo.GetPaymentsRepo())

	group.InitialiseService()
	contacts.InitialiseService()
	pendingJobs.InitialiseService()
	cronProcessing.InitialiseService()
	sms.InitialiseService()
}

func initCrons() {
	// crons initialisation : TODO - move to worker initialisation
	// todo : stop all crons at graceful shutdown

	group.InitialiseS3ContactsCron()
	sms.InitialiseCampaignCron()
}
