package main

import (
	"github.com/eqlabs/sprawl/app"
	"github.com/eqlabs/sprawl/config"
	"go.uber.org/zap"
)

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	appConfig = &config.Config{Logger: log}
	appConfig.ReadConfig("./config/default")
}

func main() {
	app := &app.App{}
	app.InitServices(appConfig, log)
	app.Run()
}
