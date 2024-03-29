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
	DebugMode bool `envconfig:"DEBUG_MODE"`

	// Channels
	DevICTSlackWebhook string `envconfig:"DEVICT_SLACK_WEBHOOK" required:"true"`
	DevICTTwitter      DevICTTwitterConfig

	// Sources
	DevICTMeetupURL string `envconfig:"DEVICT_MEETUP_URL" required:"true"`
	OzSecMeetupURL  string `envconfig:"OZSEC_MEETUP_URL" required:"true"`
	WTFMeetupURL    string `envconfig:"WTF_MEETUP_URL" required:"true"`
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

	loc, _ := time.LoadLocation("America/Chicago")

	e := engine.NewEngine(engine.EngineConfig{
		Sources: []sources.Source{
			sources.Source(sources.NewMeetupSource("devICT", c.DevICTMeetupURL)),
			sources.Source(sources.NewMeetupSource("OzSec", c.OzSecMeetupURL)),
			sources.Source(sources.NewMeetupSource("Wichita Technology Forum", c.WTFMeetupURL)),
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
				NumDaysOut: 10,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						t := e.DateTime.Format("Mon 01/02 @ 03:04 PM")
						return fmt.Sprintf("*%s*, %s is hosting <%s|%s> at %s", t, e.Source, e.URL, e.Name, e.Location)
					},
					// "twitter": func(e sources.Event) string {
					// 	t := e.DateTime.Format("Mon, 01/02 at 03:04 PM")
					// 	return fmt.Sprintf("Join %s for %s! %s\n\nRSVP at %s", e.Source, e.Name, t, e.URL)
					// },
				},
			},
			{
				NumDaysOut: 4,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						t := e.DateTime.Format("Monday @ 03:04 PM")
						return fmt.Sprintf("*%s*, %s is hosting <%s|%s> at %s", t, e.Source, e.URL, e.Name, e.Location)
					},
					// "twitter": func(e sources.Event) string {
					// 	t := e.DateTime.Format("Monday at 03:04 PM")
					// 	return fmt.Sprintf("Join %s for %s! %s\n\nMore info at %s", e.Source, e.Name, t, e.URL)
					// },
				},
			},
			{
				NumDaysOut: 1,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						return fmt.Sprintf("*Tomorrow!* %s is hosting <%s|%s> at %s", e.Source, e.URL, e.Name, e.Location)
					},
					// "twitter": func(e sources.Event) string {
					// 	t := e.DateTime.Format("03:04 PM")
					// 	return fmt.Sprintf("Tomorrow! Join %s at %s for %s\n\nMore info here: %s", e.Source, t, e.Name, e.URL)
					// },
				},
			},
			{
				NumDaysOut: 0,
				ChannelTemplates: map[string]rules.MsgFunc{
					"slack": func(e sources.Event) string {
						t := e.DateTime.Format("03:04 PM")
						return fmt.Sprintf("*Today!* %s is hosting <%s|%s> at %s, %s", e.Source, e.URL, e.Name, e.Location, t)
					},
				},
			},
		},
		Location:  loc,
		DebugMode: c.DebugMode,
	})

	e.RunOnce()
}
