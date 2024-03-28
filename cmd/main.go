package main

import (
	"flag"
	"fmt"
	"log"

	"try-on/internal/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var configPath *string = flag.String("c", "config/config.json", "Specify config path")

func main() {
	cfg, err := config.NewDynamicConfig(*configPath,
		nil, func(err error) {
			log.Println("Error parsing config:", err)
		})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Config is: %+v\n", cfg)

	loggerConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		DisableStacktrace: false,
		DisableCaller:     true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.TimeEncoderOfLayout(config.TimeFormat),
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
		},
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
	}

	logger := zap.Must(loggerConfig.Build())
	defer logger.Sync()

	app := NewApp(cfg, logger.Sugar())

	log.Println(app.Run())
}
