package storage

import (
	"github.com/kontsevoye/rentaflat/cmd/parser"
)

type wrappedFlat struct {
	parser.Flat
	isNew bool
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		make(map[string]wrappedFlat),
	}
}

type InMemoryStorage struct {
	flats map[string]wrappedFlat
}

func (s *InMemoryStorage) Store(flats []parser.Flat) int {
	newFlatsCount := 0
	for _, flat := range flats {
		if s.Has(flat.Id) {
			continue
		}
		s.flats[flat.Id] = wrappedFlat{
			Flat:  flat,
			isNew: true,
		}
		newFlatsCount++
	}

	return newFlatsCount
}

func (s *InMemoryStorage) Has(id string) bool {
	_, exist := s.flats[id]

	return exist
}

func (s *InMemoryStorage) GetAllNew() map[string]parser.Flat {
	flats := make(map[string]parser.Flat)
	for id, wFlat := range s.flats {
		if wFlat.isNew {
			flats[id] = wFlat.Flat
			wFlat.isNew = false
			s.flats[id] = wFlat
		}
	}

	return flats
}
