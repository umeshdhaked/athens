package servers

import (
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
	StartPingServer(baseServerGroup)
	userServerGroup := s.engine.Group("/api/v1/users")
	StartRegistrationServer(userServerGroup)

}
