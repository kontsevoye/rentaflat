package parser

type Repository interface {
	Add(Flat) error
	Has(string) (bool, error)
	FindByUrl(string) (Flat, error)
}
