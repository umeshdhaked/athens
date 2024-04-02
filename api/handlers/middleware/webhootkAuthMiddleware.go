package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var internalToken string = "0507febe-49a1-4c3e-9da3-7acca0f99d02"

func WebhookAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.Request.Header["Authorization"][0]

		if internalToken != authToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "INVALID_TOKEN"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}
