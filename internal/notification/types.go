package notification

import "github.com/tanvir0188/vcita-ai-agent/internal/vcita"
type SlackPayload struct {
	Text string `json:"text"`
}

// Handler or Service struct that holds your dependencies

type SlackBlockPayload struct {
	Text   string       `json:"text,omitempty"`
	Blocks []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type     string         `json:"type"`
	Text     *SlackText     `json:"text,omitempty"`
	Fields   []SlackText    `json:"fields,omitempty"`
	Elements []SlackElement `json:"elements,omitempty"`
}

type SlackText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type SlackElement struct {
	Type  string     `json:"type"`
	Text  *SlackText `json:"text,omitempty"`
	Style string     `json:"style,omitempty"`
	Value string     `json:"value,omitempty"`
}

// SlackClient holds the configuration for a Slack escalation notification.
type SlackClient struct {
	WebhookURL string
	Text       string
	Client     *vcita.ClientInfo
}