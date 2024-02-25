package servers

import (
	"github.com/FastBizTech/hastinapura/api/handlers/middleware"
	"github.com/gin-gonic/gin"
)

type server struct {
	engine *gin.Engine
}

func NewServer(eng *gin.Engine) *server {
	return &server{engine: eng}
}

func (s *server) Serve() {

	baseServerGroup := s.engine.Group("/api/v1")
	baseServerGroup.Use(middleware.TokenAuthMiddleware())
	StartPingServer(baseServerGroup)
	userServerGroup := s.engine.Group("/api/v1/users")
	StartRegistrationServer(userServerGroup)
	StartPromoServer(userServerGroup)
}
