package middleware

import (
	"context"
	"github.com/fastbiztech/hastinapura/internal"
	"github.com/fastbiztech/hastinapura/internal/pkg/jwt"
	"net/http"
	"time"

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

		userNme := claims["mobile"]
		id := claims["id"]
		role := claims["role"]

		usr, err := internal.GetRegistrationService().GetUser(context.Background(), userNme.(string))
		if err != nil || usr == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "USER_NOT_EXIST"})
			ctx.Abort()
			return
		}

		ctx.Set("mobile", userNme.(string))
		ctx.Set("id", id.(string))
		ctx.Set("role", role.(string))
		ctx.Next()
	}

}
