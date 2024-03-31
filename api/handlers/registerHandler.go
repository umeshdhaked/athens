package handlers

import (
	"net/http"
	"time"

	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/api/services/register"
	"github.com/fastbiztech/hastinapura/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/pkg/models/responses"
	"github.com/fastbiztech/hastinapura/pkg/services/jwt"
	"github.com/gin-gonic/gin"
)

func HandleSendOtp(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	if err := reg.SendOtp(user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.String(http.StatusOK, "Otp Sent Successful")
}

func HandleRegisterUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	registerResp, err := reg.RegisterUser(user)
	if nil != err {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, registerResp)
}

func HandleLoginUser(ctx *gin.Context) {
	var user requests.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reg *register.RegistrationService = di.GetRegistrationService()
	loginResp, err := reg.LoginUser(user)
	if nil != err {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, loginResp)
}

func RefreshToken(ctx *gin.Context) {
	jwtToken := ctx.Request.Header["Token"][0]

	if er := jwt.VerifyToken(jwtToken); er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "INVALID_TOKEN"})
		ctx.Abort()
		return
	}
	claims, _ := jwt.DecodeToken(jwtToken)
	exp := claims["exp"].(float64)
	currTime := time.Now().Unix()
	userNme := claims["username"].(string)
	role := claims["role"].(string)
	id := claims["id"].(string)
	if int64(exp) < currTime {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "TOKEN_EXPIRED"})
		ctx.Abort()
		return
	}
	if int64(exp)-currTime < 7200 {
		tkn, err := jwt.CreateToken(id, userNme, role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, responses.LoginSuccessResponse{MobileNumber: userNme, LoginToken: tkn})
		return
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "TOKEN_REFRESH_NOT_ALLOWED"})
		ctx.Abort()
		return
	}
}
