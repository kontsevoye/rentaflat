package flat_subscriber

import "github.com/kontsevoye/rentaflat/internal/flat_parser"

func NewPriceRangeCriteria(from uint, to uint) PriceRangeCriteria {
	return PriceRangeCriteria{from, to}
}

type PriceRangeCriteria struct {
	from uint
	to   uint
}

func (c PriceRangeCriteria) Fits(flat flat_parser.Flat) bool {
	return flat.Price() >= c.from && flat.Price() <= c.to
}
