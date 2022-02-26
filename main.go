package main

import (
	"fmt"
	"log"
)

func main() {
	sources := []Source{
		Source(NewMeetupSource("devICT", "")),
	}

	channels := []Channel{
		Channel(NewSlackChannel("devICT", "")),
		Channel(NewTwitterChannel("devICT", "")),
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

func messagesFromEvents(events []Event, rules []NotificationRule) []string {
	return []string{}
}
