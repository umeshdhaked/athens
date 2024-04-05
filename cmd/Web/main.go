package main

import (
	"log"

	"github.com/fastbiztech/hastinapura/api/di"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/internal/serve"
)

func main() {
	log.Printf("welcome to hastinapur")

	config.LoadConfig()
	db.NewDb()

	di.InitialiseDeps()

	serve.Serve()
}
