package controllers

import (
	"github.com/fastbiztech/hastinapura/internal"
	"github.com/fastbiztech/hastinapura/internal/services/register"
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/pkg/jwt"
	"github.com/fastbiztech/hastinapura/pkg/dtos"

	"github.com/gin-gonic/gin"
)

func HandleSendOtp(ctx *gin.Context) {
	var user dtos.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = internal.GetRegistrationService()
	if err := reg.SendOtp(ctx, user); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Otp Sent Successful")
}

func HandleRegisterUser(ctx *gin.Context) {
	var user dtos.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = internal.GetRegistrationService()
	registerResp, err := reg.RegisterUser(ctx, user)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleUpdateUserRoleToAdmin(ctx *gin.Context) {
	var user dtos.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = internal.GetRegistrationService()
	registerResp, err := reg.UpdateUserRoleToAdmin(ctx, user)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleLoginUser(ctx *gin.Context) {
	var user dtos.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = internal.GetRegistrationService()
	loginResp, err := reg.LoginUser(ctx, user)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, loginResp)
}

func HandleRefreshToken(ctx *gin.Context) {
	resp, err := jwt.RefreshToken(ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
