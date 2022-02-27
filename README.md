# promobot

This project aims to automate the promotion of devICT events, and events of local groups. It connects one or more sources to one or more channels.

## Sources

A `Source` represents a source of events data from a specific organization. promobot will iterate through each defined source to get a list of upcoming events for promotion.

## Channels

A `Channel` represents something that can have an event promoted through it. Examples are Slack or Twitter.

## NotifyRules

A a list of `NotifyRule`s describes when to send promotional messages about an event to various channels. A rule defines how many days out to promote the event, and what the message should look like for events promoted at that time.
