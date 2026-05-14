package store

import (
	"fmt"
)

// SaveMessage encrypts the content and persists a message row via GORM.
func (d *DB) SaveMessage(clientID, role, content, eventType string, isEscalated bool) error {
	enc, err := d.enc.Encrypt(content)
	if err != nil {
		return fmt.Errorf("store.SaveMessage: encrypt: %w", err)
	}
	result := d.gorm.Create(&Message{
		ClientID:    clientID,
		Role:        role,
		ContentEnc:  enc,
		EventType:   eventType,
		IsEscalated: isEscalated,
	})
	return result.Error
}

// GetConversationHistory returns the last n messages for a client,
// in chronological order, with content decrypted.
func (d *DB) GetConversationHistory(clientID string, limit int) ([]ConversationMessage, error) {
	var rows []Message
	result := d.gorm.
		Where("client_id = ?", clientID).
		Order("created_at DESC").
		Limit(limit).
		Find(&rows)
	if result.Error != nil {
		return nil, fmt.Errorf("store.GetConversationHistory: query: %w", result.Error)
	}

	// Reverse to chronological order and decrypt
	msgs := make([]ConversationMessage, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		r := rows[i]
		plain, err := d.enc.Decrypt(r.ContentEnc)
		if err != nil {
			return nil, fmt.Errorf("store.GetConversationHistory: decrypt id=%d: %w", r.ID, err)
		}
		msgs = append(msgs, ConversationMessage{
			ID:          r.ID,
			ClientID:    r.ClientID,
			Role:        r.Role,
			Content:     plain,
			EventType:   r.EventType,
			IsEscalated: r.IsEscalated,
			CreatedAt:   r.CreatedAt,
		})
	}
	return msgs, nil
}
