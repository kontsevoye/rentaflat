package uuid

import "github.com/google/uuid"

func NewGoogleGenerator() GoogleGenerator {
	return GoogleGenerator{}
}

type GoogleUuid struct {
	uuid.UUID
}

func (u GoogleUuid) String() string {
	return u.UUID.String()
}

type GoogleGenerator struct {
}

func (u GoogleGenerator) UuidV4() (UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return GoogleUuid{}, err
	}

	return GoogleUuid{id}, nil
}
