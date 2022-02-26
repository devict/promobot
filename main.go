package main

import (
	"text/template"
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
)

func main() {
	sources := []sources.Source{
		sources.Source(sources.NewMeetupSource("devICT", "")),
	}

	channels := []channels.Channel{
		channels.Channel(channels.NewSlackChannel("devICT", "")),
		channels.Channel(channels.NewTwitterChannel("devICT", "")),
	}

	notificationRules := []rules.NotificationRule{
		{
			NumDaysOut: 1,
			ChannelTemplates: map[string]*template.Template{
				"slack": template.Must(
					template.New("slack").Parse("Today! Join %s at %s for %s\n\nMore info at %s"),
				),
				"twitter": template.Must(
					template.New("twitter").Parse("Today! Join %s at %s for %s\n\nMore info at %s"),
				),
			},
		},
	}

	config := engine.EngineConfig{
		Channels: channels,
		Sources:  sources,
		Rules:    notificationRules,
	}

	sleepDuration := time.Second

	engine.NewEngine(config, sleepDuration).Run()
}
