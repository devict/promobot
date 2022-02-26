package main

import (
	"fmt"
	"log"

	"github.com/devict/promobot/channels"
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

	notificationRules := []NotificationRule{
		{1, "Today! Join %s at %s for %s\n\nMore info at %s"},
	}

	for _, source := range sources {
		events, err := source.Retrieve()
		if err != nil {
			// TODO: make sure this is the right error logging pattern
			log.Print(fmt.Errorf(
				"failed to retrieve events from %s source %s: %w",
				source.Type(),
				source.Name(),
				err,
			))
			continue
		}

		messages := messagesFromEvents(events, notificationRules)

		for _, message := range messages {
			for _, channel := range channels {
				if err := channel.Send(message); err != nil {
					log.Print(fmt.Errorf(
						"failed to send to %s channel %s: %w",
						channel.Type(),
						channel.Name(),
						err,
					))
				}
			}
		}
	}
}

type NotificationRule struct {
	numDaysOut      int
	messageTemplate string
}

func messagesFromEvents(events []sources.Event, rules []NotificationRule) []string {
	return []string{}
}
