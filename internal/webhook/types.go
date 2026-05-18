// Package webhook handles inTandem webhook events.
// Event payload types are defined in internal/vcita/types.go
// to keep the vcita API surface in one place.
package webhook

type WebhookEnvelope struct {
	EntityName string              `json:"entity_name"`
	EventType  string              `json:"event_type"`
	Data       ConversationMessage `json:"data"`
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


