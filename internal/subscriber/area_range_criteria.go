package subscriber

func NewAreaRangeCriteria(from uint, to uint) AreaRangeCriteria {
	return AreaRangeCriteria{from, to}
}

type AreaRangeCriteria struct {
	from uint
	to   uint
}

func (c AreaRangeCriteria) Fits(flat Flat) bool {
	return flat.Area() >= c.from && flat.Area() <= c.to
}
