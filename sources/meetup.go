package sources

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

func (c *MeetupSource) Name() string {
	return c.name
}

func (c *MeetupSource) Type() string {
	return "meetup"
}

func (c *MeetupSource) Retrieve() ([]Event, error) {
	// TODO
	return []Event{}, nil
}
