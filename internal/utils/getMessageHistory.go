package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"go.uber.org/zap"
)

type Message struct {
	UId           string `json:"uid"`
	Text          string `json:"text"`
	MailDelivered bool   `json:"mail_delivered"`
	WasRead       bool   `json:"was_read"`
	Direction     string `json:"direction"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetMessageHistory(conversationID string) ([]Message, error) {
	logger.Log.Info(
		"fetching message history",
		zap.String("conversation_id", conversationID),
	)

	url := fmt.Sprintf(
		"https://api.vcita.biz/v2/conversations/%s/messages?last_update=true",
		conversationID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	token := "Bearer " + os.Getenv("VCITA_DIRECTORY_TOKEN")
	req.Header.Add("authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error(
			"get message history error",
			zap.Error(err),
		)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	logger.Log.Info(
		"message history raw response",
		zap.String("body", string(body)),
	)

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"vcita api error: %s",
			string(body),
		)
	}

	var messages []Message

	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
