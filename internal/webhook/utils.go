package webhook

import (
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
	"go.uber.org/zap"
)

func (h *Handler) processWebhook(
	envelope WebhookEnvelope,
) {

	payload := envelope.Data

	// Ignore unrelated contacts
	if payload.ContactUID != "06dodrl3k4w5k1rd" {

		h.log.Info(
			"ignoring webhook",
			zap.String("reason", "contact uid mismatch"),
		)

		return
	}

	// CUSTOMER MESSAGE
	if payload.Direction == "client_to_business" {

		err := h.db.
			CreateOrUpdateConversationOnCustomerMessageCreate(
				store.CustomerMessageCreateParams{
					ConversationID:        payload.ConversationUID,
					LastMessageID:         payload.UID,
					LastCustomerMessageID: payload.UID,
					LastMessageFrom:       "customer",
					AIReplyPending:        true,
					PendingMessageID:      payload.UID,
				},
			)

		if err != nil {

			h.log.Error(
				"failed to update customer conversation state",
				zap.Error(err),
			)

			return
		}

		h.log.Info(
			"customer conversation updated",
			zap.String("conversation_uid", payload.ConversationUID),
		)

		latestMessages, err := utils.GetMessageHistory(
			payload.ConversationUID,
		)
		if err != nil {

			h.log.Error(
				"failed to fetch message history",
				zap.Error(err),
			)

			return
		}

		for _, message := range latestMessages {

			h.log.Info(
				"message history item",
				zap.String("uid", message.UId),
				zap.String("text", message.Text),
				zap.String("direction", message.Direction),
			)
		}

		return
	}

	// STAFF MESSAGE
	if payload.Direction == "business_to_client" {

		err := h.db.
			CreateOrUpdateConversationOnStaffMessage(
				store.StaffMessageCreateParams{
					ConversationID:     payload.ConversationUID,
					LastMessageID:      payload.UID,
					LastStaffMessageID: payload.UID,
					LastMessageFrom:    "staff",
					HumanActive:        true,
					HumanActiveAt:      payload.CreatedAt,
					AIReplyPending:     false,
					AIReplyGenerating:  false,
					PendingMessageID:   "",
					
				},
			)

		if err != nil {

			h.log.Error(
				"failed to update staff conversation state",
				zap.Error(err),
			)

			return
		}

		h.log.Info(
			"staff conversation updated",
			zap.String("conversation_uid", payload.ConversationUID),
		)

		return
	}
}
