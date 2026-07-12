package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
)

// SendMessageToSlack sends a beautifully formatted Block Kit layout to Slack
func SendMessageToSlack(client SlackClient, assignedStaffEmail string, contactID string, reasoning string) {
	if client.Client == nil {
		logger.Log.Error("SlackClient missing client info dataset")
		return
	}

	// 1. Construct the fallback text for notifications
	fallbackText := fmt.Sprintf("Escalation for Client: %s %s", client.Client.FirstName, client.Client.LastName)

	// 2. Build the structural Block Kit payload layout array
	payload := SlackBlockPayload{
		Text: fallbackText,
		Blocks: []SlackBlock{
			// Header Block
			{
				Type: "header",
				Text: &SlackText{
					Type: "plain_text",
					Text: "Client Escalation Alert",
				},
			},
			// Client Identification Fields Block
			{
				Type: "section",
				Fields: []SlackText{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Assigned staff email:*\n%s", assignedStaffEmail),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Client Name:*\n%s %s", client.Client.FirstName, client.Client.LastName),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Contact ID:*\n`%s`", contactID),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Email:*\n%s", client.Client.Email),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Phone Number:*\n%s", client.Client.MobilePhone),
					},
				},
			},
			// Divider context for better reading flow
			{
				Type: "divider",
			},
			// Escalation Analysis / Core Text Block
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Reason for Escalation:*\n%s", reasoning),
				},
			},
			// Underlying Smart Reply Text Context Block
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Message content:*\n_%s_", client.Text),
				},
			},
		},
	}

	// 3. Unmarshal validation and transmission block
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Error("Failed to marshal slack block kit layout: " + err.Error())
		return
	}

	// Using a customized client with a programmatic timeout instead of raw default http.Post
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.Post(client.WebhookURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Log.Error("Failed to connect to Slack API webhook endpoint: " + err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(res.Body)
		logger.Log.Error(fmt.Sprintf("Slack endpoint rejection status (%d): %s", res.StatusCode, string(bodyBytes)))
	}
}
