package main

import (
	"context"
	"github.com/kontsevoye/rentaflat/internal/common/logger"
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"github.com/kontsevoye/rentaflat/internal/parser"
	"github.com/kontsevoye/rentaflat/internal/parser/app"
	"github.com/kontsevoye/rentaflat/internal/parser/app/command"
	"github.com/kontsevoye/rentaflat/internal/parser/app/query"
	"github.com/kontsevoye/rentaflat/internal/parser/domain"
	"github.com/kontsevoye/rentaflat/internal/parser/infrastructure"
	"github.com/kontsevoye/rentaflat/internal/parser/port"
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
			domain.NewFlatFactory,
			infrastructure.CreateDbConnection,
			infrastructure.CreateDbxConnection,
			command.NewParseFlatListHandler,
			query.NewGetLastFlatServiceIdHandler,
			app.NewApplication,
			fx.Annotate(
				infrastructure.NewSqlRepository,
				fx.As(new(domain.Repository)),
			),
			fx.Annotate(
				uuid.NewGoogleGenerator,
				fx.As(new(uuid.Generator)),
			),
			fx.Annotate(
				port.NewScheduler,
				fx.OnStop(func(_ context.Context, s *port.Scheduler) error {
					s.Shutdown()
					return nil
				}),
			),
			fx.Annotate(
				infrastructure.NewSsGeParser,
				fx.As(new(domain.Parser)),
			),
		),
		fx.Invoke(func(scheduler *port.Scheduler) {
			go scheduler.Run()
		}),
	).Run()
}
