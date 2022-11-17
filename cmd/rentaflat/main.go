package main

import (
	"context"
	"github.com/kontsevoye/rentaflat/internal/config"
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/flat_scheduler"
	"github.com/kontsevoye/rentaflat/internal/flat_storage"
	"github.com/kontsevoye/rentaflat/internal/logger"
	"github.com/kontsevoye/rentaflat/internal/uuid"
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
			config.NewConfig,
			logger.NewZapLogger,
			fx.Annotate(
				uuid.NewGoogleGenerator,
				fx.As(new(uuid.Generator)),
			),
			fx.Annotate(
				flat_scheduler.NewScheduler,
				fx.OnStop(func(_ context.Context, s *flat_scheduler.Scheduler) error {
					s.Shutdown()
					return nil
				}),
			),
			fx.Annotate(
				flat_storage.NewInMemoryStorage,
				fx.As(new(flat_storage.Storage)),
			),
			fx.Annotate(
				flat_parser.NewSsGeParser,
				fx.As(new(flat_parser.Parser)),
			),
		),
		fx.Invoke(func(scheduler *flat_scheduler.Scheduler, storage flat_storage.Storage, logger *zap.Logger) {
			storage.Subscribe(func(flat flat_parser.Flat) {
				logger.Info(
					"got new flat",
					zap.String("id", flat.Id),
					zap.String("title", flat.Title),
					zap.Int("area", flat.Area),
					zap.Int("price", flat.Price),
					zap.Time("publishedAt", flat.PublishedAt),
				)
			})
			go scheduler.Run()
		}),
	).Run()
}
