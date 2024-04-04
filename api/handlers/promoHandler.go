package handlers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/api/services/promo"
	"github.com/fastbiztech/hastinapura/pkg/models/dtos"
	"github.com/gin-gonic/gin"
)

func HandleSaveNumber(ctx *gin.Context) {
	var user dtos.PromoUserRequest
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
	var user dtos.PromoUserRequest
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
	var user dtos.PromoUserRequest
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
