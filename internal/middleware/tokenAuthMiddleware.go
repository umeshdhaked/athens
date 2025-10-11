package middleware

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/gin-gonic/gin"
)

var internalToken string = "0507febe-49a1-4c3e-9da3-7acca0f99d02"

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.Request.Header["Token"][0]

		if internalToken != authToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "INVALID_TOKEN"})
			ctx.Abort()
			return
		}

		ctx.Set(constants.JwtTokenRole, "admin")
		ctx.Set(constants.JwtTokenMobile, "sys-admin-mobile")
		ctx.Set(constants.JwtTokenUserID, int64(-1))
		ctx.Next()
	}

}
