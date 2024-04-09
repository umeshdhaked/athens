package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetServerPing(ctx *gin.Context) {
	//content := models.Testing{Message: "Pong"}
	ctx.JSON(http.StatusOK, "Okay")
}
