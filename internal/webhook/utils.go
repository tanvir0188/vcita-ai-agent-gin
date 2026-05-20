package webhook

import (
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
	"go.uber.org/zap"
)

func (h *Handler) processWebhook(
	envelope WebhookEnvelope,
) {

	payload := envelope.Data

	// Ignore unrelated contacts
	if payload.ContactUID != "06dodrl3k4w5k1rd" || payload.MessageType != "text" {

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

		// get last message id and direction
		latestMessage := utils.GetLatestMessage(payload.ConversationUID)

		if latestMessage.Direction != "client_to_business" {

			h.log.Info(
				"latest message not customer, skipping ai evaluation",
			)

			return
		}

		if latestMessage.UId != payload.UID {

			h.log.Info(
				"payload message is no longer latest",
				zap.String("payload_uid", payload.UID),
				zap.String("latest_uid", latestMessage.UId),
			)

			return
		}

		state, err := h.db.GetConversationByID(
			payload.ConversationUID,
		)

		if err != nil {

			h.log.Error(
				"failed to fetch conversation state",
				zap.Error(err),
			)

			return
		}

		go WaitForEvaluation(
			h.db,
			&AiEvaluationParam{
				ConversationID:      payload.ConversationUID,
				ConversationVersion: state.ConversationVersion,
				PendingMessageID:    payload.UID,
				ContactId:           payload.ContactUID,
			},
		)
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
					HumanActiveUntil:   time.Now().Add(10 * time.Minute),
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

func WaitForEvaluation(db *store.DB, params *AiEvaluationParam) {
	logger.Log.Info("Starting AI evaluation cooldown",
		zap.String("conversation_id", params.ConversationID))

	// Cooldown
	time.Sleep(10 * time.Second)

	for attempt := 0; attempt < 3; attempt++ { // retry on transient DB issues
		state, err := db.GetConversationByID(params.ConversationID)
		if err != nil {
			logger.Log.Error("failed to get conversation", zap.Error(err))
			return
		}

		// === Validation ===
		if state.ConversationVersion != params.ConversationVersion {
			logger.Log.Info("conversation version changed, aborting")
			return
		}
		if state.LastMessageID != params.PendingMessageID {
			logger.Log.Info("latest message changed, aborting")
			return
		}
		if state.HumanActive && time.Now().Before(state.HumanActiveUntil) {
			logger.Log.Info("human is still active")
			return
		}
		if state.AIReplyGenerating {
			logger.Log.Info("another AI generation already in progress")
			return
		}

		// Mark as generating
		state.AIReplyGenerating = true
		if err := db.SaveConversation(state); err != nil {
			logger.Log.Error("failed to set AIReplyGenerating", zap.Error(err))
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	// Simulate / Call real AI
	reply_text, err := ai.GetSmartReplyEmail(params.ConversationID, params.ContactId)
	generatedReply := reply_text

	// === Final validation before sending ===
	state, err := db.GetConversationByID(params.ConversationID)
	if err != nil ||
		state.ConversationVersion != params.ConversationVersion ||
		state.HumanActive && time.Now().Before(state.HumanActiveUntil) ||
		state.LastMessageID != params.PendingMessageID {

		logger.Log.Info("state changed before sending, aborting")
		// Optionally reset AIReplyGenerating
		return
	}

	// Send message
	replied, err := utils.CreateMessage(params.ContactId, generatedReply)
	if err != nil {
		logger.Log.Error("failed to send AI reply", zap.Error(err))
		return
	}

	logger.Log.Info("AI reply sent successfully", zap.String("reply_id", replied))

	// Update state
	state.AIReplyPending = false
	state.AIReplyGenerating = false
	state.ConversationVersion++
	state.LastMessageID = replied // important!
	// Optionally reset HumanActive if you want AI to take over again after replying

	if err := db.SaveConversation(state); err != nil {
		logger.Log.Error("failed to update conversation after AI reply", zap.Error(err))
	}
}
