package command

import (
	"github.com/kontsevoye/rentaflat/internal/parser/domain"
	"go.uber.org/zap"
)

func NewParseFlatListHandler(
	parser domain.Parser,
	repository domain.Repository,
	logger *zap.Logger,
) ParseFlatListHandler {
	return ParseFlatListHandler{parser, repository, logger}
}

type ParseFlatList struct {
	LastId string
}

type ParseFlatListHandler struct {
	parser     domain.Parser
	repository domain.Repository
	logger     *zap.Logger
}

func (h ParseFlatListHandler) Handle(cmd ParseFlatList) []error {
	var errs []error
	flats, errsChan := h.parser.Parse(cmd.LastId)

	for flat := range flats {
		err := h.repository.Add(flat)
		if err != nil {
			errs = append(errs, err)
			h.logger.Error("error while saving flat", zap.Error(err))
		}
	}
	for err := range errsChan {
		errs = append(errs, err)
		h.logger.Error(err.Error())
	}

	return errs
}
