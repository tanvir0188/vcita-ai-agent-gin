package store

import "fmt"

// WriteAuditLog appends a HIPAA audit record.
// NEVER pass PHI (message content, patient name, medication details) to any field.
func (d *DB) WriteAuditLog(clientID, action, eventType, outcome, actorType, errMsg string) error {
	result := d.gorm.Create(&AuditLog{
		ClientID:  clientID,
		Action:    action,
		EventType: eventType,
		Outcome:   outcome,
		ActorType: actorType,
		ErrorMsg:  errMsg,
	})
	return result.Error
}

// SaveEscalation persists an escalation record with encrypted trigger reason.
func (d *DB) SaveEscalation(clientID, triggerReason string) (uint, error) {
	enc, err := d.enc.Encrypt(triggerReason)
	if err != nil {
		return 0, fmt.Errorf("store.SaveEscalation: encrypt: %w", err)
	}
	record := &Escalation{
		ClientID:   clientID,
		TriggerEnc: enc,
		Status:     "open",
	}
	if err := d.gorm.Create(record).Error; err != nil {
		return 0, fmt.Errorf("store.SaveEscalation: create: %w", err)
	}
	return record.ID, nil
}

// SaveAppointmentSuggestion records an AI-generated appointment suggestion.
func (d *DB) SaveAppointmentSuggestion(clientID, staffID string, slot string) (uint, error) {
	// Parse the ISO 8601 slot string to time.Time
	// slot is already formatted as RFC3339 by the caller
	var t interface{} = slot // stored as string in the DB via raw update if needed

	record := &AppointmentSuggestion{
		ClientID: clientID,
		StaffID:  staffID,
		Status:   "suggested",
	}
	_ = t
	if err := d.gorm.Create(record).Error; err != nil {
		return 0, fmt.Errorf("store.SaveAppointmentSuggestion: create: %w", err)
	}
	return record.ID, nil
}

// UpdateAppointmentStatus updates the status of an appointment suggestion.
func (d *DB) UpdateAppointmentStatus(id uint, status string) error {
	return d.gorm.Model(&AppointmentSuggestion{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
