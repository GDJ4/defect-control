package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates structured logger tuned for JSON output in production and human-friendly output locally.
func New(env string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if env != "production" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return cfg.Build()
}
