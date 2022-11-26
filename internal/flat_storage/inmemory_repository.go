package flat_storage

import (
	"errors"
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/uuid"
	"go.uber.org/zap"
)

func NewInMemoryRepository(logger *zap.Logger, generator uuid.Generator) *InMemoryRepository {
	return &InMemoryRepository{
		logger,
		make(map[string]flat_parser.Flat),
		generator,
	}
}

type InMemoryRepository struct {
	logger        *zap.Logger
	flats         map[string]flat_parser.Flat
	uuidGenerator uuid.Generator
}

func (s *InMemoryRepository) Add(flat flat_parser.Flat) error {
	has, err := s.Has(flat.Url().String())
	if err != nil {
		return err
	}
	if has {
		return errors.New("flat already exists")
	}
	s.flats[flat.Url().String()] = flat

	return nil
}

func (s *InMemoryRepository) Has(id string) (bool, error) {
	_, exist := s.flats[id]

	return exist, nil
}

func (s *InMemoryRepository) FindByUrl(id string) (flat_parser.Flat, error) {
	flat, exist := s.flats[id]
	if !exist {
		return flat_parser.Flat{}, errors.New("flat doesnt exist")
	}

	return flat, nil
}
