package servers

import (
	"github.com/FastBizTech/hastinapura/api/handlers"
	"github.com/gin-gonic/gin"
)

func StartRegistrationServer(serverGrp *gin.RouterGroup) {

	serverGrp.POST("/sendOtp", handlers.HandleSendOtp)
	serverGrp.POST("/registerUser", handlers.HandleRegisterUser)
	serverGrp.POST("/login", handlers.HandleLoginUser)
	serverGrp.POST("/savePromoNumber", handlers.HandleSaveNumber)
	serverGrp.POST("/fetchPromoNumbers", handlers.HandleFetchPromoNumbers)
	serverGrp.POST("/markContacted", handlers.HandleMarkContactedNumber)

}
