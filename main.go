package main

import (
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
)

func main() {
	engine.NewEngine(engine.EngineConfig{

		Channels: []channels.Channel{
			channels.Channel(channels.NewSlackChannel("devICT", "")),
			channels.Channel(channels.NewTwitterChannel("devICT", "")),
		},

		Sources: []sources.Source{
			sources.Source(sources.NewMeetupSource("devICT", "")),
		},

		Rules: []rules.NotifyRule{
			{
				NumDaysOut: 1,
				ChannelTemplates: map[string]rules.EventTemplate{
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
