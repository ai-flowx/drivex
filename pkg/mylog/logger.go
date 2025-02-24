package mylog

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLog(mode string) {
	var encoder zapcore.Encoder
	var encoderConfig zapcore.EncoderConfig
	var level zapcore.Level

	log.Println("level mode", mode)

	switch mode {
	case "prod":
		encoderConfig = zap.NewProductionEncoderConfig()
		level = zapcore.WarnLevel
	case "dev":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		level = zapcore.InfoLevel
	case "debug":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		level = zapcore.DebugLevel
	default:
		log.Println("level mode default prod")
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		level = zapcore.WarnLevel
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	log.Println("log plain-text format")

	core := zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)

	Logger = zap.New(core, zap.AddCaller())
}
