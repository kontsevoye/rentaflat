package subscriber

import "github.com/kontsevoye/rentaflat/internal/common/uuid"

type Flat struct {
	id    uuid.UUID
	area  uint
	price uint
}

func (f *Flat) Id() uuid.UUID {
	return f.id
}

func (f *Flat) Area() uint {
	return f.area
}

func (f *Flat) Price() uint {
	return f.price
}
