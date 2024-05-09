package internal

import (
	pkgAws "github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/crypto"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/internal/pkg/otp"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/pkg/rzp"
	"github.com/fastbiztech/hastinapura/internal/services/contacts"
	"github.com/fastbiztech/hastinapura/internal/services/group"
	"github.com/fastbiztech/hastinapura/internal/services/invoices"
	otcSvc "github.com/fastbiztech/hastinapura/internal/services/otp"
	"github.com/fastbiztech/hastinapura/internal/services/payments"
	"github.com/fastbiztech/hastinapura/internal/services/pendingJobs"
	"github.com/fastbiztech/hastinapura/internal/services/promo"
	"github.com/fastbiztech/hastinapura/internal/services/register"
	"github.com/fastbiztech/hastinapura/internal/services/s3Processing"
	"github.com/fastbiztech/hastinapura/internal/services/sms"
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/fastbiztech/hastinapura/pkg/mutex"
)

// InitialiseDeps *Make sure service are in correct order based on their dependency on each other* //
func InitialiseDeps() {

	// db initialisation
	db.NewDb()

	// logger initialisation
	logger.Build()

	// pkg initialisation
	pkgAws.InitialiseS3Client()

	// repos
	repo.InitialiseRepositories(db.GetDb().Client)

	// services
	crypto.NewCrypto()

	otp.NewOtpSender()
	otcSvc.NewOtpService(otp.GetOtpSender(), crypto.GetCrypto(), repo.GetOtpRepo())
	promo.NewPromoService(repo.GetPromotionRepo())
	subscription.NewSubscriptionService(repo.GetPricingRepo(), repo.GetSubscriptionRepo(), repo.GetUserRepo(), repo.GetCreditsRepo(), repo.GetCreditsAuditRepo())
	register.NewRegistrationService(repo.GetUserRepo(), otcSvc.GetOtpService(), crypto.GetCrypto(), subscription.GetSubscriptionService())
	rzp.NewRzpService()
	invoices.NewInvoiceService(repo.GetInvoiceRepo())
	payments.NewPaymentService(rzp.GetRzpService(), repo.GetPaymentsRepo(), repo.GetInvoiceRepo(), subscription.GetSubscriptionService())
	payments.NewPaymentCronService(rzp.GetRzpService(), repo.GetPaymentsRepo())

	group.InitialiseService()
	contacts.InitialiseService()
	pendingJobs.InitialiseService()
	s3Processing.InitialiseService()
	sms.InitialiseService()

	// crons initialisation : TODO - move to worker initialisation
	// todo : stop all crons at graceful shutdown
	group.InitialiseS3ContactsCron()
	payments.InitiateRefundForStuckOrdersCron(payments.GetPaymentCronService())

	// mutex connection check
	mutex.ConnectCheck()
	// mutex initialisations
	mutex.Initialise()
}
