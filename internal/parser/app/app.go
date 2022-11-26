package app

import (
	"github.com/kontsevoye/rentaflat/internal/parser/app/command"
	"github.com/kontsevoye/rentaflat/internal/parser/app/query"
)

func NewApplication(pflh command.ParseFlatListHandler, glfsih query.GetLastFlatServiceIdHandler) Application {
	return Application{
		Commands: Commands{
			ParseFlatList: pflh,
		},
		Queries: Queries{
			GetLastFlatServiceId: glfsih,
		},
	}
}

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ParseFlatList command.ParseFlatListHandler
}

type Queries struct {
	GetLastFlatServiceId query.GetLastFlatServiceIdHandler
}
