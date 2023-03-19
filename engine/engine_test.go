package engine_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/devict/promobot/channels"
	"github.com/devict/promobot/engine"
	"github.com/devict/promobot/rules"
	"github.com/devict/promobot/sources"
)

var octoberFirst = time.Date(2012, time.October, 1, 8, 0, 0, 0, time.Local)

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

func (c *TestSource) Name() string    { return c.name }
func (c *TestSource) Type() string    { return "test" }
func (c *TestSource) JsonUrl() string { return "nope" }
func (c *TestSource) HtmlUrl() string { return "https://meetup.com/devict" }
func (c *TestSource) Retrieve(t *time.Location) ([]sources.Event, error) {
	return c.events, nil
}

func testEvent(name, source string, daysAhead int) sources.Event {
	return sources.Event{
		Name:     name,
		Source:   source,
		URL:      "https://devict.org",
		Location: "Definitely not the metaverse",
		DateTime: octoberFirst.Add(time.Duration(daysAhead*24) * time.Hour),
	}
}

func TestEngine(t *testing.T) {
	testChannel1 := NewTestChannel("test1")
	testChannel2 := NewTestChannel("test2")

	config := engine.EngineConfig{
		NowFunc: func() time.Time { return octoberFirst },
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
		WeeklySummaryDay: octoberFirst.Weekday(),
		WeeklySummaryTemplates: map[string]rules.WeeklySummaryFunc{
			"test": func(events []sources.Event) string {
				eventLines := []string{}
				for _, event := range events {
					day := event.DateTime.Weekday().String()
					eventLines = append(eventLines, fmt.Sprintf("- [%s] <%s|%s>", day[:3], event.URL, event.Name))
				}
				return fmt.Sprintf("Events this week!\n\n%s", strings.Join(eventLines, "\n"))
			},
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
		strings.Join([]string{
			"Events this week!\n",
			"- [Tue] <https://devict.org|Test Event 1>",
			"- [Wed] <https://devict.org|Test Event 2>",
			"- [Thu] <https://devict.org|Test Event 3>",
			"- [Fri] <https://devict.org|Test Event 4>",
			"- [Sat] <https://devict.org|Test Event 5>",
			"- [Sun] <https://devict.org|Test Event 6>",
		}, "\n"),
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
	now := octoberFirst

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
