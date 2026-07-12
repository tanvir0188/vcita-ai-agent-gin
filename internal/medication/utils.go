package medication

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// ConvertHtmlToMarkDown converts HTML content to markdown.
func ConvertHtmlToMarkDown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Printf("failed to convert html to markdown: %v", err)
		return html
	}

	return markdown
}

// ConvertedNoteList returns note list in a format that will be fed to AI.
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

// FilterTrimmedNotes filters the note list to detect the note object with regex.
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

// HasCurrentMedications uses regex to detect notes with current medication.
func HasCurrentMedications(text string) bool {
	pattern := regexp.MustCompile(`\*\*Current Medications\*\*\n\n`)
	return pattern.MatchString(text)
}

// ExtractCurrentMedications extracts the current medications section from HTML content.
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

// DeterminePageAmount determines the number of pages needed to traverse all users.
func DeterminePageAmount(itemCount int) int {
	if itemCount <= 0 {
		return 0
	}
	return (itemCount + 99) / 100
}

// FilterClients filters clients based on a predicate function.
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

// GetClientsWithCurrentMedications returns all clients that have current medications.
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

			_, err = vcita.CreateMessage(c.UID, *res.Response)
			if err != nil {
				log.Printf("failed to create message: %v", err)
			}
		}
		result = append(result, c)
	}

	return result, nil
}

// GetAllClients fetches all clients across all available pages concurrently.
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
