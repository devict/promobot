package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MeetupSource struct {
	name string
	urlName  string

}

func NewMeetupSource(name, urlName string) *MeetupSource {
	return &MeetupSource{
		name: name,
		urlName: urlName,
	}
}

func (c *MeetupSource) Name() string { return c.name }
func (c *MeetupSource) Type() string { return "meetup" }
func (c *MeetupSource) JsonUrl() string {
	return fmt.Sprintf("https://api.meetup.com/2/events?&sign=true&photo-host=public&group_urlname=%s&limited_events=false&fields=series&status=upcoming&page=20", c.urlName)
}
func (c *MeetupSource) HtmlUrl() string {
	return fmt.Sprintf("https://meetup.com/%s", c.urlName)
}

func (c *MeetupSource) Retrieve(loc *time.Location) ([]Event, error) {
	resp, err := http.Get(c.JsonUrl())
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
