// Package webhook handles incoming inTandem webhook events via Gin.
// Every request is HMAC-SHA256 signature-verified before any processing.
// Message body content is never logged (HIPAA compliance).
package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// AIService is the interface the handler uses for AI processing.
// The AI developer implements this; the stub satisfies it until then.
type AIService interface {
	ProcessMessage(ctx context.Context, clientID string, history []store.ConversationMessage, notes []vcita.Note, medication string) (reply string, escalate bool, err error)
	ExtractMedication(ctx context.Context, noteContent string) (medicationInfo string, supplyDays int, err error)
	SuggestAppointment(ctx context.Context, req AppointmentRequest) (*AppointmentSuggestion, error)
}

// AppointmentRequest is passed to the AI service for slot selection.
type AppointmentRequest struct {
	ClientID        string
	RequestedTime   *time.Time
	ServiceName     string
	AvailableSlots  []vcita.TimeSlot
	AvailableStaff  []vcita.StaffMember
	ConversationCtx []store.ConversationMessage
}

// AppointmentSuggestion is what the AI returns after evaluating slots.
type AppointmentSuggestion struct {
	Slot                    vcita.TimeSlot
	Staff                   vcita.StaffMember
	ConfirmationMessage     string
	HumanInterventionNeeded bool
}

// Handler holds all dependencies for processing webhook events.
type Handler struct {
	webhookSecret string
	db            *store.DB
	vcitaClient   *vcita.APIClient
	ai            AIService
	auditor       *audit.Logger
	log           *zap.Logger
}

// New constructs a Handler with all its dependencies.
func New(secret string, db *store.DB, vc *vcita.APIClient, ai AIService, auditor *audit.Logger, log *zap.Logger) *Handler {
	return &Handler{
		webhookSecret: secret,
		db:            db,
		vcitaClient:   vc,
		ai:            ai,
		auditor:       auditor,
		log:           log,
	}
}

// Handle is the Gin handler for POST /webhook.
// It reads the raw body for signature verification, then dispatches async.
func (h *Handler) Handle(c *gin.Context) {
	// Read raw body — Gin normally consumes it; we need it raw for HMAC.
	body, err := c.GetRawData()
	if err != nil {
		h.log.Warn("webhook: failed to read body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	// ── 1. Verify HMAC-SHA256 signature ──────────────────────────────────────
	if !h.verifySignature(c.GetHeader("X-Vcita-Signature"), body) {
		h.log.Warn("webhook: invalid signature – request rejected")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// ── 2. Parse event envelope ───────────────────────────────────────────────
	var envelope struct {
		EventType string          `json:"event_type"`
		Payload   json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed envelope"})
		return
	}

	// ── 3. Respond 200 immediately so inTandem does not retry ─────────────────
	c.JSON(http.StatusOK, gin.H{"received": true})

	// ── 4. Dispatch asynchronously under a fresh context ─────────────────────
	payload := envelope.Payload
	eventType := envelope.EventType
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.dispatch(ctx, eventType, payload)
	}()
}

// dispatch routes each event type to its dedicated handler.
func (h *Handler) dispatch(ctx context.Context, eventType string, payload json.RawMessage) {
	switch eventType {
	case "message.client_sent_message":
		h.handleClientMessage(ctx, payload)
	case "message.business_sent_message":
		// We sent this ourselves — audit and skip.
		h.auditor.Log("webhook_skipped", "", "system", "event=business_sent_message")
	case "client.updated":
		h.handleClientUpdated(ctx, payload)
	case "appointment.requested":
		h.handleAppointmentRequested(ctx, payload)
	default:
		h.log.Debug("webhook: unhandled event", zap.String("event_type", eventType))
	}
}

// ── Event handlers ────────────────────────────────────────────────────────────

func (h *Handler) handleClientMessage(ctx context.Context, raw json.RawMessage) {
	var msg vcita.MessagePayload
	if err := json.Unmarshal(raw, &msg); err != nil {
		h.log.Error("webhook: parse MessagePayload", zap.Error(err))
		return
	}

	h.auditor.Log("webhook_received", msg.ClientID, "system",
		fmt.Sprintf("event=client_sent_message channel=%s", msg.Channel))

	// Persist incoming patient message (encrypted)
	if err := h.db.SaveMessage(msg.ClientID, "patient", msg.Body, "client_sent_message", false); err != nil {
		h.log.Error("webhook: save patient message", zap.Error(err))
		return
	}

	// Fetch conversation history and vcita notes for AI context
	history, err := h.db.GetConversationHistory(msg.ClientID, 10)
	if err != nil {
		h.log.Error("webhook: get conversation history", zap.Error(err))
		return
	}

	notes, err := h.vcitaClient.GetClientNotes(ctx, msg.ClientID)
	if err != nil {
		h.log.Warn("webhook: could not fetch notes", zap.Error(err))
		notes = nil
	}

	// Pull latest medication info from refill schedule if present
	medInfo := ""
	reminders, _ := h.db.GetDueReminders()
	for _, r := range reminders {
		if r.ClientID == msg.ClientID {
			medInfo = r.Medication
			break
		}
	}

	// ── Call AI ───────────────────────────────────────────────────────────────
	reply, escalate, err := h.ai.ProcessMessage(ctx, msg.ClientID, history, notes, medInfo)
	if err != nil {
		h.log.Error("webhook: AI ProcessMessage failed", zap.Error(err))
		h.sendFallback(ctx, msg)
		return
	}

	h.auditor.Log("ai_response_generated", msg.ClientID, "ai",
		fmt.Sprintf("escalate=%v", escalate))

	if escalate {
		h.escalate(ctx, msg.ClientID, "AI flagged message for human review")
	}

	// Persist AI reply (encrypted) then send via vcita API
	if err := h.db.SaveMessage(msg.ClientID, "assistant", reply, "ai_response", escalate); err != nil {
		h.log.Error("webhook: save assistant reply", zap.Error(err))
	}

	if err := h.vcitaClient.SendMessage(ctx, vcita.SendMessageRequest{
		ClientID:       msg.ClientID,
		ConversationID: msg.ConversationID,
		Body:           reply,
		Channel:        msg.Channel,
	}); err != nil {
		h.log.Error("webhook: vcita SendMessage failed", zap.Error(err))
		return
	}

	h.auditor.Log("message_sent", msg.ClientID, "system", "channel="+msg.Channel)
}

func (h *Handler) handleClientUpdated(ctx context.Context, raw json.RawMessage) {
	var payload struct {
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		h.log.Error("webhook: parse client.updated", zap.Error(err))
		return
	}

	notes, err := h.vcitaClient.GetClientNotes(ctx, payload.ClientID)
	if err != nil {
		h.log.Warn("webhook: client.updated – get notes failed", zap.Error(err))
		return
	}

	for _, note := range notes {
		if !strings.EqualFold(note.NoteType, "current_medications") {
			continue
		}

		medInfo, supplyDays, err := h.ai.ExtractMedication(ctx, note.Content)
		if err != nil {
			h.log.Error("webhook: ExtractMedication failed", zap.Error(err))
			h.auditor.Log("medication_extraction_failed", payload.ClientID, "system", "note_id="+note.ID)
			continue
		}

		refillDate := calculateRefillDate(supplyDays)
		if err := h.db.UpsertRefillSchedule(payload.ClientID, medInfo, refillDate); err != nil {
			h.log.Error("webhook: upsert refill schedule", zap.Error(err))
			continue
		}

		h.auditor.Log("medication_updated", payload.ClientID, "system",
			fmt.Sprintf("note_id=%s refill_date=%s", note.ID, refillDate.Format("2006-01-02")))
	}
}

func (h *Handler) handleAppointmentRequested(ctx context.Context, raw json.RawMessage) {
	var appt vcita.AppointmentPayload
	if err := json.Unmarshal(raw, &appt); err != nil {
		h.log.Error("webhook: parse AppointmentPayload", zap.Error(err))
		return
	}

	h.auditor.Log("appointment_requested", appt.ClientID, "system",
		"appointment_id="+appt.AppointmentID)

	staff, err := h.vcitaClient.GetStaff(ctx)
	if err != nil {
		h.log.Error("webhook: get staff", zap.Error(err))
		h.escalate(ctx, appt.ClientID, "Could not fetch staff for appointment")
		return
	}

	from, to := time.Now(), time.Now().Add(7*24*time.Hour)
	var allSlots []vcita.TimeSlot
	for _, s := range staff {
		slots, err := h.vcitaClient.GetAvailableSlots(ctx, s.ID, from, to)
		if err != nil {
			h.log.Warn("webhook: get slots", zap.String("staff_id", s.ID), zap.Error(err))
			continue
		}
		allSlots = append(allSlots, slots...)
	}

	if len(allSlots) == 0 {
		h.escalate(ctx, appt.ClientID, "No available slots for appointment request")
		return
	}

	history, _ := h.db.GetConversationHistory(appt.ClientID, 5)
	suggestion, err := h.ai.SuggestAppointment(ctx, AppointmentRequest{
		ClientID:        appt.ClientID,
		RequestedTime:   appt.RequestedTime,
		ServiceName:     appt.ServiceName,
		AvailableSlots:  allSlots,
		AvailableStaff:  staff,
		ConversationCtx: history,
	})
	if err != nil {
		h.log.Error("webhook: SuggestAppointment failed", zap.Error(err))
		h.escalate(ctx, appt.ClientID, "AI could not suggest appointment slot")
		return
	}

	if suggestion.HumanInterventionNeeded {
		h.escalate(ctx, appt.ClientID, "AI flagged appointment for human review")
		return
	}

	sugID, err := h.db.SaveAppointmentSuggestion(
		appt.ClientID,
		suggestion.Staff.ID,
		suggestion.Slot.StartTime.Format(time.RFC3339),
	)
	if err != nil {
		h.log.Error("webhook: save appointment suggestion", zap.Error(err))
	}

	if err := h.vcitaClient.SendMessage(ctx, vcita.SendMessageRequest{
		ClientID: appt.ClientID,
		Body:     suggestion.ConfirmationMessage,
	}); err != nil {
		h.log.Error("webhook: send appointment confirmation", zap.Error(err))
		_ = h.db.UpdateAppointmentStatus(sugID, "rejected")
		return
	}

	_ = h.db.UpdateAppointmentStatus(sugID, "confirmed")
	h.auditor.Log("appointment_confirmed", appt.ClientID, "system",
		fmt.Sprintf("suggestion_id=%d", sugID))
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// verifySignature checks the HMAC-SHA256 webhook signature.
// inTandem sends: X-Vcita-Signature: sha256=<hex_digest>
func (h *Handler) verifySignature(sigHeader string, body []byte) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(sigHeader, prefix) {
		return false
	}
	expected := sigHeader[len(prefix):]
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	computed := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(computed), []byte(expected))
}

// escalate saves an escalation record and logs the event.
// reason must be a system-level description — never patient content.
func (h *Handler) escalate(ctx context.Context, clientID, reason string) {
	id, err := h.db.SaveEscalation(clientID, reason)
	if err != nil {
		h.log.Error("webhook: save escalation", zap.Error(err))
	}
	h.auditor.Log("escalation_created", clientID, "system",
		fmt.Sprintf("escalation_id=%d", id))
	h.log.Warn("escalation triggered",
		zap.String("client_id", clientID),
		zap.Uint("escalation_id", id),
	)
}

// sendFallback sends a safe holding message to the patient when AI fails.
func (h *Handler) sendFallback(ctx context.Context, msg vcita.MessagePayload) {
	const fallback = "Thank you for your message. A member of our team will be in touch with you shortly."
	if err := h.vcitaClient.SendMessage(ctx, vcita.SendMessageRequest{
		ClientID:       msg.ClientID,
		ConversationID: msg.ConversationID,
		Body:           fallback,
		Channel:        msg.Channel,
	}); err != nil {
		h.log.Error("webhook: send fallback failed", zap.Error(err))
	}
}

// calculateRefillDate returns the date supplyDays from now, minus a 7-day buffer.
func calculateRefillDate(supplyDays int) time.Time {
	const buffer = 7
	if supplyDays <= buffer {
		return time.Now().AddDate(0, 0, 1)
	}
	return time.Now().AddDate(0, 0, supplyDays-buffer)
}
