package handlers

import (
	"fmt"
	"net/http"

	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/pkg/models/requests"
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
