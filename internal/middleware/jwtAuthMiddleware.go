package middleware

import (
	"net/http"
	"time"

	"github.com/fastbiztech/hastinapura/internal/pkg/jwt"
	"github.com/fastbiztech/hastinapura/internal/services/register"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		jwtToken := ctx.Request.Header["Authorization"][0]

		if er := jwt.VerifyToken(jwtToken); er != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "INVALID_TOKEN"})
			ctx.Abort()
			return
		}
		claims, _ := jwt.DecodeToken(jwtToken)

		exp := claims["exp"].(float64)
		currTime := time.Now().Unix()
		if int64(exp)-currTime < 3600 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "REFRESH_TOKEN"})
			ctx.Abort()
			return
		}
		if int64(exp) < currTime {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "TOKEN_EXPIRED"})
			ctx.Abort()
			return
		}

		userName := claims[constants.JwtTokenMobile].(string)
		id := claims[constants.JwtTokenUserID].(float64)
		role := claims[constants.JwtTokenRole].(string)

		usr, err := register.GetRegistrationService().GetUser(ctx, userName)
		if err != nil || usr == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "USER_NOT_EXIST"})
			ctx.Abort()
			return
		}

		ctx.Set(constants.JwtTokenMobile, userName)
		ctx.Set(constants.JwtTokenUserID, int64(id))
		ctx.Set(constants.JwtTokenRole, role)
		ctx.Next()
	}

}
