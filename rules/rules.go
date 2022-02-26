package rules

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/devict/promobot/sources"
)

type NotificationRule struct {
	NumDaysOut       int
	ChannelTemplates map[string]*template.Template
}

func MessagesFromEvent(event sources.Event, rules []NotificationRule) (map[string]string, error) {
	channelMessages := make(map[string]string)
	for _, rule := range rules {
		for chanType, tmpl := range rule.ChannelTemplates {
			var out bytes.Buffer
			if err := tmpl.Execute(&out, event); err != nil {
				return channelMessages, fmt.Errorf("failed to compile template: %w", err)
			}
			channelMessages[chanType] = out.String()
		}
	}
	return channelMessages, nil
}
