package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
)

type EngineConfig struct {
	Channels []channels.Channel
	Sources  []sources.Source
	Rules    []rules.NotifyRule
	RunAt    RunAt
	Location *time.Location
}

type RunAt struct {
	Hour   int
	Minute int
}

type Engine struct {
	config EngineConfig
}

func NewEngine(config EngineConfig) *Engine {
	if config.Location == nil {
		config.Location = time.UTC
	}
	return &Engine{config}
}

func (e *Engine) Run() {
	for {
		if e.ShouldRun(time.Now()) {
			e.RunOnce()
		}
		time.Sleep(time.Minute)
	}
}

func (e *Engine) ShouldRun(now time.Time) bool {
	now = now.Round(time.Minute)
	return now.Hour() == e.config.RunAt.Hour && now.Minute() == e.config.RunAt.Minute
}

func (e *Engine) RunOnce() {
	for _, source := range e.config.Sources {
		events, err := source.Retrieve()
		if err != nil {
			// TODO: make sure this is the right error logging pattern
			log.Println(fmt.Errorf(
				"failed to retrieve events from %s source %s: %w",
				source.Type(),
				source.Name(),
				err,
			))
			continue
		}

		for _, event := range events {
			for _, rule := range e.config.Rules {
				if !rule.EventIsApplicable(event) {
					continue
				}

				channelMessages, err := rule.MessagesFromEvent(event)
				if err != nil {
					log.Println(fmt.Errorf("failed to parse channel messages: %w", err))
				}

				for _, channel := range e.config.Channels {
					msg, ok := channelMessages[channel.Type()]
					if !ok {
						log.Println(fmt.Errorf("did not find message for %s channel %s", channel.Type(), channel.Name()))
					}

					log.Printf("sending event to %s on %s: %s\n", channel.Name(), channel.Type(), event.Name)
					if err := channel.Send(msg); err != nil {
						log.Println(fmt.Errorf(
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
}
