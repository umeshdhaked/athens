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
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *promo.PromoService = di.GetPromoService()
	err := reg.SavePhoneNo(ctx, user.MobileNumber)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Saved phone number")
}

func HandleFetchPromoNumbers(ctx *gin.Context) {
	var user dtos.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	var reg *promo.PromoService = di.GetPromoService()
	list, err := reg.FetchPromoNumbers(ctx, user.Is_already_contacted)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func HandleMarkContactedNumber(ctx *gin.Context) {
	var user dtos.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *promo.PromoService = di.GetPromoService()
	err := reg.MarkContacted(ctx, user.MobileNumber, user.Comment)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Marked as contacted")
}
