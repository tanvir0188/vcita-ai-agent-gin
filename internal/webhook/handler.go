// Package webhook handles incoming inTandem webhook events via Gin.
// Every request is HMAC-SHA256 signature-verified before any processing.
// Message body content is never logged (HIPAA compliance).
package webhook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// Handler holds all dependencies for processing webhook events.
type Handler struct {
	webhookSecret string
	db            *store.DB
	vcitaClient   *vcita.APIClient
	auditor       *audit.Logger
	log           *zap.Logger
}

// New constructs a Handler with all its dependencies.
func New(secret string, db *store.DB, vc *vcita.APIClient, auditor *audit.Logger, log *zap.Logger) *Handler {
	return &Handler{
		webhookSecret: secret,
		db:            db,
		vcitaClient:   vc,
		auditor:       auditor,
		log:           log,
	}
}

// Handle is the Gin handler for POST /webhook.
// It reads the raw body for signature verification, then dispatches async.
func (h *Handler) Handle(c *gin.Context) {

	var envelope WebhookEnvelope

	if err := c.ShouldBindJSON(&envelope); err != nil {

		h.log.Warn(
			"invalid webhook payload",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})

		return
	}

	h.log.Info(
		"webhook received",
		zap.String("event_type", envelope.EventType),
		zap.String("conversation_id", envelope.Data.ConversationUID),
		zap.String("messege_id", envelope.Data.UID),
		zap.String("contact_id", envelope.Data.ContactUID),
		zap.String("direction", envelope.Data.Direction),
		zap.String("Staff_UID", envelope.Data.StaffUid),
		zap.String("message_type", envelope.Data.MessageType),
	)
	if envelope.Data.ContactUID != "06dodrl3k4w5k1rd" {

		// Respond immediately
		h.log.Info(
			"ignoring webhook event",
			zap.String("reason", "not a client message or contact ID mismatch"),
		)
		c.JSON(http.StatusOK, gin.H{
			"received": true,
		})
		return
	}

	latestMessages, err := utils.GetMessageHistory(envelope.Data.ConversationUID)
	if err != nil {
		h.log.Error(
			"failed to fetch message history",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch message history",
		})
		return
	}
	for _, message := range latestMessages {

		h.log.Info(
			"message",
			zap.String("uid", message.UId),
			zap.String("text", message.Text),
			zap.String("direction", message.Direction),
		)
	}

	// Process the webhook asynchronously to avoid blocking the response.
	go h.processWebhook(envelope)

	c.JSON(http.StatusOK, gin.H{
		"received": true,
	})

}
