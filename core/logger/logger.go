package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(env string) {
	var encoderConfig zapcore.EncoderConfig
	if env == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create a core that writes to both stdout and a file for errors
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	// Ensure logs directory exists
	_ = os.MkdirAll("logs", 0755)
	file, _ := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(file),
		zap.NewAtomicLevelAt(zap.ErrorLevel), // Only log errors to file
	)

	Log = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())
	zap.ReplaceGlobals(Log)
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
