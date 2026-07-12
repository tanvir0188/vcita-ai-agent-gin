package vcita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DoRequest performs an HTTP request without a body and decodes the JSON response.
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

// DoRequestWithBody performs an HTTP request with a JSON body and decodes the response.
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
