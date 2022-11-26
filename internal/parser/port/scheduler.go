package port

import (
	"github.com/kontsevoye/rentaflat/internal/parser"
	"github.com/kontsevoye/rentaflat/internal/parser/app"
	"github.com/kontsevoye/rentaflat/internal/parser/app/command"
	"github.com/kontsevoye/rentaflat/internal/parser/app/query"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Scheduler struct {
	app    app.Application
	logger *zap.Logger
	ticker *time.Ticker
	done   chan interface{}
}

func NewScheduler(app app.Application, c *parser.AppConfig, log *zap.Logger) *Scheduler {
	return &Scheduler{
		app,
		log,
		time.NewTicker(c.PollInterval),
		make(chan interface{}),
	}
}

func (s *Scheduler) task(mutex *sync.Mutex) {
	if !mutex.TryLock() {
		s.logger.Warn("lock unavailable, skipping task")
		return
	} else {
		s.logger.Debug("lock acquired")
	}
	defer func() {
		mutex.Unlock()
		s.logger.Debug("lock released")
	}()

	lastId, err := s.app.Queries.GetLastFlatServiceId.Handle(query.GetLastFlatServiceId{})
	if err != nil {
		return
	}
	s.app.Commands.ParseFlatList.Handle(command.ParseFlatList{LastId: lastId})
}

func (s *Scheduler) Run() {
	mutex := &sync.Mutex{}
	go s.task(mutex)
	for {
		select {
		case <-s.ticker.C:
			s.logger.Debug("scheduler tick")
			go s.task(mutex)
		case <-s.done:
			s.logger.Debug("scheduler shutdown")
			return
		}
	}
}

func (s *Scheduler) Shutdown() {
	s.logger.Debug("shutdown call")
	s.done <- struct{}{}
}
