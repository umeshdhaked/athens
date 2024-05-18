package serve

import (
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"net/http"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/gin-gonic/gin"
)

func Serve() {
	// Initialise gin engine
	engine := gin.Default()

	// Register routes
	for _, routeGroup := range routeList {
		registerRouteGroup(engine, routeGroup)
	}

	// Initialise go server
	srv := &http.Server{
		Addr:           ":" + config.GetConfig().App.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.GetLogger().Info("application running on port : " + config.GetConfig().App.Port)

	err := srv.ListenAndServe()
	if err != nil {
		panic("failed to launch")
	}
}

func registerRouteGroup(router *gin.Engine, routeGroup route) {
	var middlewareList []gin.HandlerFunc

	if !utils.IsEmpty(routeGroup.middleware) {
		middlewareList = routeGroup.middleware
	}

	group := router.Group(routeGroup.group, middlewareList...)

	for _, endpoint := range routeGroup.endpoints {
		group.Handle(endpoint.method, endpoint.path, endpoint.handler)
	}
}
