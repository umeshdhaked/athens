package controllers

import (
	"github.com/fastbiztech/hastinapura/internal/services/subscription"
	"net/http"

	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func HandleCreateNewPricingSystem(ctx *gin.Context) {
	var pricingRequest dtos.PricingRequest
	if err := ctx.ShouldBindJSON(&pricingRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()

	resp, err := sub.CreateNewPricingSystem(ctx, &pricingRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func HandleDeactivatePricing(ctx *gin.Context) {
	var pricingRequest dtos.DeactivatePricingRequest
	if err := ctx.ShouldBindJSON(&pricingRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.PricingStatusUpdate(ctx, &pricingRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func HandleFetchAllActivePricingModel(ctx *gin.Context) {

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	resp, err := sub.FetchAllActivePricingModel(ctx)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func HandleAddDefaultSubscriptionToUser(ctx *gin.Context) {
	var subRequest *dtos.UserDefaultSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.AddDefaultSubscriptionToUser(ctx, subRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func HandleAddSubscriptionToUser(ctx *gin.Context) {
	var subRequest *dtos.UserSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.UpdateSubscriptionToUser(ctx, subRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func HandleFetchAllActiveSubscriptionsForUser(ctx *gin.Context) {
	var subRequest *dtos.FetchSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	resp, err := sub.FetchAllActiveSubscriptionsForUser(ctx, subRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func HandleDeactivateSubscriptionsForUser(ctx *gin.Context) {
	var subRequest *dtos.DeactivateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.SubscriptionsStatusUpdate(ctx, subRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "OK")

}

func HandleAddCreditToUser(ctx *gin.Context) {
	var creditRequest *dtos.AddCreditsRequest
	if err := ctx.ShouldBindJSON(&creditRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.AddCreditToUser(ctx, creditRequest)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func HandleFetchCredits(ctx *gin.Context) {
	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	resp, err := sub.FetchCredit(ctx)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleChargeUser(ctx *gin.Context) {
	var chargeRequest *dtos.ChargeUserRequest
	if err := ctx.ShouldBindJSON(&chargeRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub *subscription.SubscriptionService = subscription.GetSubscriptionService()
	err := sub.ChargeUser(ctx, chargeRequest.UserId, chargeRequest.Category, chargeRequest.SubCategory, chargeRequest.UnitCount)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}
