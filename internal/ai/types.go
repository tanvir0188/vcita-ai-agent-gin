package ai

import "time"

type SmartReplyRequest struct {
	MatterUID string `json:"matter_uid"`
	ClientUID string `json:"client_uid"`
}

type SmartReplyResponse struct {
	Success bool           `json:"success"`
	Data    SmartReplyData `json:"data"`
}

type SmartReplyData struct {
	UID       string            `json:"uid"`
	MatterUID string            `json:"matter_uid"`
	Payload   SmartReplyPayload `json:"payload"`
}

type SmartReplyPayload struct {
	EmailMessage string `json:"email_message"`
	SMSMessage   string `json:"sms_message"`
	FBMessage    string `json:"FB_message"`
	OtherMessage string `json:"other_message"`
}

type AppointmentAiResponse struct {
	NeedsScheduling bool       `json:"needs_scheduling"`
	ServiceName     string     `json:"service_name"`
	PreferredTime   *time.Time `json:"preferred_time"`
	BackupTime      *time.Time `json:"backup_time"`
}