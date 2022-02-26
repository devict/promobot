package main

type Channel interface {
	Send(string) error
	Type() string
	Name() string
}

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

type TwitterConfig struct {
	token string
}

type TwitterChannel struct {
	name   string
	config TwitterConfig
}

func NewTwitterChannel(name, token string) *TwitterChannel {
	return &TwitterChannel{name: name, config: TwitterConfig{token: token}}
}

func (c *TwitterChannel) Send(message string) error {
	// TODO: implement twitter sending
	return nil
}

func (c *TwitterChannel) Type() string {
	return "twitter"
}

func (c *TwitterChannel) Name() string {
	return c.name
}
