package serve

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/api/handlers"
	"github.com/fastbiztech/hastinapura/api/handlers/middleware"
	"github.com/fastbiztech/hastinapura/internal/controllers"
	"github.com/gin-gonic/gin"
)

// type handlerFunc func(ctx *gin.Context) (interface{}, error, int)
//type handlerFunc func(ctx *gin.Context)

type route struct {
	group      string
	middleware []gin.HandlerFunc
	endpoints  []endpoint
}

type endpoint struct {
	method  string
	path    string
	handler func(ctx *gin.Context)
}

var routeList = [...]route{
	{
		group:      "/ping",
		middleware: []gin.HandlerFunc{},
		endpoints: []endpoint{
			{http.MethodGet, "", handlers.HandleGetServerPing},
		},
	},
	{
		group:      "/v1",
		middleware: []gin.HandlerFunc{middleware.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodGet, "", nil},
		},
	},
	{
		group:      "/v1/users",
		middleware: []gin.HandlerFunc{middleware.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/savePromoNumber", handlers.HandleSaveNumber},
			{http.MethodPost, "/fetchPromoNumbers", handlers.HandleFetchPromoNumbers},
			{http.MethodPost, "/markContacted", handlers.HandleMarkContactedNumber},
		},
	},
	{
		group: "/v1/users",
		endpoints: []endpoint{
			{http.MethodPost, "/sendOtp", handlers.HandleSendOtp},
			{http.MethodPost, "/registerUser", handlers.HandleRegisterUser},
			{http.MethodPost, "/login", handlers.HandleLoginUser},
			{http.MethodPost, "/refresh_token", handlers.HandleRefreshToken},
		},
	},
	{
		group:      "/v1/subscriptions",
		middleware: []gin.HandlerFunc{middleware.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/createNewPricingSystem", handlers.HandleCreateNewPricingSystem}, //admin api
			{http.MethodPost, "/deactivatePricing", handlers.HandleDeactivatePricing},           //admin api
			{http.MethodGet, "/fetchAllActivePricingModel", handlers.HandleFetchAllActivePricingModel},
			{http.MethodPost, "/addDefaultSubscriptions", handlers.HandleAddDefaultSubscriptionToUser}, //admin api
			{http.MethodPost, "/addCustomSubscriptions", handlers.HandleAddSubscriptionToUser},         //admin api
			{http.MethodPost, "/fetchAllActiveActiveSubscriptionsForUser", handlers.HandleFetchAllActiveSubscriptionsForUser},
			{http.MethodPost, "/deactivateUserSubscription", handlers.HandleDeactivateSubscriptionsForUser},
			// {http.MethodPost, "/addCreditToUser", handlers.HandleAddCreditToUser}, // admin api
			{http.MethodPost, "/fetchCredits", handlers.HandleFetchCredits},
			// {http.MethodPost, "/chargeUser", handlers.HandleChargeUser},
		},
	},
	{
		group:      "/v1/subscriptions",
		middleware: []gin.HandlerFunc{middleware.WebhookAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/addCreditToUser", handlers.HandleAddCreditToUser}, // admin api
			{http.MethodPost, "/chargeUser", handlers.HandleChargeUser},           //this was created just for testing
		},
	},
	{
		group:      "/v1/group",
		middleware: []gin.HandlerFunc{},
		endpoints: []endpoint{
			{http.MethodPost, "/contacts", controllers.UploadGroupContacts},
			{http.MethodGet, "/contacts", controllers.GetGroupContacts},
		},
	},
	{
		group:      "/v1/sms",
		middleware: []gin.HandlerFunc{},
		endpoints: []endpoint{
			// Sender Id related apis
			{http.MethodPost, "/senderid", controllers.PostSenderCode},
			{http.MethodGet, "/senderid", controllers.GetSenderCode},
			{http.MethodPost, "/senderid/approve", controllers.ApproveSenderCode},
			{http.MethodPatch, "/senderid/deactivate", controllers.DeActivateSenderCode},

			// Template related apis
			{http.MethodPost, "/template", controllers.PostSmsTemplate},
			{http.MethodPost, "/template/approve", controllers.ApproveSmsTemplate},
			{http.MethodGet, "/template", controllers.GetSmsTemplate},
			{http.MethodPatch, "/template", controllers.UpdateSmsTemplate},
			{http.MethodPatch, "/template/deactivate", controllers.DeActivateSmsTemplate},

			// instant sms api
			{http.MethodPost, "", controllers.PostSms},
			{http.MethodPost, "/retry", controllers.PostSms}, // TODO validate if new method needed for retry sms

			// Sms Reporting
			//{http.MethodGet, "", controllers.UploadGroupContacts},
		},
	},
}
