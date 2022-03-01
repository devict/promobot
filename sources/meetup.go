package sources

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type MeetupSource struct {
	name string
	url  string
}

func NewMeetupSource(name, url string) *MeetupSource {
	return &MeetupSource{
		name: name,
		url:  url,
	}
}

func (c *MeetupSource) Name() string { return c.name }
func (c *MeetupSource) Type() string { return "meetup" }

func (c *MeetupSource) Retrieve(loc *time.Location) ([]Event, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return []Event{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Event{}, err
	}

	var meetupResp meetupResponse
	if err = json.Unmarshal(respBytes, &meetupResp); err != nil {
		return []Event{}, err
	}

	events := make([]Event, 0)
	for _, evt := range meetupResp.Results {
		if haveNextInSeries(events, evt.Name) {
			continue
		}
		events = append(events, Event{
			Name:     evt.Name,
			Source:   c.name,
			URL:      evt.URL,
			Location: evt.Venue.Name,
			DateTime: time.UnixMilli(evt.Time).In(loc),
		})
	}

	return events, nil
}

func haveNextInSeries(events []Event, eventName string) bool {
	for _, e := range events {
		if e.Name == eventName {
			return true
		}
	}
	return false
}

type meetupResponse struct {
	Results []meetupEvent `json:"results"`
}

type meetupEvent struct {
	Name  string      `json:"name"`
	Time  int64       `json:"time"`
	URL   string      `json:"event_url"`
	Venue meetupVenue `json:"venue"`
}

type meetupVenue struct {
	Name string `json:"name"`
}
