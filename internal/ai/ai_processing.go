package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
