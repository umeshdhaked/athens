package main

import (
	"github.com/fastbiztech/hastinapura/internal"
	"log"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/serve"
)

func main() {
	log.Printf("welcome to hastinapur")

	config.LoadConfig()

	internal.InitialiseDeps()

	serve.Serve()
}
