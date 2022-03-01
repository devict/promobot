package sources

import "time"

type Source interface {
	Name() string
	Type() string
	Retrieve(*time.Location) ([]Event, error)
}

type Event struct {
	Name     string
	Source   string
	URL      string
	DateTime time.Time
	Location string
}
