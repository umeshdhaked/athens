package handlers

import (
	"net/http"

	"github.com/FastBizTech/hastinapura/api/di"
	"github.com/FastBizTech/hastinapura/api/services/promo"
	"github.com/FastBizTech/hastinapura/pkg/models/requests"
	"github.com/gin-gonic/gin"
)

func HandleSaveNumber(ctx *gin.Context) {
	var user requests.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *promo.PromoService = di.GetPromoService()
	err := reg.SavePhoneNo(user.MobileNumber)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.String(http.StatusOK, "Saved phone number")
}

func HandleFetchPromoNumbers(ctx *gin.Context) {
	var user requests.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var reg *promo.PromoService = di.GetPromoService()
	list, err := reg.FetchPromoNumbers(user.Is_already_contacted)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.JSON(http.StatusOK, list)
}

func HandleMarkContactedNumber(ctx *gin.Context) {
	var user requests.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *promo.PromoService = di.GetPromoService()
	err := reg.MarkContacted(user.MobileNumber, user.Comment)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.String(http.StatusOK, "Marked as contacted")
}
