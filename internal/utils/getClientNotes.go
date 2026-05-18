package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// here conversation uid is actually mater uid
type NotesResponse struct {
	Data struct {
		Notes []struct {
			Content string `json:"content"`
		} `json:"notes"`
	} `json:"data"`
}

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
