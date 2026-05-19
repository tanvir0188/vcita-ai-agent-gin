package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"go.uber.org/zap"
)

type ConversationHandler struct {
	webhookSecret string
	db            *store.DB

	auditor *audit.Logger
	log     *zap.Logger
}

// New constructs a Handler with all its dependencies.
func NewConversation(secret string, db *store.DB, auditor *audit.Logger, log *zap.Logger) *ConversationHandler {
	return &ConversationHandler{
		webhookSecret: secret,
		db:            db,
		auditor:       auditor,
		log:           log,
	}
}
func (h *ConversationHandler) ConversationReadHandle(c *gin.Context) {

	h.log.Info("conversation read webhook received")

	h.auditor.Log(
		"conversation_read_webhook_received",
		"",
		"system",
		"conversation read webhook received",
	)

	c.JSON(200, gin.H{
		"success": true,
	})
}
