package medicationreminder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tanvir0188/vcita-ai-agent/internal/config"
)

func GetClientNotes(matterUID string) (*TrimmedNotesResponse, error) {
	url := fmt.Sprintf(
		"https://api.vcita.biz/v3/clients/client_notes?matter_uid=%s",
		matterUID,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	accessToken := "Bearer " + config.Envs.VcitaDirectoryToken

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var response ClientNotesResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return ConvertedNoteList(&response)
}

func ConvertedNoteList(response *ClientNotesResponse) (*TrimmedNotesResponse, error) {
	if len(response.Data.ClientNotes) == 0 {
		return nil, fmt.Errorf("no notes found")
	}

	result := &TrimmedNotesResponse{
		MatterUID: response.Data.ClientNotes[0].MatterUID,
	}

	for _, note := range response.Data.ClientNotes {
		result.TrimmedNotes = append(result.TrimmedNotes, TrimmedNote{
			NoteID:  note.UID,
			Content: ConvertHtmlToMarkDown(note.Content),
		})
	}

	return result, nil
}
