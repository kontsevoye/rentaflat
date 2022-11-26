package parser

import (
	"github.com/kontsevoye/rentaflat/internal/common/uuid"
	"net/url"
	"time"
)

type Flat struct {
	id          uuid.UUID
	serviceId   string
	url         *url.URL
	photoUrls   []*url.URL
	title       string
	description string
	area        uint
	rooms       uint
	floor       uint
	price       uint
	contactName string
	phone       string
	isAgency    bool
	publishedAt time.Time
	createdAt   time.Time
}

func (f *Flat) Id() uuid.UUID {
	return f.id
}

func (f *Flat) ServiceId() string {
	return f.serviceId
}

func (f *Flat) Url() *url.URL {
	return f.url
}

func (f *Flat) PhotoUrls() []*url.URL {
	return f.photoUrls
}

func (f *Flat) PhotoUrlsAsStrings() []string {
	var urls []string
	for _, photoUrl := range f.photoUrls {
		urls = append(urls, photoUrl.String())
	}

	return urls
}

func (f *Flat) Title() string {
	return f.title
}

func (f *Flat) Description() string {
	return f.description
}

func (f *Flat) Area() uint {
	return f.area
}

func (f *Flat) Rooms() uint {
	return f.rooms
}

func (f *Flat) Floor() uint {
	return f.floor
}

func (f *Flat) Price() uint {
	return f.price
}

func (f *Flat) ContactName() string {
	return f.contactName
}

func (f *Flat) Phone() string {
	return f.phone
}

func (f *Flat) IsAgency() bool {
	return f.isAgency
}

func (f *Flat) PublishedAt() time.Time {
	return f.publishedAt
}

func (f *Flat) CreatedAt() time.Time {
	return f.createdAt
}
