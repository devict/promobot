package channels

type TwitterConfig struct {
	AccessToken       string
	AccessTokenSecret string
	APIKey            string
	APISecretKey      string
}

type TwitterChannel struct {
	name   string
	config TwitterConfig
}

func NewTwitterChannel(name string, config) *TwitterChannel {
	return &TwitterChannel{name: name, config: config}
}

func (c *TwitterChannel) Type() string { return "twitter" }
func (c *TwitterChannel) Name() string { return c.name }

func (c *TwitterChannel) Send(message string) error {
	// TODO: implement twitter sending
	return nil
}
