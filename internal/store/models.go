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

type User struct {
	gorm.Model
	StaffUID    string `gorm:"uniqueIndex;size:255"`
	FullName    string `gorm:"size:255"`
	Email       string `gorm:"uniqueIndex;size:255"`
	PhoneNumber string `gorm:"size:20"`
	Password    string `gorm:"not null"`

	IsActive   bool `gorm:"default:false"`
	IsVerified bool `gorm:"default:false"`
	IsAdmin    bool `gorm:"default:false"`

	OtpCode      string
	OtpExpiresAt time.Time
}

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

	AutoReplyOffUntil *time.Time

	PendingMessageID string
}

// RefillSchedule holds a medication refill reminder for one patient.

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
