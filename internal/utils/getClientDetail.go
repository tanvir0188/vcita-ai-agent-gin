package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetClientDetailResponse struct {
	Status string           `json:"status"`
	Data   ClientDetailData `json:"data"`
}

type ClientDetailData struct {
	Client ClientInfo `json:"client"`
}

type ClientInfo struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	MobilePhone   string `json:"mobile_phone"`
}

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
