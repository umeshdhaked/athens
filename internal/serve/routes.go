package serve

import (
	middleware2 "github.com/fastbiztech/hastinapura/internal/middleware"
	"net/http"

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
			{http.MethodGet, "", controllers.HandleGetServerPing},
		},
	},
	{
		group:      "/v1",
		middleware: []gin.HandlerFunc{middleware2.JwtAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodGet, "", nil},
		},
	},
	{
		group:      "/v1/users",
		middleware: []gin.HandlerFunc{middleware2.JwtAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/savePromoNumber", controllers.HandleSaveNumber},
			{http.MethodPost, "/fetchPromoNumbers", controllers.HandleFetchPromoNumbers},
			{http.MethodPost, "/markContacted", controllers.HandleMarkContactedNumber},
		},
	},
	{
		group: "/v1/users",
		endpoints: []endpoint{
			{http.MethodPost, "/sendOtp", controllers.HandleSendOtp},
			{http.MethodPost, "/registerUser", controllers.HandleRegisterUser},
			{http.MethodPost, "/login", controllers.HandleLoginUser},
			{http.MethodPost, "/refresh_token", controllers.HandleRefreshToken},
		},
	},
	{
		group:      "/v1/subscriptions",
		middleware: []gin.HandlerFunc{middleware2.JwtAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/createNewPricingSystem", controllers.HandleCreateNewPricingSystem}, //admin api
			{http.MethodPost, "/deactivatePricing", controllers.HandleDeactivatePricing},           //admin api
			{http.MethodGet, "/fetchAllActivePricingModel", controllers.HandleFetchAllActivePricingModel},
			{http.MethodPost, "/addDefaultSubscriptions", controllers.HandleAddDefaultSubscriptionToUser}, //admin api
			{http.MethodPost, "/addCustomSubscriptions", controllers.HandleAddSubscriptionToUser},         //admin api
			{http.MethodPost, "/fetchAllActiveActiveSubscriptionsForUser", controllers.HandleFetchAllActiveSubscriptionsForUser},
			{http.MethodPost, "/deactivateUserSubscription", controllers.HandleDeactivateSubscriptionsForUser},
			// {http.MethodPost, "/addCreditToUser", handlers.HandleAddCreditToUser}, // admin api
			{http.MethodPost, "/fetchCredits", controllers.HandleFetchCredits},
			// {http.MethodPost, "/chargeUser", handlers.HandleChargeUser},
		},
	},
	{
		group:      "/v1/subscriptions",
		middleware: []gin.HandlerFunc{middleware2.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/addCreditToUser", controllers.HandleAddCreditToUser}, // admin api
			{http.MethodPost, "/chargeUser", controllers.HandleChargeUser},           //this was created just for testing
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
