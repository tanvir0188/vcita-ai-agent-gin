package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"go.uber.org/zap"
)

type Message struct {
	UId             string `json:"uid"`
	Text            string `json:"text"`
	ConversationUID string `json:"conversation_uid"`

	Staff struct {
		Uid string `json:"uid"`
	} `json:"staff"`

	WasRead   bool   `json:"was_read"`
	Direction string `json:"direction"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

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
