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
	NowFunc                func() time.Time
	Channels               []channels.Channel
	Sources                []sources.Source
	Rules                  []rules.NotifyRule
	WeeklySummaryTemplates map[string]rules.WeeklySummaryFunc
	WeeklySummaryDay       time.Weekday
	RunAt                  RunAt
	Location               *time.Location
	DebugMode              bool
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
		events, err := source.Retrieve(e.config.Location)
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

		eventsInNextWeek := []sources.Event{}
		dayOfWeek := e.config.NowFunc().Weekday()

		for _, event := range events {
			if dayOfWeek == event.DateTime.Weekday() && eventIsWithinNextWeek(event, e.config.NowFunc()) {
				eventsInNextWeek = append(eventsInNextWeek, event)
			}

			for _, rule := range e.config.Rules {
				if !rule.EventIsApplicable(event, e.config.Location, e.config.NowFunc()) {
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
					e.sendMessage(channel, msg)
				}
			}
		}

		if dayOfWeek == e.config.WeeklySummaryDay {
			for _, channel := range e.config.Channels {
				msgFunc, ok := e.config.WeeklySummaryTemplates[channel.Type()]
				if !ok {
					log.Println(fmt.Errorf("did not find message for %s channel %s", channel.Type(), channel.Name()))
				} else {
					e.sendMessage(channel, msgFunc(eventsInNextWeek))
				}
			}
		}
	}
}

func (e *Engine) sendMessage(channel channels.Channel, message string) {
	log.Printf("sending weekly summary to %s on %s\n", channel.Name(), channel.Type())

	if e.config.DebugMode {
		fmt.Printf("%s\n\n", message)
	} else {
		if err := channel.Send(message); err != nil {
			log.Println(fmt.Errorf(
				"failed to send to %s channel %s: %w",
				channel.Type(),
				channel.Name(),
				err,
			))
		}
	}
}

func eventIsWithinNextWeek(event sources.Event, now time.Time) bool {
	diff := now.Sub(event.DateTime)
	return diff < 7*24*time.Hour
}
