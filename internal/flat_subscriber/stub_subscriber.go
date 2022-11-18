package flat_subscriber

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/uuid"
	"go.uber.org/zap"
)

func NewStubSubscriberFactory(logger *zap.Logger, uuidGenerator uuid.Generator) StubSubscriberFactory {
	return StubSubscriberFactory{
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

type StubSubscriberFactory struct {
	logger        *zap.Logger
	uuidGenerator uuid.Generator
}

func (s StubSubscriberFactory) NewStubSubscriber() (StubSubscriber, error) {
	id, err := s.uuidGenerator.UuidV4()
	if err != nil {
		return StubSubscriber{}, err
	}

	return StubSubscriber{
		logger: s.logger,
		id:     id,
	}, nil
}

type StubSubscriber struct {
	id           uuid.UUID
	criteriaList []Criteria
	logger       *zap.Logger
}

func (s *StubSubscriber) GetId() uuid.UUID {
	return s.id
}

func (s *StubSubscriber) AddCriteria(criteria Criteria) {
	s.criteriaList = append(s.criteriaList, criteria)
}

func (s *StubSubscriber) Fits(flat flat_parser.Flat) bool {
	for _, criteria := range s.criteriaList {
		fits := criteria.Fits(flat)
		if !fits {
			return false
		}
	}

	return true
}

func (s *StubSubscriber) Notify(flat flat_parser.Flat) {
	s.logger.Info(
		"got new flat",
		zap.String("uid", s.id.String()),
		zap.String("id", flat.Id),
		zap.String("title", flat.Title),
		zap.Uint("area", flat.Area),
		zap.Uint("price", flat.Price),
		zap.Time("publishedAt", flat.PublishedAt),
	)
}
