package uuid

type UUID interface {
	String() string
}

type Generator interface {
	UuidV4() (UUID, error)
	FromString(string) (UUID, error)
}
