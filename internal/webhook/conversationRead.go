package webhook

import (
	"fmt"

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

type ConversationReadPayload struct {
	EntityName string `json:"entity_name"`
	EventType  string `json:"event_type"`
	Data       struct {
		MatterUID   string `json:"matter_uid"`
		UpdatedAt   string `json:"updated_at"`
		BusinessUID string `json:"business_uid"`
		ContactUID  string `json:"contact_uid"`
	} `json:"data"`
}

func (h *ConversationHandler) ConversationReadHandle(c *gin.Context) {

	h.log.Info("conversation read webhook received")

	var payload ConversationReadPayload

	if err := c.ShouldBindJSON(&payload); err != nil {

		h.log.Error("failed to parse conversation read webhook payload", zap.Error(err))

		h.auditor.Log(
			"conversation_read_webhook_parse_failed",
			"",
			"system",
			err.Error(),
		)

		c.JSON(400, gin.H{
			"success": false,
			"message": "invalid payload",
		})

		return
	}

	h.log.Info(
		"conversation read",
		zap.String("matter_uid", payload.Data.MatterUID),
		zap.String("contact_uid", payload.Data.ContactUID),
		zap.String("business_uid", payload.Data.BusinessUID),
		zap.String("updated_at", payload.Data.UpdatedAt),
	)

	h.auditor.Log(
		"conversation_read_webhook_received",
		payload.Data.MatterUID,
		"system",
		fmt.Sprintf(
			"conversation read by staff. contact_uid=%s business_uid=%s updated_at=%s",
			payload.Data.ContactUID,
			payload.Data.BusinessUID,
			payload.Data.UpdatedAt,
		),
	)

	// TODO:
	// add your business logic here
	// examples:
	// - mark unread messages as seen
	// - update conversation last_seen_at
	// - notify websocket clients
	// - sync CRM status
	// - store read event in database

	c.JSON(200, gin.H{
		"success": true,
		"message": "conversation read webhook processed",
	})
}
