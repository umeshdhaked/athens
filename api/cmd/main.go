package main

import (
	"log"
	"net/http"
	"time"

	"github.com/FastBizTech/hastinapura/api/di"
	servers "github.com/FastBizTech/hastinapura/api/servers"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("welcome to hastinapur")

	di.InitialiseServices()

	engine := gin.Default()
	server := servers.NewServer(engine)
	server.Serve()

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	srv.ListenAndServe()
}
