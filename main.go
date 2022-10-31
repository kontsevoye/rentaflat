package main

import (
	"fmt"
	"github.com/kontsevoye/rentaflat/cmd/logger"
	"github.com/kontsevoye/rentaflat/cmd/parser"
	ssge "github.com/kontsevoye/rentaflat/cmd/parser/ss.ge"
	"github.com/kontsevoye/rentaflat/cmd/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func main() {
	fx.New(
		fx.Provide(
			logger.NewLogger,
			fx.Annotate(
				storage.NewInMemoryStorage,
				fx.As(new(storage.Storage)),
			),
			fx.Annotate(
				ssge.NewParser,
				fx.As(new(parser.Parser)),
			),
		),
		fx.Invoke(func(p parser.Parser, log *zap.SugaredLogger, storage storage.Storage) {
			ticker := time.NewTicker(1 * time.Minute)
			task := func() {
				url := "https://ss.ge/en/real-estate/l/For-Rent?Sort.SortExpression=%22OrderDate%22%20DESC&RealEstateDealTypeId=1&CommercialRealEstateType=&PriceType=false&CurrencyId=1&Context.Request.Query[Query]=&WithImageOnly=true"
				flats, err := p.Parse(url, 10)
				if err != nil {
					log.Fatal(err)
				}
				storage.Store(flats)
				storedFlats := storage.GetAllNew()
				for _, flat := range storedFlats {
					fmt.Printf("#%s %s %dm^2, $%d %s\n", flat.Id, flat.Title, flat.Area, flat.Price, flat.PublishedAt.Format(time.RFC3339))
				}
				fmt.Printf("There are %d new flats\n", len(storedFlats))
			}
			go func() {
				for t := range ticker.C {
					fmt.Println("Tick at", t)
					task()
				}
			}()
			task()
		}),
	).Run()
}
