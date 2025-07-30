package logger

import (
	"log"

	"go.uber.org/zap"
)

var L *zap.Logger

func Init() {
	var err error
	L, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
}

func Info(msg string, fields ...zap.Field) {
	L.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	L.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	L.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	L.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	L.Fatal(msg, fields...)
}
