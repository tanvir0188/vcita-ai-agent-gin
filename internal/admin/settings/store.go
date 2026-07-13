package settings

import (
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

// NewStore creates a new instance of Store.
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetSystemSetting() (*store.SystemSetting, error) {
	var setting store.SystemSetting
	err := s.db.FirstOrCreate(&setting, store.SystemSetting{ID: 1, SystemEnabled: true}).Error
	return &setting, err
}

func (s *Store) UpdateSystemSetting(enabled bool) error {
	var setting store.SystemSetting
	err := s.db.FirstOrCreate(&setting, store.SystemSetting{ID: 1, SystemEnabled: true}).Error
	if err != nil {
		return err
	}
	setting.SystemEnabled = enabled
	return s.db.Save(&setting).Error
}
