package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Make init for passing tests.
// Without init() tests fall with SIGSEGV as *zap.Logger is nil.
// Better do with environment variables.
func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot init logger")
	}
}

func InitLogger(developMode bool) error {
	var err error
	if developMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}
	return nil
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}
