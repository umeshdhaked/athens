package serve

import (
	"net/http"

	"github.com/fastbiztech/hastinapura/api/handlers"
	"github.com/fastbiztech/hastinapura/api/handlers/middleware"
	"github.com/fastbiztech/hastinapura/internal/controllers"
	"github.com/gin-gonic/gin"
)

// type handlerFunc func(ctx *gin.Context) (interface{}, error, int)
//type handlerFunc func(ctx *gin.Context)

type route struct {
	group      string
	middleware []gin.HandlerFunc
	endpoints  []endpoint
}

type endpoint struct {
	method  string
	path    string
	handler func(ctx *gin.Context)
}

var routeList = [...]route{
	{
		group:      "/ping",
		middleware: []gin.HandlerFunc{},
		endpoints: []endpoint{
			{http.MethodGet, "", handlers.HandleGetServerPing},
		},
	},
	{
		group:      "/v1",
		middleware: []gin.HandlerFunc{middleware.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodGet, "", nil},
		},
	},
	{
		group:      "/v1/users",
		middleware: []gin.HandlerFunc{middleware.TokenAuthMiddleware()},
		endpoints: []endpoint{
			{http.MethodPost, "/sendOtp", handlers.HandleSendOtp},
			{http.MethodPost, "/registerUser", handlers.HandleRegisterUser},
			{http.MethodPost, "/login", handlers.HandleLoginUser},
			{http.MethodPost, "/savePromoNumber", handlers.HandleSaveNumber},
			{http.MethodPost, "/fetchPromoNumbers", handlers.HandleFetchPromoNumbers},
			{http.MethodPost, "/markContacted", handlers.HandleMarkContactedNumber},
		},
	},
	{
		group:      "/v1/group",
		middleware: []gin.HandlerFunc{},
		endpoints: []endpoint{
			{http.MethodPost, "/contacts", controllers.UploadGroupContacts},
		},
	},
}
