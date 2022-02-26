package channels

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
