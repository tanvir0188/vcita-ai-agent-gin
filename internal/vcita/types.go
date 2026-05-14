// Package vcita provides a typed HTTP client for the inTandem / vcita platform API.
// Base URL: https://api.vcita.biz
// Auth: Bearer token in Authorization header
package vcita

import "time"

// ── Webhook payload types ─────────────────────────────────────────────────────

// WebhookEnvelope is the top-level structure inTandem POSTs to our endpoint.
type WebhookEnvelope struct {
	EventType string      `json:"event_type"` // e.g. "message.client_sent_message"
	BusinessUID string    `json:"business_uid"`
	Payload   interface{} `json:"payload"` // decoded based on EventType
}

// MessagePayload is the payload for message.client_sent_message and
// message.business_sent_message events.
type MessagePayload struct {
	ClientID        string    `json:"client_id"`
	MessageID       string    `json:"message_id"`
	ConversationID  string    `json:"conversation_id"`
	MessageType     string    `json:"message_type"` // "text", "api", "general_question", etc.
	Direction       string    `json:"direction"`    // "client_to_pivot" or "pivot_to_client"
	Body            string    `json:"body"`
	Channel         string    `json:"channel"` // "portal", "sms", "email"
	CreatedAt       time.Time `json:"created_at"`
}

// NotePayload is a simplified representation of a vcita client note.
type NotePayload struct {
	NoteID    string    `json:"id"`
	ClientID  string    `json:"client_id"`
	NoteType  string    `json:"note_type"` // "current_medications", "general", etc.
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AppointmentPayload represents the appointment.requested webhook payload.
type AppointmentPayload struct {
	ClientID        string    `json:"client_id"`
	AppointmentID   string    `json:"appointment_id"`
	ServiceName     string    `json:"service_name"`
	RequestedTime   *time.Time `json:"requested_time,omitempty"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// ── API request/response types ────────────────────────────────────────────────

// SendMessageRequest is the body for POST /v3/communication/messages (or equivalent).
type SendMessageRequest struct {
	ClientID       string `json:"client_id"`
	ConversationID string `json:"conversation_id,omitempty"`
	Body           string `json:"body"`
	Channel        string `json:"channel,omitempty"` // portal | sms | email
}

// Client represents a vcita client record.
type Client struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// Note represents a vcita client note.
type Note struct {
	ID       string `json:"id"`
	NoteType string `json:"note_type"`
	Content  string `json:"content"`
}

// StaffMember represents a vcita staff member.
type StaffMember struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// TimeSlot represents an available scheduling slot.
type TimeSlot struct {
	StaffID   string    `json:"staff_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ScheduleAppointmentRequest is the body to confirm/create an appointment via API.
type ScheduleAppointmentRequest struct {
	ClientID    string    `json:"client_id"`
	StaffID     string    `json:"staff_id"`
	ServiceName string    `json:"service_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

// APIResponse is the generic inTandem response envelope.
type APIResponse[T any] struct {
	Status string `json:"status"` // "OK" or "ERROR"
	Data   T      `json:"data"`
}
