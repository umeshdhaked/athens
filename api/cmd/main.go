package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FastBizTech/hastinapura/api/di"
	servers "github.com/FastBizTech/hastinapura/api/servers"
	"github.com/FastBizTech/hastinapura/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	log.Printf("welcome to hastinapur")

	viper.AddConfigPath("./config")
	viper.SetConfigName(os.Getenv("ENV"))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &models.ApplicationConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	di.InitialiseServices(conf)

	engine := gin.Default()
	server := servers.NewServer(engine)
	server.Serve()

	srv := &http.Server{
		Addr:           ":" + conf.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("application running on port : " + conf.Port)
	srv.ListenAndServe()
}
