package main

import (
	"github.com/sprawl/sprawl/app"
	"github.com/sprawl/sprawl/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger
var configPath = "./config/default"

func init() {
	// Read config
	appConfig = &config.Config{}
	appConfig.ReadConfig(configPath)

	// Read logLevel from appConfig
	var logLevel zapcore.Level
	switch appConfig.GetLogLevel() {
	case "DEBUG":
		logLevel = zapcore.DebugLevel
	case "INFO":
		logLevel = zapcore.InfoLevel
	case "WARN":
		logLevel = zapcore.WarnLevel
	case "ERROR":
		logLevel = zapcore.ErrorLevel
	case "PANIC":
		logLevel = zapcore.PanicLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	// Read logFormat ("console"/"json") from appConfig
	logFormat := appConfig.GetLogFormat()
	if logFormat == "" {
		logFormat = "json"
	}

	// Create logger
	cfg := zap.Config{
		Encoding:         logFormat,
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ = cfg.Build()
	log = logger.Sugar()
}

func main() {
	app := &app.App{}
	app.InitServices(appConfig, log)
	app.Run()
}
