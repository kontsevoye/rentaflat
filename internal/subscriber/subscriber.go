package subscriber

import (
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
)

type Subscriber interface {
	AddCriteria(Criteria)
	GetId() uuid.UUID
	Fits(Flat) bool
	Notify(Flat)
}
