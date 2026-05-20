// Package webhook handles inTandem webhook events.
// Event payload types are defined in internal/vcita/types.go
// to keep the vcita API surface in one place.
package webhook

import "time"

type WebhookEnvelope struct {
	EntityName string         `json:"entity_name"`
	EventType  string         `json:"event_type"`
	Data       WebhookPayload `json:"data"`
}

type WebhookPayload struct {
	UID             string `json:"uid"`
	ConversationUID string `json:"conversation_uid"`
	ContactUID      string `json:"contact_uid"`
	BusinessUID     string `json:"business_uid"`

	StaffUid string `json:"staff_uid"`

	MessageType string `json:"message_type"`

	Text      string `json:"text"`
	Direction string `json:"direction"`

	CreatedAt time.Time `json:"created_at"`
}

type ConversationMessage struct {
	UID               string  `json:"uid"`
	ConversationUID   string  `json:"conversation_uid"`
	ContactUID        string  `json:"contact_uid"`
	BusinessUID       string  `json:"business_uid"`
	StaffUID          string  `json:"staff_uid"`
	MessageType       string  `json:"message_type"`
	Direction         string  `json:"direction"`
	Text              string  `json:"text"`
	MessageCategory   *string `json:"message_category"`
	MessageEntityUID  *string `json:"message_entity_uid"`
	CreatedAt         string  `json:"created_at"`
	ExternalMessageID *string `json:"external_message_id"`
	ChannelUID        string  `json:"channel_uid"`
}

type AiEvaluationParam struct {
	ConversationID      string
	ConversationVersion int64
	PendingMessageID    string
	ContactId           string
}
