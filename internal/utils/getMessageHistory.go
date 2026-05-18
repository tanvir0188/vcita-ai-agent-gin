package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"go.uber.org/zap"
)

func GetMessageHistory(conversationID string) (string, error) {
	fmt.Println("Fetching message history for conversation:", conversationID)
	messageHistoryURL := "https://api.vcita.biz/v2/conversations/" + conversationID + "/messages"
	fmt.Println("Constructed message history URL:", messageHistoryURL)

	req, _ := http.NewRequest("GET", messageHistoryURL, nil)

	req.Header.Add("accept", "application/json")
	token := "Bearer " + os.Getenv("VCITA_DIRECTORY_TOKEN")
	req.Header.Add("authorization", token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		logger.Log.Error("Get message history error",
			zap.Error(err),
		)
		return "", err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	logger.Log.Info("message history fetched:", zap.String("history", string(body)))

	return string(body), nil
}
