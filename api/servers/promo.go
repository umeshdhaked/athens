package servers

import (
	"github.com/FastBizTech/hastinapura/api/handlers"
	"github.com/gin-gonic/gin"
)

func StartPromoServer(serverGrp *gin.RouterGroup) {

	serverGrp.POST("/savePromoNumber", handlers.HandleSaveNumber)
	serverGrp.POST("/fetchPromoNumbers", handlers.HandleFetchPromoNumbers)
	serverGrp.POST("/markContacted", handlers.HandleMarkContactedNumber)

}
