package main

import "time"

type Event struct {
	Name     string
	Source   string
	URL      string
	DateTime time.Time
	Location string
}
