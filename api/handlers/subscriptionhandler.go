package handlers

import (
	"fmt"
	"net/http"

	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/api/services/subscription"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/requests"
	"github.com/gin-gonic/gin"
)

func HandleCreateNewPricingSystem(ctx *gin.Context) {
	var pricingRequest requests.PricingRequest
	if err := ctx.ShouldBindJSON(&pricingRequest); err != nil {
		fmt.Print(pricingRequest)
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var sub *subscription.SubscriptionService = di.GetSubscriptionService()

	resp, err := sub.CreateNewPricingSystem(ctx, &pricingRequest)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func HandleFetchAllActivePricingModel(ctx *gin.Context) {
	var sub *subscription.SubscriptionService = di.GetSubscriptionService()
	resp, err := sub.FetchAllActivePricingModel(ctx)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func HandleAddDefaultSubscriptionToUser(ctx *gin.Context) {
	var subRequest *requests.UserDefaultSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		fmt.Print(subRequest)
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var sub *subscription.SubscriptionService = di.GetSubscriptionService()
	err := sub.AddDefaultSubscriptionToUser(ctx, subRequest)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func HandleAddSubscriptionToUser(ctx *gin.Context) {
	var subRequest *requests.UserSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		fmt.Print(subRequest)
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var sub *subscription.SubscriptionService = di.GetSubscriptionService()
	err := sub.AddSubscriptionToUser(ctx, subRequest)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func HandleFetchAllActiveSubscriptionsForUser(ctx *gin.Context) {
	var subRequest *requests.FetchSubscriptionRequest
	if err := ctx.ShouldBindJSON(&subRequest); err != nil {
		fmt.Print(subRequest)
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var sub *subscription.SubscriptionService = di.GetSubscriptionService()
	resp, err := sub.FetchAllActiveSubscriptionsForUser(ctx, subRequest)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
