package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured zap.Logger.
func New(env string) *zap.Logger {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, err := cfg.Build()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	return l
}



