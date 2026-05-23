package webhook

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) processWebhook(
	envelope WebhookEnvelope,
) {

	payload := envelope.Data

	// Ignore unrelated contacts

	// Querying the specific column using the contactid string
	state, err := h.db.GetConversationByID(payload.ConversationUID)

	if err != nil {
		logger.Log.Error("failded to get conversation state")
	}
	if state.AutoReplyOffUntil != nil && time.Now().Before(*state.AutoReplyOffUntil) {
		h.log.Info("AI replies are disabled for this conversation due to escalation",
			zap.String("conversation_id", payload.ConversationUID))
		return
	}

	if err != nil {
		logger.Log.Error("failed to get AutoReplyOffUntil timestamp", zap.Error(err))
		return
	}
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

		// latestMessages, err := utils.GetMessageHistory(
		// 	payload.ConversationUID,
		// )
		// if err != nil {

		// 	h.log.Error(
		// 		"failed to fetch message history",
		// 		zap.Error(err),
		// 	)

		// 	return
		// }

		// for _, message := range latestMessages {

		// 	h.log.Info(
		// 		"message history item",
		// 		zap.String("uid", message.UId),
		// 		zap.String("text", message.Text),
		// 		zap.String("direction", message.Direction),
		// 	)
		// }

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
			h.db, &h.webhookSecret, &h.slackWebhookUrl,
			&AiEvaluationParam{
				ConversationID:      payload.ConversationUID,
				ConversationVersion: state.ConversationVersion,
				PendingMessageID:    payload.UID,
				ContactId:           payload.ContactUID,
				Text:                payload.Text,
				AssignedStaffEmail:  latestMessage.Staff.Email,
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

type HumanInterventionResponse struct {
	HumanInterventionNeeded bool   `json:"human_intervention_needed"`
	EscalationReason        string `json:"escalation_reason"`
}

func WaitForEvaluation(db *store.DB, webhookSecret *string, whUrl *string, params *AiEvaluationParam) {
	logger.Log.Info("Starting AI evaluation cooldown",
		zap.String("conversation_id", params.ConversationID))

	// Cooldown
	result, err := ai.HumanInterventionNeeded(
		params.Text,
	)

	if err != nil {
		logger.Log.Error("openai_error", zap.Error(err))
	}
	// === Human Intervention / Slack Escalation Trigger ===
	// Moving this here allows us to send the actual generated text to Slack!
	var parsed HumanInterventionResponse
	cleaned := strings.TrimSpace(result)

	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	err = json.Unmarshal([]byte(cleaned), &parsed)

	if err != nil {
		logger.Log.Error(
			"failed to parse AI response",
			zap.Error(err),
			zap.String("raw_response", result),
		)
		return
	}

	humanInterventionNeeded := parsed.HumanInterventionNeeded
	escalationReason := parsed.EscalationReason

	if humanInterventionNeeded {
		logger.Log.Warn("Human intervention required - escalating",
			zap.String("conversation_id", params.ConversationID))

		clientDetails, err := utils.GetClientDetail(webhookSecret, &params.ContactId)
		if err != nil {
			logger.Log.Error("failed to fetch client details for slack alert", zap.Error(err))
		}
		fetchedClient := clientDetails.Data.Client

		// 2. Build the configuration payload
		slackConfig := utils.SlackClient{
			WebhookURL: *whUrl,         // Dereference *string to get the raw string URL
			Text:       params.Text,    // Pass the generated text draft we just got from AI
			Client:     &fetchedClient, // Pass the address of the ClientInfo struct

		}

		// 3. Fire the custom block kit formatter function we built
		err = db.GetGorm().
			Model(&store.Conversation{}).
			Where("conversation_id = ?", params.ConversationID).
			Updates(map[string]interface{}{
				"has_escalated":        true,
				"auto_reply_off_until": time.Now().Add(20 * time.Second),
				"conversation_version": gorm.Expr("conversation_version + 1"),
			}).Error

		if err != nil {
			logger.Log.Error(
				"failed to update escalation state",
				zap.Error(err),
			)
			return
		}
		utils.SendMessageToSlack(slackConfig, params.AssignedStaffEmail, params.ContactId, escalationReason)

	}

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

	// Call real AI (Note: updated arguments to match our dynamic smart reply function)
	reply_text, err := ai.GetSmartReplyEmail(params.ConversationID, params.ContactId)
	if err != nil {
		logger.Log.Error("failed to generate smart reply", zap.Error(err))
		return
	}
	generatedReply := reply_text

	// === Final validation before sending ===
	state, err := db.GetConversationByID(params.ConversationID)
	if err != nil ||
		state.ConversationVersion != params.ConversationVersion ||
		state.HumanActive && time.Now().Before(state.HumanActiveUntil) ||
		state.LastMessageID != params.PendingMessageID {

		logger.Log.Info("state changed before sending, aborting")
		return
	}

	// Send message via API
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
	state.LastMessageID = replied

	if err := db.SaveConversation(state); err != nil {
		logger.Log.Error("failed to update conversation after AI reply", zap.Error(err))
	}
}
