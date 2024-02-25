package servers

import (
	"github.com/FastBizTech/hastinapura/api/handlers"
	"github.com/gin-gonic/gin"
)

func StartPingServer(serverGrp *gin.RouterGroup) {
	serverGrp.POST("/ping", handlers.HandleGetServerPing)
}
