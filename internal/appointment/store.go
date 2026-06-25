package appointment

import (
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Store interface
type AppointmentStore interface {
	SetAutoReplyOff(conversationID string) error
}

// GORM implementation
func (s *Store) SetAutoReplyOff(conversationID string) error {
	return s.db.
		Model(&store.Conversation{}).
		Where("matter_uid = ?", conversationID).
		Update("auto_reply_off", true).
		Error
}
