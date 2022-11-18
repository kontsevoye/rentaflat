package main

import (
	"context"
	"github.com/kontsevoye/rentaflat/internal/config"
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/flat_scheduler"
	"github.com/kontsevoye/rentaflat/internal/flat_storage"
	"github.com/kontsevoye/rentaflat/internal/flat_subscriber"
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
			flat_subscriber.NewStubSubscriberFactory,
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
		fx.Invoke(func(scheduler *flat_scheduler.Scheduler, storage flat_storage.Storage, subscriberFactory flat_subscriber.StubSubscriberFactory) {
			subscriber, _ := subscriberFactory.NewStubSubscriber()
			subscriber.AddCriteria(flat_subscriber.NewPriceRangeCriteria(600, 700))
			subscriber.AddCriteria(flat_subscriber.NewAreaRangeCriteria(50, 70))
			storage.Subscribe(&subscriber)
			go scheduler.Run()
		}),
	).Run()
}
