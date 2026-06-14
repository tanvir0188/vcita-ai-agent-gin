package store

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
