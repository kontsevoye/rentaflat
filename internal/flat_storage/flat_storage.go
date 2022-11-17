package flat_storage

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
	"github.com/kontsevoye/rentaflat/internal/uuid"
)

type Storage interface {
	Store(flat flat_parser.Flat)
	Has(id string) bool
	Subscribe(func(flat_parser.Flat)) uuid.UUID
	Unsubscribe(uuid.UUID)
	Count() int
}
