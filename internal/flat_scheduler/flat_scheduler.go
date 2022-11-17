package flat_scheduler

import (
	"github.com/kontsevoye/rentaflat/internal/config"
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/flat_storage"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Scheduler struct {
	logger  *zap.Logger
	parser  flat_parser.Parser
	storage flat_storage.Storage
	ticker  *time.Ticker
	done    chan interface{}
	lastId  string
}

func NewScheduler(p flat_parser.Parser, log *zap.Logger, s flat_storage.Storage, c *config.AppConfig) *Scheduler {
	return &Scheduler{
		log,
		p,
		s,
		time.NewTicker(c.PollInterval),
		make(chan interface{}),
		"0",
	}
}

func (s *Scheduler) Run() {
	mutex := &sync.Mutex{}
	task := func(mutex *sync.Mutex) {
		if !mutex.TryLock() {
			s.logger.Warn("lock unavailable, skipping task")
			return
		} else {
			s.logger.Debug("lock acquired")
		}
		flats, errs := s.parser.Parse(flat_parser.Request{LastId: s.lastId})

		for flat := range flats {
			s.storage.Store(flat)
			if flat.Id > s.lastId {
				s.lastId = flat.Id
			}
		}
		for err := range errs {
			s.logger.Error(err.Error())
		}
		mutex.Unlock()
		s.logger.Debug("lock released")
	}
	go task(mutex)
	for {
		select {
		case <-s.ticker.C:
			s.logger.Debug("scheduler tick")
			task(mutex)
		case <-s.done:
			s.logger.Debug("scheduler shutdown")
			s.done <- struct{}{}
			return
		}
	}
}

func (s *Scheduler) Shutdown() {
	s.logger.Debug("shutdown call")
	s.done <- struct{}{}
	<-s.done
}
