package handlers

import (
	"fmt"
	"net/http"

	"github.com/FastBizTech/hastinapura/api/di"
	"github.com/FastBizTech/hastinapura/api/services"
	models "github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/gin-gonic/gin"
)

func HandleSendOtp(ctx *gin.Context) {
	var user models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Print(user)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *services.RegistrationService = di.GetRegistrationService()
	resp, _ := reg.SendOtp(user)
	m := make(map[string]bool)
	m["otpSent"] = resp
	ctx.JSON(http.StatusOK, m)
}

func HandleRegisterUser(ctx *gin.Context) {
	var user models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Print(user)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *services.RegistrationService = di.GetRegistrationService()
	registerResp, _ := reg.RegisterUser(user)

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleLoginUser(ctx *gin.Context) {
	var user models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Print(user)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var reg *services.RegistrationService = di.GetRegistrationService()
	loginResp, _ := reg.LoginUser(user)

	ctx.JSON(http.StatusOK, loginResp)
}
