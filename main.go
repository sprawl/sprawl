package main

import (
	"github.com/sprawl/sprawl/app"
	"github.com/sprawl/sprawl/config"
	"go.uber.org/zap"
)

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger
var configPath = "./config/default"

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	appConfig = &config.Config{Logger: log}
	appConfig.ReadConfig(configPath)
}

func main() {
	app := &app.App{}
	app.InitServices(appConfig, log)
	app.Run()
}
