package medication

import (
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"gorm.io/gorm"
)

type DBStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) CreateOrUpdateReminder(matterUID string, resp *ai.MedicationRefillResponse) error {
	// Find existing sync state
	var syncState store.ClientSyncState
	err := s.db.Where("matter_uid = ?", matterUID).First(&syncState).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	now := time.Now()
	// Prepare fields based on response
	remind := resp.ReminderNeeded
	partial := resp.PartialReminder
	hasMeds := len(resp.Medications) > 0
	if err == gorm.ErrRecordNotFound {
		// Create new record
		syncState = store.ClientSyncState{
			MatterUID:             matterUID,
			RemindMedication:      remind,
			HasMedications:        hasMeds,
			PartialReminderNeeded: partial,
			LastAICheckAt:         &now,
		}
		return s.db.Create(&syncState).Error
	}
	// Update existing record
	syncState.RemindMedication = remind
	syncState.HasMedications = hasMeds
	syncState.PartialReminderNeeded = partial
	syncState.LastAICheckAt = &now
	return s.db.Save(&syncState).Error
}
