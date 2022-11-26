package parser

type Request struct {
	LastId string
}

type Parser interface {
	Parse(Request) (<-chan Flat, <-chan error)
}
