package subscriber

import (
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"go.uber.org/zap"
)

func NewLogSubscriberFactory(logger *zap.Logger, uuidGenerator uuid.Generator) LogSubscriberFactory {
	return LogSubscriberFactory{
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

type LogSubscriberFactory struct {
	logger        *zap.Logger
	uuidGenerator uuid.Generator
}

func (s LogSubscriberFactory) NewLogSubscriber() (LogSubscriber, error) {
	id, err := s.uuidGenerator.UuidV4()
	if err != nil {
		return LogSubscriber{}, err
	}

	return LogSubscriber{
		logger: s.logger,
		id:     id,
	}, nil
}

type LogSubscriber struct {
	id           uuid.UUID
	criteriaList []Criteria
	logger       *zap.Logger
}

func (s *LogSubscriber) GetId() uuid.UUID {
	return s.id
}

func (s *LogSubscriber) AddCriteria(criteria Criteria) {
	s.criteriaList = append(s.criteriaList, criteria)
}

func (s *LogSubscriber) Fits(flat Flat) bool {
	for _, criteria := range s.criteriaList {
		fits := criteria.Fits(flat)
		if !fits {
			return false
		}
	}

	return true
}

func (s *LogSubscriber) Notify(flat Flat) {
	s.logger.Info(
		"got new flat",
		zap.String("uid", s.id.String()),
		zap.String("id", flat.Id().String()),
		zap.Uint("area", flat.Area()),
		zap.Uint("price", flat.Price()),
	)
}
