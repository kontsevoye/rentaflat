package main

import (
	"database/sql"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kontsevoye/rentaflat/internal/common/logger"
	"github.com/kontsevoye/rentaflat/internal/parser"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	up := flag.Bool("up", false, "execute all available migrations")
	down := flag.Bool("down", false, "revert all migrations")
	flag.Parse()

	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			parser.NewConfig,
			logger.NewZapLogger,
			parser.CreateDbConnection,
		),
		fx.Invoke(func(log *zap.Logger, connection *sql.DB, shutdowner fx.Shutdowner) {
			if !*up && !*down {
				log.Fatal("you must specify one of [\"-up\", \"-down\"] options")
				return
			}
			if *up && *down {
				log.Fatal("you must specify only one of [\"-up\", \"-down\"] options")
				return
			}

			driver, err := postgres.WithInstance(connection, &postgres.Config{})
			if err != nil {
				log.Fatal("error while injecting db connection", zap.Error(err))
				return
			}
			m, err := migrate.NewWithDatabaseInstance(
				"file://migrations",
				"postgres",
				driver,
			)
			if err != nil {
				log.Fatal("error while initiating migrate", zap.Error(err))
				return
			}

			if *up {
				err = m.Up()
			} else if *down {
				err = m.Down()
			}

			if err != nil && err != migrate.ErrNoChange {
				log.Fatal("error while running command", zap.Error(err))
				return
			}
			shutdowner.Shutdown()
		}),
	).Run()
}
