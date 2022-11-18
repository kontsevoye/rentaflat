package flat_subscriber

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/uuid"
)

type Criteria interface {
	Fits(flat flat_parser.Flat) bool
}

type Subscriber interface {
	AddCriteria(Criteria)
	GetId() uuid.UUID
	Fits(flat_parser.Flat) bool
	Notify(flat_parser.Flat)
}
