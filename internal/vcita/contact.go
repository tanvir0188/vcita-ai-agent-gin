package vcita

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GetClientDetail fetches a client's details from the vcita platform API.
func GetClientDetail(webhookSecret *string, contactId *string) (GetClientDetailResponse, error) {
	var response GetClientDetailResponse

	// Safeguard: Ensure required pointer arguments are not nil or empty
	if contactId == nil || *contactId == "" {
		return response, fmt.Errorf("contactId parameter is required")
	}
	if webhookSecret == nil || *webhookSecret == "" {
		return response, fmt.Errorf("webhookSecret parameter is required for authorization")
	}

	// Dynamic URL generation using the contactId pointer value
	url := "https://api.vcita.biz/platform/v1/clients/" + *contactId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}

	req.Header.Add("accept", "application/json")

	// Dereference the webhookSecret pointer to use as your Bearer token
	bearerToken := "Bearer " + *webhookSecret
	req.Header.Add("authorization", bearerToken)

	// Using a dedicated client with a timeout to prevent hanging connections
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	// Catch any API errors (like unauthorized 401s or bad requests)
	if res.StatusCode >= 400 {
		return response, fmt.Errorf("vcita api error (status %d): %s", res.StatusCode, string(body))
	}

	// Unmarshal the raw JSON body straight into your response struct structures
	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("failed to parse client details json: %w", err)
	}

	return response, nil
}

// GetClientNotes fetches notes for a conversation (matter) from the business API.
func GetClientNotes(conversationUID string) ([]string, error) {

	url := fmt.Sprintf(
		"https://api.vcita.biz/business/clients/v1/matters/%s/notes",
		conversationUID,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	token := "Bearer " + os.Getenv("VCITA_DIRECTORY_TOKEN")
	req.Header.Add(
		"authorization",
		token,
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request failed: %d - %s",
			res.StatusCode,
			string(body),
		)
	}

	var response NotesResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var notes []string

	for _, note := range response.Data.Notes {
		notes = append(notes, note.Content)
	}

	return notes, nil
}
