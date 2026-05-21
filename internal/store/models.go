// Package store handles all database operations using GORM.
// PHI fields (message content, medication data, escalation triggers) are
// encrypted with AES-256-GCM before insert and decrypted after select.
package store

import (
	"time"

	"gorm.io/gorm"
)

// ── GORM models ───────────────────────────────────────────────────────────────
// Each struct maps directly to a database table.
// PHI columns are named *Enc and hold AES-256-GCM ciphertext.
// Non-PHI columns (IDs, timestamps, booleans, dates) are stored plaintext
// so GORM can query and index them without decryption.

// Message is one turn in a patient conversation.
type Conversation struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`

	ConversationID string `gorm:"uniqueIndex;not null"`

	LastMessageID         string
	LastCustomerMessageID string
	LastStaffMessageID    string

	LastMessageFrom string

	LastCustomerMessageAt time.Time
	LastStaffMessageAt    time.Time

	HumanActive   bool      `gorm:"default:false"`
	HumanActiveAt time.Time `gorm:"index"`

	HumanActiveUntil time.Time `gorm:"index"`

	ConversationVersion int64 `gorm:"default:0"`

	AIReplyPending    bool `gorm:"default:false;index"`
	AIReplyGenerating bool `gorm:"default:false"`

	HasEscalated bool `gorm:"default:false"`

	AutoReplyOffUntil time.Time

	PendingMessageID string
}

// RefillSchedule holds a medication refill reminder for one patient.
type RefillSchedule struct {
	gorm.Model
	ClientID      string    `gorm:"not null;index"`
	MedicationEnc string    `gorm:"not null"` // encrypted: name/dosage/frequency JSON
	RefillDate    time.Time `gorm:"not null;index"`
	ReminderDate  time.Time `gorm:"not null;index"` // RefillDate - 7 days
	ReminderSent  bool      `gorm:"not null;default:false"`
}

// AuditLog records every system action for HIPAA traceability.
// No PHI is stored here — only IDs, action types, and outcomes.
type AuditLog struct {
	gorm.Model
	ClientID  string `gorm:"not null;index"` // vcita client_id only, no name/contact
	Action    string `gorm:"not null;index"`
	EventType string `gorm:"not null;default:''"`
	Outcome   string `gorm:"not null"` // "success" | "error" | "escalated"
	ActorType string `gorm:"not null;default:'system'"`
	ErrorMsg  string `gorm:"default:''"` // never contains PHI
}

// Escalation records an escalation event requiring human review.
type Escalation struct {
	gorm.Model
	ClientID   string `gorm:"not null;index"`
	TriggerEnc string `gorm:"not null"`                // encrypted: reason for escalation
	Status     string `gorm:"not null;default:'open'"` // "open" | "acknowledged" | "resolved"
	ResolvedAt *time.Time
}

// AppointmentSuggestion logs an AI-generated appointment recommendation.
type AppointmentSuggestion struct {
	gorm.Model
	ClientID      string    `gorm:"not null;index"`
	SuggestedSlot time.Time `gorm:"not null"`
	StaffID       string    `gorm:"not null"`
	Status        string    `gorm:"not null;default:'suggested'"` // "suggested" | "confirmed" | "rejected" | "escalated"
}

// ── Decrypted view types (used in application logic, never persisted) ─────────

// ConversationMessage is the decrypted form of a Message row.
type ConversationMessage struct {
	ID          uint
	ClientID    string
	Role        string
	Content     string // decrypted plaintext
	EventType   string
	IsEscalated bool
	CreatedAt   time.Time
}

// RefillScheduleView is the decrypted form of a RefillSchedule row.
type RefillScheduleView struct {
	ID           uint
	ClientID     string
	Medication   string // decrypted
	RefillDate   time.Time
	ReminderDate time.Time
	ReminderSent bool
}
