package flat_scheduler

import (
	"github.com/kontsevoye/rentaflat/internal/config"
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/flat_storage"
	"go.uber.org/zap"
	"time"
)

type Scheduler struct {
	logger  *zap.Logger
	parser  flat_parser.Parser
	storage flat_storage.Storage
	ticker  *time.Ticker
	done    chan interface{}
}

func NewScheduler(p flat_parser.Parser, log *zap.Logger, s flat_storage.Storage, c *config.AppConfig) *Scheduler {
	return &Scheduler{
		log,
		p,
		s,
		time.NewTicker(c.PollInterval),
		make(chan interface{}),
	}
}

func (s *Scheduler) Run() {
	task := func() {
		url := "https://ss.ge/en/real-estate/l/For-Rent?Sort.SortExpression=%22OrderDate%22%20DESC&RealEstateDealTypeId=1&CommercialRealEstateType=&PriceType=false&CurrencyId=1&Context.Request.Query[Query]=&WithImageOnly=true"
		flats, errs := s.parser.Parse(url)
		for flat := range flats {
			s.storage.Store(flat)
		}
		for err := range errs {
			s.logger.Error(err.Error())
		}
	}
	go task()
	for {
		select {
		case <-s.ticker.C:
			s.logger.Debug("Scheduler tick")
			task()
		case <-s.done:
			s.logger.Debug("Scheduler shutdown")
			s.done <- struct{}{}
			return
		}
	}
}

func (s *Scheduler) Shutdown() {
	s.logger.Debug("Shutdown call")
	s.done <- struct{}{}
	<-s.done
}
