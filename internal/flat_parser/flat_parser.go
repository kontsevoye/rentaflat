package flat_parser

import "time"

type Flat struct {
	Id          string
	Url         string
	PhotoUrls   []string
	Title       string
	Description string
	Area        int
	Rooms       int
	Floor       int
	Price       int
	ContactName string
	Phone       string
	IsAgency    bool
	PublishedAt time.Time
}

type Request struct {
	LastId string
}

type Parser interface {
	Parse(Request) (<-chan Flat, <-chan error)
}
