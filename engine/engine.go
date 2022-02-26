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
	Rules    []rules.NotificationRule
}

type Engine struct {
	config    EngineConfig
	sleepTime time.Duration
}

func NewEngine(config EngineConfig, sleepTime time.Duration) *Engine {
	return &Engine{config, sleepTime}
}

func (e *Engine) Run() {
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
			channelMessages, err := rules.MessagesFromEvent(event, e.config.Rules)
			if err != nil {
				log.Print(fmt.Errorf("failed to parse channel messages: %w", err))
			}

			for _, channel := range e.config.Channels {
				msg, ok := channelMessages[channel.Type()]
				if !ok {
					log.Print(fmt.Errorf("did not find message for %s channel %s", channel.Type(), channel.Name()))
				}

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
