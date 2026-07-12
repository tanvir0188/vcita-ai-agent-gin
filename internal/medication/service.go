package medication

import (
	"fmt"
	"net/http"

	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// GetClientNotes fetches client notes by matter UID via vcita v3 API.
func GetClientNotes(matterUID string) (*TrimmedNotesResponse, error) {
	url := fmt.Sprintf(
		"https://api.vcita.biz/v3/clients/client_notes?matter_uid=%s",
		matterUID,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaDirectoryToken,
	}

	var response ClientNotesResponse

	if err := vcita.DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return ConvertedNoteList(&response, &matterUID)
}

// GetClientNote fetches a single note by its UID.
func GetClientNote(noteUID string) (*TrimmedNote, error) {
	url := fmt.Sprintf(
		"https://api.vcita.biz/v3/clients/client_notes/%s",
		noteUID,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaDirectoryToken,
	}

	var response ClientNoteResponse

	if err := vcita.DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return &TrimmedNote{
		NoteID:  response.Data.UID,
		Content: ConvertHtmlToMarkDown(response.Data.Content),
	}, nil
}

// GetClientList fetches a paginated list of clients via vcita search API.
func GetClientList(page int) (*ClientListResponse, error) {
	url := fmt.Sprintf(
		"https://api.vcita.biz/v2/search?entity=client&entities=client&page=%d&per_page=100&search_filter[tags_relation]=or",
		page,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaDirectoryToken,
	}

	var response ClientListResponse

	if err := vcita.DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
