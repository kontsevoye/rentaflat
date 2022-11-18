package flat_subscriber

import "github.com/kontsevoye/rentaflat/internal/flat_parser"

func NewAreaRangeCriteria(from uint, to uint) AreaRangeCriteria {
	return AreaRangeCriteria{from, to}
}

type AreaRangeCriteria struct {
	from uint
	to   uint
}

func (c AreaRangeCriteria) Fits(flat flat_parser.Flat) bool {
	return flat.Area >= c.from && flat.Area <= c.to
}
