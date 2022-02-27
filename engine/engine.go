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
	Channels  []channels.Channel
	Sources   []sources.Source
	Rules     []rules.NotifyRule
	SleepTime time.Duration
}

type Engine struct {
	config EngineConfig
}

func NewEngine(config EngineConfig) *Engine {
	return &Engine{config}
}

func (e *Engine) Run() {
	for {
		e.RunOnce()
		time.Sleep(e.config.SleepTime)
	}
}

func (e *Engine) RunOnce() {
	for _, source := range e.config.Sources {
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

		for _, event := range events {
			for _, rule := range e.config.Rules {
				if !rule.EventIsApplicable(event) {
					continue
				}

				channelMessages, err := rule.MessagesFromEvent(event)
				if err != nil {
					log.Print(fmt.Errorf("failed to parse channel messages: %w", err))
				}

				for _, channel := range e.config.Channels {
					msg, ok := channelMessages[channel.Type()]
					if !ok {
						log.Print(fmt.Errorf("did not find message for %s channel %s", channel.Type(), channel.Name()))
					}

					log.Printf("sending event to %s on %s: %s\n", channel.Name(), channel.Type(), event.Name)
					if err := channel.Send(msg); err != nil {
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
}
