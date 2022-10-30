package main

import (
	"fmt"
	"github.com/kontsevoye/rentaflat/cmd/logger"
	ssge "github.com/kontsevoye/rentaflat/cmd/parser/ss.ge"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(ssge.NewParser),
		fx.Provide(logger.NewLogger),
		fx.Invoke(func(p *ssge.Parser, log *zap.SugaredLogger) {
			url := "https://ss.ge/en/real-estate/l/For-Rent?Sort.SortExpression=%22OrderDate%22%20DESC&RealEstateDealTypeId=1&CommercialRealEstateType=&PriceType=false&CurrencyId=1&Context.Request.Query[Query]=&WithImageOnly=true"
			flats, err := p.Parse(url, 10)
			if err != nil {
				log.Fatal(err)
			}
			for _, flat := range flats {
				fmt.Printf("%s %dm^2, $%d\n", flat.Title, flat.Area, flat.Price)
			}
		}),
	).Run()
}
