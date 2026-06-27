package medicationreminder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/utils"
)

// reusable api request function
func DoRequest(method, url string, headers map[string]string, out any) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("api error: %s", string(body))
	}

	return json.NewDecoder(res.Body).Decode(out)
}

func DoRequestWithBody(method, url string, headers map[string]string, body any, out any) error {
	var reqBody io.Reader

	// Handle different body types
	if body != nil {
		switch v := body.(type) {
		case []byte:
			reqBody = bytes.NewReader(v)
		case string:
			reqBody = strings.NewReader(v)
		default:
			// Marshal to JSON for structs, maps, etc.
			jsonData, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(jsonData)

			// Auto-set Content-Type if not already set
			if headers == nil {
				headers = make(map[string]string)
			}
			if _, exists := headers["Content-Type"]; !exists {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	// Handle non-success status codes
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("api error [%d]: %s", res.StatusCode, string(bodyBytes))
	}

	// Decode response if 'out' is provided
	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// converts html to markdown
func ConvertHtmlToMarkDown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Printf("failed to convert html to markdown: %v", err)
		return html
	}

	return markdown
}

// returns note list in a format that will be feeded to ai
func ConvertedNoteList(response *ClientNotesResponse, mu *string) (*TrimmedNotesResponse, error) {
	if len(response.Data.ClientNotes) == 0 {
		return nil, fmt.Errorf("no notes found")
	}

	result := &TrimmedNotesResponse{
		MatterUID: *mu,
	}

	for _, note := range response.Data.ClientNotes {
		result.TrimmedNotes = append(result.TrimmedNotes, TrimmedNote{
			NoteID:  note.UID,
			Content: ConvertHtmlToMarkDown(note.Content),
		})
	}

	return result, nil
}

// Filters the note list to detect the note object with regex
func FilterTrimmedNotes(notes *TrimmedNotesResponse) *TrimmedNotesResponse {
	filtered := make([]TrimmedNote, 0)

	for _, note := range notes.TrimmedNotes {
		if HasCurrentMedications(note.Content) {
			filtered = append(filtered, note)
		}
	}

	return &TrimmedNotesResponse{
		MatterUID:    notes.MatterUID,
		TrimmedNotes: filtered,
	}
}

// regex to detect notes with current medication
func HasCurrentMedications(text string) bool {
	pattern := regexp.MustCompile(`\*\*Current Medications\*\*\n\n`)
	return pattern.MatchString(text)
}
func ExtractCurrentMedications(content string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return ""
	}

	var medicationSection *goquery.Selection

	doc.Find("strong").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) == "Current Medications" {
			medicationSection = s
			return false
		}
		return true
	})

	if medicationSection == nil {
		return ""
	}

	headerDiv := medicationSection.Closest("div")
	if headerDiv.Length() == 0 {
		return ""
	}

	var htmlSection strings.Builder

	for node := headerDiv.Next(); node.Length() > 0; node = node.Next() {
		html, err := goquery.OuterHtml(node)
		if err != nil {
			continue
		}

		htmlSection.WriteString(html)
	}

	return strings.TrimSpace(
		ConvertHtmlToMarkDown(htmlSection.String()),
	)
}

// determines the number of pages needs to be traversd for all the users
func DeterminePageAmount(itemCount int) int {
	if itemCount <= 0 {
		return 0
	}
	return (itemCount + 99) / 100
}

func FilterClients(
	clients []Client,
	filter func(Client) bool,
) []Client {
	filtered := make([]Client, 0)

	for _, client := range clients {
		if filter(client) {
			filtered = append(filtered, client)
		}
	}

	return filtered
}

func GetClientsWithCurrentMedications() ([]Client, error) {
	clients, err := GetAllClients()
	if err != nil {
		return nil, err
	}

	var result []Client

	for _, c := range clients {
		meds := ExtractCurrentMedications(c.Notes)

		if meds == "" {

			continue
		}

		c.Notes = meds
		if c.MatterUid == "2y1ncrj8uh7qk8ye" {
			res, err := ai.PredictMedicationRefill(c.Notes)
			if err != nil {
				return nil, err
			}
			resJson, _ := json.MarshalIndent(res, "", "  ")
			log.Printf("Prediction for %s:\n%s\n", c.MatterUid, string(resJson))

			_, err = utils.CreateMessage(c.UID, *res.Response)
			if err != nil {
				log.Printf("failed to create message: %v", err)
			}
		}
		result = append(result, c)
	}

	return result, nil
}

// GetAllClients fetches all clients across all available pages concurrently.
// It first requests the first page to determine the total number of pages,
// then spins up goroutines to fetch the remaining pages in parallel.
//
// Returns a slice of all Clients and an error if any page fetch fails.
func GetAllClients() ([]Client, error) {
	firstPage, err := GetClientList(1)
	if err != nil {
		return nil, err
	}

	totalPages := DeterminePageAmount(firstPage.Counts.Client)

	clients := make([]Client, 0, firstPage.Counts.Client)
	clients = append(clients, firstPage.TopHits.Client...)

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		errChan = make(chan error, totalPages)
	)

	// Start from page 2 because page 1 is already loaded
	for page := 2; page <= totalPages; page++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()

			resp, err := GetClientList(page)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			clients = append(clients, resp.TopHits.Client...)
			mu.Unlock()
		}(page)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return clients, nil
}
