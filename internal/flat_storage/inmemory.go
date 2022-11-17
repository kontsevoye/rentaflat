package flat_storage

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/uuid"
	"go.uber.org/zap"
)

func NewInMemoryStorage(logger *zap.Logger, g uuid.Generator) *InMemoryStorage {
	return &InMemoryStorage{
		logger,
		make(map[string]flat_parser.Flat),
		make(map[uuid.UUID]func(flat_parser.Flat)),
		g,
	}
}

type InMemoryStorage struct {
	logger        *zap.Logger
	flats         map[string]flat_parser.Flat
	subscribers   map[uuid.UUID]func(flat_parser.Flat)
	uuidGenerator uuid.Generator
}

func (s *InMemoryStorage) Store(flat flat_parser.Flat) {
	if s.Has(flat.Id) {
		return
	}
	s.flats[flat.Id] = flat

	for _, subscriber := range s.subscribers {
		go subscriber(flat)
	}
}

func (s *InMemoryStorage) Has(id string) bool {
	_, exist := s.flats[id]

	return exist
}

func (s *InMemoryStorage) Subscribe(subscriber func(flat_parser.Flat)) uuid.UUID {
	id, err := s.uuidGenerator.UuidV4()
	if err != nil {
		s.logger.Error(err.Error())
	}
	s.subscribers[id] = subscriber

	return id
}

func (s *InMemoryStorage) Unsubscribe(id uuid.UUID) {
	s.logger.Debug("unsubscribe from in-memory flat_storage " + id.String())
	delete(s.subscribers, id)
}

func (s *InMemoryStorage) Count() int {
	return len(s.flats)
}
