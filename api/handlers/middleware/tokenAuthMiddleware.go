package middleware

import (
	"net/http"
	"time"

	"github.com/fastbiztech/hastinapura/internal/pkg/services/jwt"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		jwtToken := ctx.Request.Header["Token"][0]

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

		userNme := claims["username"]
		id := claims["id"]
		role := claims["role"]
		ctx.AddParam("username", userNme.(string))
		ctx.AddParam("id", id.(string))
		ctx.AddParam("role", role.(string))
		ctx.Next()
	}

}
