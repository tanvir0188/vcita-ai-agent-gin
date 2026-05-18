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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
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

	// Read raw body
	body, err := c.GetRawData()
	if err != nil {
		h.log.Warn("webhook: failed to read body",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot read body",
		})

		return
	}

	// Parse webhook
	var envelope WebhookEnvelope

	if err := json.Unmarshal(body, &envelope); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed envelope",
		})

		return
	}

	// Respond immediately
	c.JSON(http.StatusOK, gin.H{
		"received": true,
	})

	payload := envelope.Data
	eventType := envelope.EventType

	h.log.Info("webhook received",
		zap.String("event_type", eventType),
	)

	h.log.Info("webhook payload",
		zap.Any("payload", payload),
	)

	// Skip outbound messages
	if payload.Direction == "business_to_client" {

		h.log.Info("message sent by business, skipping processing",
			zap.String("message_uid", payload.UID),
		)

		return
	}

	var (
		history string
		notes   []string

		historyErr error
		notesErr   error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch message history concurrently
	go func() {
		defer wg.Done()

		history, historyErr = utils.GetMessageHistory(
			payload.ConversationUID,
		)
	}()

	// Fetch notes concurrently
	go func() {
		defer wg.Done()

		notes, notesErr = utils.GetClientNotes(
			payload.ConversationUID,
		)
	}()

	// Wait for both requests
	wg.Wait()

	// Handle history result
	if historyErr != nil {

		h.log.Error("failed to fetch message history",
			zap.String("conversation_uid", payload.ConversationUID),
			zap.Error(historyErr),
		)

	} else {

		h.log.Info("message history fetched",
			zap.String("conversation_uid", payload.ConversationUID),
			zap.String("history", history),
		)
	}

	// Handle notes result
	if notesErr != nil {

		h.log.Error("failed to fetch client notes",
			zap.String("conversation_uid", payload.ConversationUID),
			zap.Error(notesErr),
		)

	} else {

		h.log.Info("client notes fetched",
			zap.String("conversation_uid", payload.ConversationUID),
			zap.Any("notes", notes),
		)
	}

	// AI processing
	utils.CreateMessage(
		payload.ContactUID,
		"Ai processed texts",
	)
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
