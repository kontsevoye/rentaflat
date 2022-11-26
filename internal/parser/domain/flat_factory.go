package domain

import (
	"errors"
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"net/url"
	"time"
)

func NewFlatFactory(generator uuid.Generator) FlatFactory {
	return FlatFactory{generator: generator}
}

type FlatFactory struct {
	generator uuid.Generator
}

var (
	ErrEmptyServiceId = errors.New("empty service id")
	ErrEmptyTitle     = errors.New("empty title")
	ErrEmptyArea      = errors.New("empty area")
	ErrEmptyRooms     = errors.New("empty rooms")
	ErrEmptyFloor     = errors.New("empty floor")
	ErrEmptyPrice     = errors.New("empty price")
	ErrEmptyPhone     = errors.New("empty phone")
)

func (f FlatFactory) NewFlat(
	serviceId string,
	url *url.URL,
	photoUrls []*url.URL,
	title string,
	description string,
	area uint,
	rooms uint,
	floor uint,
	price uint,
	contactName string,
	phone string,
	isAgency bool,
	publishedAt time.Time,
) (*Flat, error) {
	id, err := f.generator.UuidV4()
	if serviceId == "" {
		return nil, ErrEmptyServiceId
	}
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if area == 0 {
		return nil, ErrEmptyArea
	}
	if rooms == 0 {
		return nil, ErrEmptyRooms
	}
	if floor == 0 {
		return nil, ErrEmptyFloor
	}
	if price == 0 {
		return nil, ErrEmptyPrice
	}
	if phone == "" {
		return nil, ErrEmptyPhone
	}

	return &Flat{
		id:          id,
		serviceId:   serviceId,
		url:         url,
		photoUrls:   photoUrls,
		title:       title,
		description: description,
		area:        area,
		rooms:       rooms,
		floor:       floor,
		price:       price,
		contactName: contactName,
		phone:       phone,
		isAgency:    isAgency,
		publishedAt: publishedAt,
		createdAt:   time.Now(),
	}, err
}

func (f FlatFactory) LoadFlat(
	rawId string,
	serviceId string,
	rawUrl string,
	rawPhotoUrls []string,
	title string,
	description string,
	area uint,
	rooms uint,
	floor uint,
	price uint,
	contactName string,
	phone string,
	isAgency bool,
	publishedAt time.Time,
	createdAt time.Time,
) (Flat, error) {
	parsedId, err := f.generator.FromString(rawId)
	if err != nil {
		return Flat{}, err
	}
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return Flat{}, err
	}
	var photoUrls []*url.URL
	for _, rawUrl := range rawPhotoUrls {
		parsedPhotoUrl, err := url.Parse(rawUrl)
		if err != nil {
			return Flat{}, err
		}
		photoUrls = append(photoUrls, parsedPhotoUrl)
	}

	return Flat{
		id:          parsedId,
		serviceId:   serviceId,
		url:         parsedUrl,
		photoUrls:   photoUrls,
		title:       title,
		description: description,
		area:        area,
		rooms:       rooms,
		floor:       floor,
		price:       price,
		contactName: contactName,
		phone:       phone,
		isAgency:    isAgency,
		publishedAt: publishedAt,
		createdAt:   createdAt,
	}, nil
}
