package main

import (
	"github.com/umeshdhaked/athens/internal"
	"github.com/umeshdhaked/athens/internal/config"
	"github.com/umeshdhaked/athens/internal/serve"
	"github.com/umeshdhaked/athens/pkg/logger"
)

func main() {
	// logger initialisation
	logger.Build()

	logger.GetLogger().Info("welcome to athens")

	config.LoadConfig()

	internal.InitialiseDeps()

	serve.Serve()
}
