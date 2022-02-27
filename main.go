package main

import (
	"fmt"
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DevICTSlackWebhook string `envconfig:"DEVICT_SLACK_WEBHOOK" required:"true"`
	DevICTTwitter      DevICTTwitterConfig
	DevICTMeetupURL    string `envconfig:"DEVICT_MEETUP_URL" required:"true"`
	OzSecMeetupURL     string `envconfig:"OZSEC_MEETUP_URL" required:"true"`
}

type DevICTTwitterConfig struct {
	AccessToken       string `envconfig:"DEVICT_TW_ACCESS_TOKEN" required:"true"`
	AccessTokenSecret string `envconfig:"DEVICT_TW_ACCESS_TOKEN_SECRET" required:"true"`
	APIKey            string `envconfig:"DEVICT_TW_API_KEY" required:"true"`
	APISecretKey      string `envconfig:"DEVICT_TW_API_SECRET_KEY" required:"true"`
}

func main() {
	var c config
	if err := envconfig.Process("", &c); err != nil {
		panic(err)
	}

	engine.NewEngine(engine.EngineConfig{
		Sources: []sources.Source{
			sources.Source(sources.NewMeetupSource("devICT", c.DevICTMeetupURL)),
			sources.Source(sources.NewMeetupSource("OzSec", c.OzSecMeetupURL)),
		},
		Channels: []channels.Channel{
			channels.Channel(channels.NewSlackChannel("devICT", c.DevICTSlackWebhook)),
			channels.Channel(channels.NewTwitterChannel("devICT", channels.TwitterConfig{
				AccessToken:       c.DevICTTwitter.AccessToken,
				AccessTokenSecret: c.DevICTTwitter.AccessTokenSecret,
				APIKey:            c.DevICTTwitter.APIKey,
				APISecretKey:      c.DevICTTwitter.APISecretKey,
			})),
		},
		Rules: []rules.NotifyRule{
			{
				NumDaysOut: 1,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						return fmt.Sprintf("*Tomorrow!* %s is hosting <%s|%s> at %s", e.Source, e.URL, e.Name, e.Location)
					},
					"twitter": func(e sources.Event) string {
						return fmt.Sprintf("Tomorrow! %s is hosting %s at %s\n\nMore info here: %s", e.Source, e.Name, e.Location, e.URL)
					},
				},
			},
		},
		SleepTime: 1 * time.Hour,
	}).RunOnce()
	// }).Run()
}
