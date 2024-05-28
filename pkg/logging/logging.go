package logging

import (
	"go.uber.org/zap"
	"os"
)

var logger *zap.SugaredLogger

func Init() {
	zapLogger, _ := zap.NewProduction()
	logger = zapLogger.Sugar()
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			os.Exit(1)
		}
	}(zapLogger)
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func Info(msg string) {
	logger.Info(msg)
}

func Error(msg string) {
	logger.Error(msg)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}
