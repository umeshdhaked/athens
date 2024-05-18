package main

import (
	"github.com/fastbiztech/hastinapura/internal"
	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/serve"
	"github.com/fastbiztech/hastinapura/pkg/logger"
)

func main() {
	// logger initialisation
	logger.Build()

	logger.GetLogger().Info("welcome to Hastinapur")

	config.LoadConfig()

	internal.InitialiseDeps()

	serve.Serve()
}
