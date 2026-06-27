// internal/appointment/handler.go

package appointment

import (
	"fmt"
	"strings"
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
	"go.uber.org/zap"
)

// internal/appointment/handler.go

type SchedulingParams struct {
	ConversationID string // this is matter_uid
	ContactID      string // this is client_id from the webhook
	StaffMessage   string
	Log            *zap.Logger
}

func HandleAppointmentScheduling(db store.Store, p SchedulingParams) error {
	// 1. Extract scheduling intent from the platform AI's message.
	aiResp, err := ai.ConfirmedAppointmentDate(p.StaffMessage)
	if err != nil {
		return fmt.Errorf("AI extraction failed: %w", err)
	}

	startTimeStr := "nil"
	if aiResp.StartTime != nil {
		startTimeStr = aiResp.StartTime.Format(time.RFC3339)
	}
	endTimeStr := "nil"
	if aiResp.EndTime != nil {
		endTimeStr = aiResp.EndTime.Format(time.RFC3339)
	}
	p.Log.Info("appointment AI result",
		zap.Bool("needs_scheduling", aiResp.NeedsScheduling),
		zap.String("service_name", aiResp.ServiceName),
		zap.String("start_date", startTimeStr),
		zap.String("end_date", endTimeStr),
	)

	if !aiResp.NeedsScheduling {
		return nil
	}

	if err := db.SetAutoReplyOff(p.ConversationID); err != nil {
		return fmt.Errorf("failed to pause auto-reply: %w", err)
	}

	// 3. Match service name → vcita service ID.
	serviceID, err := matchService(aiResp.ServiceName)
	if err != nil {
		return fmt.Errorf("service match failed: %w", err)
	}

	// 4. Try preferred time, then backup.
	slot, err := findFirstAvailableSlot(serviceID, aiResp.PreferredTime, aiResp.BackupTime)
	if err != nil {
		// No slot found — leave auto-reply paused for manual follow-up.

		p.Log.Warn("no available slot found, leaving auto-reply paused",
			zap.String("conversation_id", p.ConversationID),
			zap.Error(err),
		)
		utils.CreateMessage(p.ContactID, "No slot found in the given time")
		if err := db.SetAutoReplyOn(p.ConversationID); err != nil {
			return fmt.Errorf("failed to pause auto-reply: %w", err)
		}
		return nil

	}

	p.Log.Info("slot found",
		zap.String("start_time", slot.StartTime),
		zap.String("staff_id", slot.StaffID),
	)

	// 5. Book the appointment.
	booking, err := BookAppointment(config.Envs.VcitaBusinessToken, BookingRequest{
		BusinessID: config.Envs.BusinessUid,
		ClientID:   p.ContactID,
		MatterUID:  p.ConversationID,
		ServiceID:  serviceID,
		StaffID:    slot.StaffID,
		StartTime:  slot.StartTime,
	})
	if err != nil {
		_ = db.UpsertAppointment(&store.Appointment{
			ConversationID: p.ConversationID,
			ContactID:      p.ContactID,
			ServiceID:      serviceID,
			ServiceName:    aiResp.ServiceName,
			Status:         "failed",
		})
		return fmt.Errorf("booking API failed: %w", err)
	}

	// 6. Re-enable auto-replies now that booking is done.
	if err := db.SetAutoReplyOn(p.ConversationID); err != nil {
		p.Log.Error("booked successfully but failed to re-enable auto-reply",
			zap.String("conversation_id", p.ConversationID),
			zap.Error(err),
		)
	}

	// 7. Persist the appointment.
	bookedStart, err := time.Parse(time.RFC3339, booking.StartTime)
	if err != nil {
		bookedStart = *aiResp.PreferredTime
	}

	vcitaID := booking.Title
	if err := db.UpsertAppointment(&store.Appointment{
		ConversationID:     p.ConversationID,
		ContactID:          p.ContactID,
		ServiceID:          serviceID,
		ServiceName:        aiResp.ServiceName,
		Status:             "scheduled",
		StartTime:          &bookedStart,
		VcitaAppointmentID: &vcitaID,
		VcitaClientID:      p.ContactID,
	}); err != nil {
		return fmt.Errorf("booked in vcita but DB upsert failed: %w", err)
	}

	p.Log.Info("appointment booked and auto-reply re-enabled",
		zap.String("conversation_id", p.ConversationID),
		zap.String("start_time", booking.StartTime),
	)
	return nil
}

func findFirstAvailableSlot(serviceID string, preferred, backup *time.Time) (*AvailabilitySlot, error) {
	if preferred == nil {
		return nil, fmt.Errorf("no preferred time provided")
	}

	// Search a window around the preferred date (same day + 1 day covers timezone shifts).
	const layout = "2006-01-02"
	startDate := preferred.Format(layout)
	endDate := preferred.AddDate(0, 0, 1).Format(layout)

	availability, err := GetAvailabilitySlots(serviceID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("availability fetch failed: %w", err)
	}

	if slot := matchSlotToTime(availability, preferred); slot != nil {
		return slot, nil
	}

	// Try backup time if preferred had no match.
	if backup != nil {
		backupStart := backup.Format(layout)
		backupEnd := backup.AddDate(0, 0, 1).Format(layout)

		backupAvailability, err := GetAvailabilitySlots(serviceID, backupStart, backupEnd)
		if err != nil {
			return nil, fmt.Errorf("backup availability fetch failed: %w", err)
		}

		if slot := matchSlotToTime(backupAvailability, backup); slot != nil {
			return slot, nil
		}
	}

	return nil, fmt.Errorf("no open slot found for preferred or backup time")
}

func matchSlotToTime(availability *AvailabilityResponse, target *time.Time) *AvailabilitySlot {
	for _, slots := range availability.Data.Availabilities {
		for i := range slots {
			slot := &slots[i]
			if slot.SpotsOpen == 0 {
				continue
			}

			slotTime, err := time.Parse(time.RFC3339, slot.StartTime)
			if err != nil {
				continue
			}

			diff := slotTime.Sub(*target)
			if diff < 0 {
				diff = -diff
			}
			if diff <= 15*time.Minute {
				return slot
			}
		}
	}
	return nil
}

func matchService(aiName string) (string, error) {
	services, err := GetServices()
	if err != nil {
		return "", err
	}

	lower := strings.ToLower(aiName)
	for _, s := range services {
		if strings.Contains(strings.ToLower(s.Name), lower) ||
			strings.Contains(lower, strings.ToLower(s.Name)) {
			return s.ID, nil
		}
	}

	return "", fmt.Errorf("no vcita service matched %q (checked %d services)", aiName, len(services))
}
