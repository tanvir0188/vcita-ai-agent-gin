package types

import "time"

type UserListResponse struct {
	ID          uint      `json:"id"`
	StaffUID    string    `json:"staff_uid"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
}

type Conversation struct {
	Id             uint   `json:"id"`
	ConversationID string `json:"conversation_id"`

	LastMessageID         string `json:"last_message_id"`
	LastCustomerMessageID string `json:"last_customer_message_id"`
	LastStaffMessageID    string `json:"last_staff_message_id"`

	LastMessageFrom string `json:"last_message_from"`

	LastCustomerMessageAt time.Time `json:"last_customer_message_at"`
	LastStaffMessageAt    time.Time `json:"last_staff_message_at"`

	HumanActive   bool      `json:"human_active"`
	HumanActiveAt time.Time `json:"human_active_at"`

	HumanActiveUntil time.Time `json:"human_active_until"`

	ConversationVersion int64 `json:"conversation_version"`

	AIReplyPending    bool `json:"ai_reply_pending"`
	AIReplyGenerating bool `json:"ai_reply_generating"`

	HasEscalated bool `json:"has_escalated"`

	AutoReplyOffUntil *time.Time `json:"auto_reply_off_until"`

	PendingMessageID string `json:"pending_message_id"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterPayload struct {
	StaffUID    string `json:"staff_uid" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,min=3,max=130"`
}
