package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logConfig.Level = zap.NewAtomicLevelAt(zapLevel)

	return logConfig.Build()
}
