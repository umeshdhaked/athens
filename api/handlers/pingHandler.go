package handlers

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/internal/pkg/models"
	"github.com/gin-gonic/gin"
)

func HandleGetServerPing(ctx *gin.Context) {
	content := models.Testing{Message: "Pong"}
	ctx.JSON(http.StatusOK, content)
}
