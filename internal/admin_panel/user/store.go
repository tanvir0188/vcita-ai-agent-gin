package user

import (
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func (s *Store) DB() *gorm.DB {
	return s.db
}
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}
