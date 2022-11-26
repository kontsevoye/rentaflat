package domain

import (
	"errors"
	"fmt"
)

var ErrFlatNotFound = errors.New("flat not found")

func ErrFlatNotFoundWithId(id string) error {
	return fmt.Errorf("%w: %s", ErrFlatNotFound, id)
}

type Repository interface {
	Add(Flat) error
	Has(string) (bool, error)
	FindByUrl(string) (Flat, error)
	FindLatest() (Flat, error)
}
