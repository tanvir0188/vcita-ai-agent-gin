package utils

type ConversationMessage struct {
	UID             string `json:"uid"`
	ConversationUID string `json:"conversation_uid"`
	MessageType     string `json:"message_type"`
	Direction       string `json:"direction"`
	Text            string `json:"text"`

	CreatedAt  string `json:"created_at"`
	ChannelUID string `json:"channel_uid"`
}

type GetClientDetailParams struct {
	webhookSecret string
	contactId     string
}
