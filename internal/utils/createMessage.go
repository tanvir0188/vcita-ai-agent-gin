package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CreateMessageRequest struct {
	Message MessagePayload `json:"message"`
}

type MessagePayload struct {
	Direction string `json:"direction"`
	ClientID  string `json:"client_id"`
	Text      string `json:"text"`
}

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
