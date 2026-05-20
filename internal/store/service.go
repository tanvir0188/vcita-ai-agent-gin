package store

import (
	"time"

	"gorm.io/gorm"
)

type CustomerMessageCreateParams struct {
	ConversationID        string
	LastMessageID         string
	LastCustomerMessageID string
	LastMessageFrom       string

	LastCustomerMessageAt time.Time

	ConversationVersion int64

	AIReplyPending   bool
	PendingMessageID string
}

func (s *DB) CreateOrUpdateConversationOnCustomerMessageCreate(
	params CustomerMessageCreateParams,
) error {

	updates := map[string]interface{}{
		"last_message_id":          params.LastMessageID,
		"last_customer_message_id": params.LastCustomerMessageID,
		"last_message_from":        params.LastMessageFrom,
		"last_customer_message_at": params.LastCustomerMessageAt,
		"conversation_version":     gorm.Expr("conversation_version + 1"),
		"ai_reply_pending":         params.AIReplyPending,
		"pending_message_id":       params.PendingMessageID,
	}

	return s.gorm.
		Model(&Conversation{}).
		Where("conversation_id = ?", params.ConversationID).
		Assign(updates).
		FirstOrCreate(&Conversation{
			ConversationID: params.ConversationID,
		}).
		Error
}

type StaffMessageCreateParams struct {
	ConversationID     string
	LastMessageID      string
	LastStaffMessageID string
	LastMessageFrom    string

	LastStaffMessageAt time.Time

	ConversationVersion int64

	HumanActive   bool
	HumanActiveAt time.Time

	AIReplyPending    bool
	AIReplyGenerating bool

	PendingMessageID string
}

func (s *DB) CreateOrUpdateConversationOnStaffMessage(
	params StaffMessageCreateParams,
) error {

	updates := map[string]interface{}{
		"last_message_id":       params.LastMessageID,
		"last_staff_message_id": params.LastStaffMessageID,
		"last_message_from":     params.LastMessageFrom,
		"last_staff_message_at": params.LastStaffMessageAt,
		"conversation_version":  gorm.Expr("conversation_version + 1"),
		"human_active":          params.HumanActive,
		"human_active_at":       params.HumanActiveAt,
		"ai_reply_pending":      params.AIReplyPending,
		"ai_reply_generating":   params.AIReplyGenerating,
		"pending_message_id":    params.PendingMessageID,
	}

	return s.gorm.
		Model(&Conversation{}).
		Where("conversation_id = ?", params.ConversationID).
		Assign(updates).
		FirstOrCreate(&Conversation{
			ConversationID: params.ConversationID,
		}).
		Error
}

type ConversationReadParams struct {
	ConversationID      string
	HumanActive         bool
	HumanActiveAt       string
	ConversationVersion int
	UpdatedAt           string
}

func (d *DB) CreateOrUpdateConversationRead(
	params ConversationReadParams,
) error {

	updates := map[string]interface{}{
		"human_active":         params.HumanActive,
		"human_active_at":      params.HumanActiveAt,
		"conversation_version": gorm.Expr("conversation_version + 1"),
	}

	return d.gorm.
		Model(&Conversation{}).
		Where("conversation_id = ?", params.ConversationID).
		Assign(updates).
		FirstOrCreate(&Conversation{
			ConversationID: params.ConversationID,
		}).
		Error
}

// GetConversationByID retrieves a conversation by its ConversationID.
func (d *DB) GetConversationByID(conversationID string) (*Conversation, error) {

	var conversation Conversation

	err := d.gorm.
		Where("conversation_id = ?", conversationID).
		First(&conversation).
		Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}
