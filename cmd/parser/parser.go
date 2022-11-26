package main

import (
	"context"
	"github.com/kontsevoye/rentaflat/internal/common/logger"
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"github.com/kontsevoye/rentaflat/internal/parser"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			parser.NewConfig,
			logger.NewZapLogger,
			parser.NewFlatFactory,
			parser.CreateDbConnection,
			parser.CreateDbxConnection,
			fx.Annotate(
				parser.NewSqlRepository,
				fx.As(new(parser.Repository)),
			),
			fx.Annotate(
				uuid.NewGoogleGenerator,
				fx.As(new(uuid.Generator)),
			),
			fx.Annotate(
				parser.NewScheduler,
				fx.OnStop(func(_ context.Context, s *parser.Scheduler) error {
					s.Shutdown()
					return nil
				}),
			),
			fx.Annotate(
				parser.NewSsGeParser,
				fx.As(new(parser.Parser)),
			),
		),
		fx.Invoke(func(scheduler *parser.Scheduler) {
			go scheduler.Run()
		}),
	).Run()
}
