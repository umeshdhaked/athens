package handlers

import (
	"fmt"
	"net/http"

	models "github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/gin-gonic/gin"
)

func HandleGetServerPing(ctx *gin.Context) {
	mobile, _ := ctx.Params.Get("username")
	id, _ := ctx.Params.Get("id")

	content := models.TestingResponse{Message: fmt.Sprintf("Ping Pong to %s %s", id, mobile)}
	ctx.JSON(http.StatusOK, content)
}
