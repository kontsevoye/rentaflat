package flat_storage

import (
	"github.com/kontsevoye/rentaflat/internal/flat_parser"
)

type Repository interface {
	Add(flat_parser.Flat) error
	Has(string) (bool, error)
	FindByUrl(string) (flat_parser.Flat, error)
}
