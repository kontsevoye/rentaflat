package flat_storage

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/flat_subscriber"
	"github.com/kontsevoye/rentaflat/internal/uuid"
	"go.uber.org/zap"
)

func NewInMemoryStorage(logger *zap.Logger, generator uuid.Generator) *InMemoryStorage {
	return &InMemoryStorage{
		logger,
		make(map[string]flat_parser.Flat),
		make(map[uuid.UUID]flat_subscriber.Subscriber),
		generator,
	}
}

type InMemoryStorage struct {
	logger        *zap.Logger
	flats         map[string]flat_parser.Flat
	subscribers   map[uuid.UUID]flat_subscriber.Subscriber
	uuidGenerator uuid.Generator
}

func (s *InMemoryStorage) Store(flat flat_parser.Flat) {
	if s.Has(flat.Id) {
		return
	}
	s.flats[flat.Id] = flat

	for _, subscriber := range s.subscribers {
		if subscriber.Fits(flat) {
			go subscriber.Notify(flat)
		}
	}
}

func (s *InMemoryStorage) Has(id string) bool {
	_, exist := s.flats[id]

	return exist
}

func (s *InMemoryStorage) Subscribe(subscriber flat_subscriber.Subscriber) {
	s.subscribers[subscriber.GetId()] = subscriber
}

func (s *InMemoryStorage) Unsubscribe(id uuid.UUID) {
	s.logger.Debug("unsubscribe from in-memory flat_storage " + id.String())
	delete(s.subscribers, id)
}

func (s *InMemoryStorage) Count() int {
	return len(s.flats)
}
