package logger

import "go.uber.org/zap"

type env uint8

const (
	EnvTesting env = 0
)

func Setup(env env) *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}
