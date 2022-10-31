package storage

import "github.com/kontsevoye/rentaflat/cmd/parser"

type Storage interface {
	Store(flats []parser.Flat) int
	Has(id string) bool
	GetAllNew() map[string]parser.Flat
}
