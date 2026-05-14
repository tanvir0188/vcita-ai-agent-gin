package vcita

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 15 * time.Second
	apiVersion     = "/platform/v1"
)

// APIClient is a thread-safe HTTP client for the inTandem platform API.
// It holds the base URL and auth token; all methods accept a context for cancellation.
type APIClient struct {
	baseURL     string
	bearerToken string
	http        *http.Client
}

// NewAPIClient creates a new inTandem API client.
// bearerToken is the business-level token from the platform.
func NewAPIClient(baseURL, bearerToken string) *APIClient {
	return &APIClient{
		baseURL:     baseURL,
		bearerToken: bearerToken,
		http: &http.Client{
			Timeout: defaultTimeout,
			// TLS 1.2+ is enforced by Go's default transport; no additional config needed.
		},
	}
}

// ── Messages ──────────────────────────────────────────────────────────────────

// SendMessage posts a reply into an existing patient conversation thread.
// Uses the same conversation_id from the incoming webhook to stay in-thread.
func (c *APIClient) SendMessage(ctx context.Context, req SendMessageRequest) error {
	// The inTandem v3 messaging endpoint – adjust path when you verify in API reference
	path := "/v3/communication/messages"
	return c.post(ctx, path, req, nil)
}

// ── Clients ───────────────────────────────────────────────────────────────────

// GetClient fetches a patient's profile by client_id.
func (c *APIClient) GetClient(ctx context.Context, clientID string) (*Client, error) {
	type respData struct {
		Client Client `json:"client"`
	}
	var resp APIResponse[respData]
	if err := c.get(ctx, fmt.Sprintf("%s/clients/%s", apiVersion, clientID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data.Client, nil
}

// GetClientNotes returns all notes for a patient.
func (c *APIClient) GetClientNotes(ctx context.Context, clientID string) ([]Note, error) {
	type respData struct {
		Notes []Note `json:"notes"`
	}
	var resp APIResponse[respData]
	if err := c.get(ctx, fmt.Sprintf("%s/clients/%s/notes", apiVersion, clientID), &resp); err != nil {
		return nil, err
	}
	return resp.Data.Notes, nil
}

// ── Staff & Availability ──────────────────────────────────────────────────────

// GetStaff returns the list of active staff members for the business.
func (c *APIClient) GetStaff(ctx context.Context) ([]StaffMember, error) {
	type respData struct {
		Staff []StaffMember `json:"staff"`
	}
	var resp APIResponse[respData]
	if err := c.get(ctx, fmt.Sprintf("%s/staffs", apiVersion), &resp); err != nil {
		return nil, err
	}
	return resp.Data.Staff, nil
}

// GetAvailableSlots fetches open scheduling slots.
// fromTime and toTime define the window to check.
func (c *APIClient) GetAvailableSlots(ctx context.Context, staffID string, from, to time.Time) ([]TimeSlot, error) {
	type respData struct {
		Slots []TimeSlot `json:"slots"`
	}
	var resp APIResponse[respData]
	path := fmt.Sprintf("%s/scheduling/slots?staff_id=%s&from=%s&to=%s",
		apiVersion, staffID, from.Format(time.RFC3339), to.Format(time.RFC3339))
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data.Slots, nil
}

// ScheduleAppointment creates a confirmed appointment via the API.
func (c *APIClient) ScheduleAppointment(ctx context.Context, req ScheduleAppointmentRequest) error {
	return c.post(ctx, fmt.Sprintf("%s/appointments", apiVersion), req, nil)
}

// ── HTTP helpers ──────────────────────────────────────────────────────────────

func (c *APIClient) get(ctx context.Context, path string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("vcita.get: build request: %w", err)
	}
	return c.do(req, out)
}

func (c *APIClient) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("vcita.post: marshal body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("vcita.post: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, out)
}

func (c *APIClient) do(req *http.Request, out interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("vcita: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MB limit
	if err != nil {
		return fmt.Errorf("vcita: reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("vcita: unexpected status %d: %s", resp.StatusCode, sanitiseBody(body))
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("vcita: unmarshal response: %w", err)
		}
	}
	return nil
}

// sanitiseBody truncates the response body for error messages.
// Never log the full body as it may contain PHI.
func sanitiseBody(b []byte) string {
	if len(b) > 200 {
		return string(b[:200]) + "...[truncated]"
	}
	return string(b)
}
