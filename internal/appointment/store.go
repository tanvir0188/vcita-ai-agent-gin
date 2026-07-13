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
	SetAutoReplyOn(conversationID string) error
	UpsertAppointment(appt *store.Appointment) error
}

// GORM implementation
func (s *Store) SetAutoReplyOff(conversationID string) error {
	return s.db.
		Model(&store.Conversation{}).
		Where("conversation_id = ?", conversationID).
		Update("auto_reply_off", true).
		Error
}

func (s *Store) SetAutoReplyOn(conversationID string) error {
	return s.db.
		Model(&store.Conversation{}).
		Where("conversation_id = ?", conversationID).
		Update("auto_reply_off", false).
		Error
}

func (s *Store) UpsertAppointment(appt *store.Appointment) error {
	return s.db.
		Where(store.Appointment{
			ConversationID: appt.ConversationID,
			ServiceID:      appt.ServiceID,
		}).
		Assign(store.Appointment{
			ContactID:          appt.ContactID,
			ServiceName:        appt.ServiceName,
			Status:             appt.Status,
			StartTime:          appt.StartTime,
			EndTime:            appt.EndTime,
			VcitaAppointmentID: appt.VcitaAppointmentID,
			VcitaClientID:      appt.VcitaClientID,
		}).
		FirstOrCreate(appt).
		Error
}
