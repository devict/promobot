package channels

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterConfig struct {
	AccessToken       string
	AccessTokenSecret string
	APIKey            string
	APISecretKey      string
}

type TwitterChannel struct {
	name   string
	client *twitter.Client
}

func NewTwitterChannel(name string, config TwitterConfig) *TwitterChannel {
	oa := oauth1.NewConfig(config.APIKey, config.APISecretKey)
	token := oauth1.NewToken(
		config.AccessToken,
		config.AccessTokenSecret,
	)
	httpClient := oa.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return &TwitterChannel{name: name, client: client}
}

func (c *TwitterChannel) Type() string { return "twitter" }
func (c *TwitterChannel) Name() string { return c.name }

func (c *TwitterChannel) Send(message string) error {
	// TODO: check for failures in the resp object?
	_, _, err := c.client.Statuses.Update(message, nil)
	if err != nil {
		return fmt.Errorf("failed to post to twitter: %w", err)
	}

	return nil
}
