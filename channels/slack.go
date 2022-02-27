package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackChannel struct {
	name string
	url  string
}

func NewSlackChannel(name, url string) *SlackChannel {
	return &SlackChannel{name: name, url: url}
}

func (c *SlackChannel) Type() string { return "slack" }
func (c *SlackChannel) Name() string { return c.name }

func (c *SlackChannel) Send(message string) error {
	messageStr, err := json.Marshal(SlackMessage{Text: message})
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	_, err = http.Post(c.url, "application/json", bytes.NewReader(messageStr))
	if err != nil {
		return fmt.Errorf("failed to post to slack: %w", err)
	}

	return nil
}

type SlackMessage struct {
	Text string `json:"text"`
}
