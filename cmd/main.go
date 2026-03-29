package main

import (
	"go.uber.org/zap"

	"github.com/DangeL187/erax"

	"back/internal/app"
	"back/internal/infra/http/server"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	application, err := app.NewApp()
	if err != nil {
		err = erax.Wrap(err, "Failed to create App")
		zap.L().Fatal("\n" + erax.Format(err))
	}

	httpServer := server.NewServer(application)
	err = httpServer.Run("0.0.0.0:8000")
	if err != nil {
		err = erax.Wrap(err, "Failed to run HTTP server")
		zap.L().Fatal("\n" + erax.Format(err))
	}
}
