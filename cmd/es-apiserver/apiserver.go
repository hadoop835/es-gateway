package main

import (
	"github.com/es-gateway/cmd/es-apiserver/app"
	"github.com/es-gateway/pkg/log"
	"go.uber.org/zap"
)

func main() {
	cmd := app.WithApiServerCommand()
	err := cmd.Execute()
	if err != nil {
		log.Log().Fatal("start server error", zap.Error(err))
	}
}
