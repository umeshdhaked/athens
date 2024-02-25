package handlers

import (
	"net/http"

	models "github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/gin-gonic/gin"
)

func HandleGetServerPing(ctx *gin.Context) {
	content := models.TesingServer{Message: "Pong"}
	ctx.JSON(http.StatusOK, content)
}
