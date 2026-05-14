package store

import (
	"fmt"
	"time"
)

// UpsertRefillSchedule inserts or updates the refill schedule for a client.
// Medication info is encrypted before storage.
// GORM's Save() performs an upsert when the primary key is present;
// here we use FirstOrCreate + Updates for explicit control.
func (d *DB) UpsertRefillSchedule(clientID, medicationInfo string, refillDate time.Time) error {
	enc, err := d.enc.Encrypt(medicationInfo)
	if err != nil {
		return fmt.Errorf("store.UpsertRefillSchedule: encrypt: %w", err)
	}

	reminderDate := refillDate.AddDate(0, 0, -7)

	var existing RefillSchedule
	result := d.gorm.Where("client_id = ?", clientID).First(&existing)

	if result.Error != nil {
		// No existing record — create one
		return d.gorm.Create(&RefillSchedule{
			ClientID:      clientID,
			MedicationEnc: enc,
			RefillDate:    refillDate.UTC(),
			ReminderDate:  reminderDate.UTC(),
			ReminderSent:  false,
		}).Error
	}

	// Update existing record and reset reminder_sent
	return d.gorm.Model(&existing).Updates(map[string]interface{}{
		"medication_enc": enc,
		"refill_date":    refillDate.UTC(),
		"reminder_date":  reminderDate.UTC(),
		"reminder_sent":  false,
	}).Error
}

// GetDueReminders returns all refill schedules whose reminder_date is today or past
// and reminder has not yet been sent.
func (d *DB) GetDueReminders() ([]RefillScheduleView, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var rows []RefillSchedule
	result := d.gorm.
		Where("reminder_date <= ? AND reminder_sent = ?", today, false).
		Find(&rows)
	if result.Error != nil {
		return nil, fmt.Errorf("store.GetDueReminders: query: %w", result.Error)
	}

	views := make([]RefillScheduleView, 0, len(rows))
	for _, r := range rows {
		med, err := d.enc.Decrypt(r.MedicationEnc)
		if err != nil {
			return nil, fmt.Errorf("store.GetDueReminders: decrypt id=%d: %w", r.ID, err)
		}
		views = append(views, RefillScheduleView{
			ID:           r.ID,
			ClientID:     r.ClientID,
			Medication:   med,
			RefillDate:   r.RefillDate,
			ReminderDate: r.ReminderDate,
			ReminderSent: r.ReminderSent,
		})
	}
	return views, nil
}

// MarkReminderSent sets reminder_sent = true for the given record ID.
func (d *DB) MarkReminderSent(id uint) error {
	return d.gorm.Model(&RefillSchedule{}).
		Where("id = ?", id).
		Update("reminder_sent", true).
		Error
}
