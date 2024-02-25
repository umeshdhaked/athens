package handlers

import (
	"fmt"
	"net/http"

	"github.com/FastBizTech/hastinapura/api/di"
	"github.com/FastBizTech/hastinapura/api/services/register"
	"github.com/FastBizTech/hastinapura/pkg/models/requests"
	"github.com/gin-gonic/gin"
)

func HandleSendOtp(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Print(user)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	if err := reg.SendOtp(user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.String(http.StatusOK, "Otp Sent Successful")
}

func HandleRegisterUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	registerResp, err := reg.RegisterUser(user)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleLoginUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Print(user)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	loginResp, err := reg.LoginUser(user)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.JSON(http.StatusOK, loginResp)
}

func HandleSaveNumber(ctx *gin.Context) {
	var user requests.PromoUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
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
	var reg *register.RegistrationService = di.GetRegistrationService()
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

	var reg *register.RegistrationService = di.GetRegistrationService()
	err := reg.MarkContacted(user.MobileNumber, user.Comment)
	if nil != err {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.String(http.StatusOK, "Marked as contacted")
}
