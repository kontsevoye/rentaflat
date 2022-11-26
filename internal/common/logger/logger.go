package logger

import (
	"github.com/kontsevoye/rentaflat/internal/parser"
	"go.uber.org/zap"
)

func NewZapLogger(c *parser.AppConfig) *zap.Logger {
	var logger *zap.Logger
	if c.Environment.IsDev() {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	return logger
}
