package handlers

import (
	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/internal/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/internal/pkg/services/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleSendOtp(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	if err := reg.SendOtp(user); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Otp Sent Successful")
}

func HandleRegisterUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	registerResp, err := reg.RegisterUser(user)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleLoginUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	loginResp, err := reg.LoginUser(user)
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
