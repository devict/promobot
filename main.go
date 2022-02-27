package main

import (
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	devICTSlackWebhook string `envconfig:"DEVICT_SLACK_WEBHOOK" required:"true"`
	devICTTwitter      devICTTwitterConfig
	devICTMeetupURL    string `envconfig:"DEVICT_MEETUP_URL" required:"true"`
}

type devICTTwitterConfig struct {
	AccessToken       string `envconfig:"DEVICT_TW_ACCESS_TOKEN" required:"true"`
	AccessTokenSecret string `envconfig:"DEVICT_TW_ACCESS_TOKEN_SECRET" required:"true"`
	APIKey            string `envconfig:"DEVICT_TW_API_KEY" required:"true"`
	APISecretKey      string `envconfig:"DEVICT_TW_API_SECRET_KEY" required:"true"`
}

func main() {
	var c config
	envconfig.MustProcess("", &c)

	engine.NewEngine(engine.EngineConfig{
		Channels: []channels.Channel{
			channels.Channel(channels.NewSlackChannel("devICT", c.devICTSlackWebhook)),
			channels.Channel(channels.NewTwitterChannel("devICT", channels.TwitterConfig{
				AccessToken:       c.devICTTwitter.AccessToken,
				AccessTokenSecret: c.devICTTwitter.AccessTokenSecret,
				APIKey:            c.devICTTwitter.APIKey,
				APISecretKey:      c.devICTTwitter.APISecretKey,
			})),
		},
		Sources: []sources.Source{
			sources.Source(sources.NewMeetupSource("devICT", c.devICTMeetupURL)),
		},
		Rules: []rules.NotifyRule{
			{
				NumDaysOut: 1,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						return "TODO!"
					},
					"twitter": func(e sources.Event) string {
						return "TODO!"
					},
				},
			},
		},
		SleepTime: 1 * time.Second,
	}).Run()
}
