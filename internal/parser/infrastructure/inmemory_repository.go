package infrastructure

import (
	"errors"
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"github.com/kontsevoye/rentaflat/internal/parser/domain"
	"go.uber.org/zap"
)

func NewInMemoryRepository(logger *zap.Logger, generator uuid.Generator) *InMemoryRepository {
	return &InMemoryRepository{
		logger,
		make(map[string]domain.Flat),
		generator,
	}
}

type InMemoryRepository struct {
	logger        *zap.Logger
	flats         map[string]domain.Flat
	uuidGenerator uuid.Generator
}

func (s *InMemoryRepository) Add(flat domain.Flat) error {
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

func (s *InMemoryRepository) FindByUrl(id string) (domain.Flat, error) {
	flat, exist := s.flats[id]
	if !exist {
		return domain.Flat{}, domain.ErrFlatNotFoundWithId(id)
	}

	return flat, nil
}

func (s *InMemoryRepository) FindLatest() (domain.Flat, error) {
	var latest domain.Flat
	for _, flat := range s.flats {
		if latest.PublishedAt().Before(flat.PublishedAt()) || latest.PublishedAt().Equal(flat.PublishedAt()) {
			latest = flat
		}
	}
	if latest.Url().String() == "" {
		return domain.Flat{}, domain.ErrFlatNotFound
	}

	return latest, nil
}
