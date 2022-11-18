package flat_parser

import "time"

type Flat struct {
	Id          string
	Url         string
	PhotoUrls   []string
	Title       string
	Description string
	Area        uint
	Rooms       uint
	Floor       uint
	Price       uint
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
