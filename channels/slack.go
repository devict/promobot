package channels

type SlackChannel struct {
	name string
	url  string
}

func NewSlackChannel(name, url string) *SlackChannel {
	return &SlackChannel{name: name, url: url}
}

func (c *SlackChannel) Send(message string) error {
	// TODO: implement slack sending
	return nil
}

func (c *SlackChannel) Type() string {
	return "slack"
}

func (c *SlackChannel) Name() string {
	return c.name
}
