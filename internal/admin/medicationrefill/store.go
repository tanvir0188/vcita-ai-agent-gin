package medicationrefill

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

func (s *Store) ToggleReminder(id string) error {
	var clientSyncState store.ClientSyncState
	err := s.db.Where("id = ?", id).First(&clientSyncState).Error
	if err != nil {
		return err
	}
	clientSyncState.RemindMedication = !clientSyncState.RemindMedication
	return s.db.Save(&clientSyncState).Error
}
