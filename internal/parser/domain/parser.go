package domain

type Parser interface {
	Parse(lastServiceId string) (<-chan Flat, <-chan error)
}
