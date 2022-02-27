package rules

import (
	"time"

	"github.com/devict/promobot/sources"
)

type MsgFunc func(sources.Event) string

type NotifyRule struct {
	NumDaysOut       int
	ChannelTemplates map[string]MsgFunc
}

func (rule NotifyRule) MessagesFromEvent(event sources.Event) (map[string]string, error) {
	channelMessages := make(map[string]string)
	for chanType, msgFunc := range rule.ChannelTemplates {
		channelMessages[chanType] = msgFunc(event)
	}
	return channelMessages, nil
}

func (rule NotifyRule) EventIsApplicable(event sources.Event) bool {
	// don't promote the event if it's already started, in case of 0 days out rule
	if time.Now().After(event.DateTime) {
		return false
	}

	checkDate := dateFromTime(time.Now().Add(time.Duration(rule.NumDaysOut*24) * time.Hour))
	eventDate := dateFromTime(event.DateTime)
	return eventDate.Equal(checkDate)
}

func dateFromTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
