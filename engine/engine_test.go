package engine_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
)

type TestChannel struct {
	name     string
	sentMsgs []string
}

func NewTestChannel(name string) *TestChannel {
	return &TestChannel{name: name, sentMsgs: make([]string, 0)}
}

func (c *TestChannel) Name() string { return c.name }
func (c *TestChannel) Type() string { return "test" }
func (c *TestChannel) Send(msg string) error {
	c.sentMsgs = append(c.sentMsgs, msg)
	return nil
}

type TestSource struct {
	name   string
	events []sources.Event
}

func NewTestSource(name string, events []sources.Event) *TestSource {
	return &TestSource{name: name, events: events}
}

func (c *TestSource) Name() string { return c.name }
func (c *TestSource) Type() string { return "test" }
func (c *TestSource) Retrieve() ([]sources.Event, error) {
	return c.events, nil
}

func testEvent(name, source string, daysAhead int) sources.Event {
	return sources.Event{
		Name:     name,
		Source:   source,
		URL:      "https://devict.org",
		Location: "Definitely not the metaverse",
		DateTime: time.Now().Add(time.Duration(daysAhead*24) * time.Hour),
	}
}

func TestEngine(t *testing.T) {
	testChannel1 := NewTestChannel("test1")
	testChannel2 := NewTestChannel("test2")

	config := engine.EngineConfig{
		Channels: []channels.Channel{
			channels.Channel(testChannel1),
			channels.Channel(testChannel2),
		},
		Sources: []sources.Source{
			sources.Source(NewTestSource("test-source-1", []sources.Event{
				testEvent("Test Event 1", "test-source-1", 1),
				testEvent("Test Event 2", "test-source-1", 2),
				testEvent("Test Event 3", "test-source-1", 3),
				testEvent("Test Event 4", "test-source-1", 4),
				testEvent("Test Event 5", "test-source-1", 5),
				testEvent("Test Event 6", "test-source-1", 6),
				testEvent("Test Event 7", "test-source-1", 7),
				testEvent("Test Event 8", "test-source-1", 8),
			})),
		},
		Rules: []rules.NotifyRule{
			{
				NumDaysOut: 1,
				ChannelTemplates: map[string]rules.MsgFunc{
					"test": func(e sources.Event) string {
						return fmt.Sprintf("1 day until %s from %s", e.Name, e.Source)
					},
				},
			},
			{
				NumDaysOut: 3,
				ChannelTemplates: map[string]rules.MsgFunc{
					"test": func(e sources.Event) string {
						return fmt.Sprintf("3 days until %s from %s", e.Name, e.Source)
					},
				},
			},
			{
				NumDaysOut: 7,
				ChannelTemplates: map[string]rules.MsgFunc{
					"test": func(e sources.Event) string {
						return fmt.Sprintf("7 days until %s from %s", e.Name, e.Source)
					},
				},
			},
		},
	}

	engine.NewEngine(config).RunOnce()

	// assertions
	expected := []string{
		"1 day until Test Event 1 from test-source-1",
		"3 days until Test Event 3 from test-source-1",
		"7 days until Test Event 7 from test-source-1",
	}

	if len(testChannel1.sentMsgs) != len(expected) {
		t.Fatal("testChannel1 did not receive expected number of messages")
	}
	if len(testChannel2.sentMsgs) != len(expected) {
		t.Fatal("testChannel2 did not receive expected number of messages")
	}

	for _, expectedMsg := range expected {
		if !containsStr(testChannel1.sentMsgs, expectedMsg) {
			t.Fatalf("testChannel1 did not receive message: %s", expectedMsg)
		}
		if !containsStr(testChannel2.sentMsgs, expectedMsg) {
			t.Fatalf("testChannel2 did not receive message: %s", expectedMsg)
		}
	}
}

func containsStr(slice []string, str string) bool {
	found := false
	for _, s := range slice {
		if s == str {
			found = true
		}
	}
	return found
}

func TestShouldRun(t *testing.T) {
	now := time.Now()

	e := engine.NewEngine(engine.EngineConfig{
		RunAt: engine.RunAt{
			Hour:   now.Hour(),
			Minute: now.Minute(),
		},
	})

	if e.ShouldRun(now.Add(-time.Minute)) != false {
		t.Fatal("expected e.ShouldRun to be false a minute before")
	}
	if e.ShouldRun(now) != true {
		t.Fatal("expected e.ShouldRun to be true")
	}
	if e.ShouldRun(now.Add(time.Minute)) != false {
		t.Fatal("expected e.ShouldRun to be false after a minute")
	}
}
