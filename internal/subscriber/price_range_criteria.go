package subscriber

func NewPriceRangeCriteria(from uint, to uint) PriceRangeCriteria {
	return PriceRangeCriteria{from, to}
}

type PriceRangeCriteria struct {
	from uint
	to   uint
}

func (c PriceRangeCriteria) Fits(flat Flat) bool {
	return flat.Price() >= c.from && flat.Price() <= c.to
}
