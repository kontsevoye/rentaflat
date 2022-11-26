package query

import (
	"errors"
	"github.com/kontsevoye/rentaflat/internal/parser/domain"
	"go.uber.org/zap"
	"time"
)

func NewGetLastFlatServiceIdHandler(repository domain.Repository, logger *zap.Logger) GetLastFlatServiceIdHandler {
	return GetLastFlatServiceIdHandler{repository, logger}
}

type GetLastFlatServiceId struct {
}

type GetLastFlatServiceIdHandler struct {
	repository domain.Repository
	logger     *zap.Logger
}

func (h GetLastFlatServiceIdHandler) Handle(query GetLastFlatServiceId) (string, error) {
	flat, err := h.repository.FindLatest()
	if errors.Is(err, domain.ErrFlatNotFound) {
		return "0", nil
	}
	if err != nil {
		return "", err
	}

	timeGap := time.Now().Sub(flat.CreatedAt())
	if timeGap.Minutes() > 15 {
		return "0", nil
	}

	return flat.ServiceId(), nil
}
