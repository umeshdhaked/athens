package middleware

import (
	"net/http"

	"github.com/FastBizTech/hastinapura/pkg/services/jwt"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		// if err := ctx.ShouldBindJSON(&req); err != nil {
		// 	fmt.Print(req)
		// 	ctx.Error(err)
		// 	ctx.AbortWithStatus(http.StatusBadRequest)
		// 	return
		// }
		jwtToken := ctx.Request.Header["Token"][0]

		if er := jwt.VerifyToken(jwtToken); er != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
			ctx.Abort()
			return
		}
		claims, _ := jwt.DecodeToken(jwtToken)
		userNme := claims["username"]
		id := claims["id"]
		ctx.AddParam("username", userNme.(string))
		ctx.AddParam("id", id.(string))
		ctx.Next()
	}

}
