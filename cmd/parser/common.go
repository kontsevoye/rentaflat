package parser

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
