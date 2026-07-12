package ai

import "time"

type AppointmentAiResponse struct {
	NeedsScheduling    bool       `json:"needs_scheduling"`
	ServiceName        string     `json:"service_name"`
	PreferredTime      *time.Time `json:"preferred_time"`
	InteractionDetails string     `json:"interaction_details"`
	BackupTime         *time.Time `json:"backup_time"`
	StartTime          *time.Time `json:"start_time"`
	EndTime            *time.Time `json:"end_time"`
}

// MedicationEntry represents a single medication in the refill prediction response.
type MedicationEntry struct {
	Name           string  `json:"name"`
	Strength       *string `json:"strength"`
	Quantity       *int    `json:"quantity"`
	Frequency      *string `json:"frequency"`
	NextRefillDate *string `json:"next_refill_date"` // YYYY-MM-DD
	Confidence     string  `json:"confidence"`       // "high" | "estimated" | "unknown"
	NeedsRefill    bool    `json:"needs_refill"`
}

// MedicationRefillResponse is the root object returned by PredictMedicationRefill.
type MedicationRefillResponse struct {
	SourceType      string            `json:"source_type"`
	Provider        *string           `json:"provider"`
	OrderDate       *string           `json:"order_date"`
	Medications     []MedicationEntry `json:"medications"`
	Response        *string           `json:"response"`
	PartialReminder bool              `json:"partial_reminder"`
	ReminderNeeded  bool              `json:"reminder_needed"`
}
