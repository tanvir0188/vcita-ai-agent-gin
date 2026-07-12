package vcita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"go.uber.org/zap"
)

// CreateMessage sends a business-to-client message via vcita API.
func CreateMessage(clientID string, text string) (string, error) {
	url := "https://api.vcita.biz/platform/v1/messages"

	requestBody := CreateMessageRequest{
		Message: MessagePayload{
			Direction: "business_to_client",
			ClientID:  clientID,
			Text:      text,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	bearerToken := "Bearer " + os.Getenv("VCITA_DIRECTORY_TOKEN")
	req.Header.Add("authorization", bearerToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf(
			"vcita api error: %s",
			string(body),
		)
	}

	return string(body), nil
}

// GetMessageHistory fetches the last 5 messages for a conversation.
func GetMessageHistory(conversationID string) ([]Message, error) {

	logger.Log.Info(
		"fetching message history",
		zap.String("conversation_id", conversationID),
	)

	url := fmt.Sprintf(
		"https://api.vcita.biz/v2/conversations/%s/messages?per_page=5",
		conversationID,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set(
		"authorization",
		"Bearer "+os.Getenv("VCITA_DIRECTORY_TOKEN"),
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {

		logger.Log.Error(
			"failed to fetch message history",
			zap.Error(err),
		)

		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {

		var body map[string]interface{}

		_ = json.NewDecoder(res.Body).Decode(&body)

		logger.Log.Error(
			"vcita api error",
			zap.Int("status_code", res.StatusCode),
			zap.Any("response", body),
		)

		return nil, fmt.Errorf(
			"vcita api returned status %d",
			res.StatusCode,
		)
	}

	var messages []Message

	err = json.NewDecoder(res.Body).Decode(&messages)
	if err != nil {

		logger.Log.Error(
			"failed to decode message history",
			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"message history fetched successfully",
		zap.Int("message_count", len(messages)),
	)

	return messages, nil
}

// GetLatestMessage returns the most recent message in a conversation.
func GetLatestMessage(conversationID string) Message {
	latestMessages, err := GetMessageHistory(conversationID)
	if err != nil {
		logger.Log.Error("failed to get message history: ")

		return Message{}
	}

	if len(latestMessages) == 0 {
		logger.Log.Info("no message history found")

		return Message{}
	}

	// assuming oldest -> newest
	latestMessage := latestMessages[0]

	return latestMessage
}

// GetSmartReplyEmail fetches an AI-generated smart reply from vcita.
func GetSmartReplyEmail(matterUID string, clientUID string) (string, error) {
	url := "https://api.vcita.biz/v3/ai/ai_smart_replies"

	requestBody := SmartReplyRequest{
		MatterUID: matterUID,
		ClientUID: clientUID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	bearerToken := "Bearer " + os.Getenv("VCITA_DIRECTORY_TOKEN")
	req.Header.Add("authorization", bearerToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf(
			"vcita api error (status %d): %s",
			res.StatusCode,
			string(body),
		)
	}

	// Unmarshal the successful response JSON into our structs
	var apiResponse SmartReplyResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Return only the email message string
	return apiResponse.Data.Payload.EmailMessage, nil
}
