//go:build e2e

package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// SearchResponse mirrors the structure from the orchestrator's proto file
// for unmarshalling the JSON response from the API Gateway.
type SearchResponse struct {
	Summary string `json:"summary"`
	Sources []struct {
		URL     string `json:"url"`
		Title   string `json:"title"`
		Snippet string `json:"snippet"`
	} `json:"sources"`
}

// TestE2ESearch performs a simple end-to-end test of the /search endpoint.
// It requires the entire Docker Compose stack to be running.
//
// To run this test:
// 1. In the project root, run `docker-compose up -d --build`.
// 2. Wait a few seconds for all services to start up.
// 3. Run the test with the e2e build tag: `go test -v -tags=e2e ./tests/...`
func TestE2ESearch(t *testing.T) {
	// A short delay to give services a moment to initialize.
	// In a more robust setup, we'd use a proper health check loop.
	time.Sleep(10 * time.Second)

	apiURL := "http://localhost:8080/search"
	query := `{"query": "what is gocolly?"}`

	reqBody := bytes.NewBuffer([]byte(query))
	resp, err := http.Post(apiURL, "application/json", reqBody)
	if err != nil {
		t.Fatalf("failed to send request to API gateway: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200 OK, got %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	expectedSummary := "This is a mock answer from the expert."
	if searchResp.Summary != expectedSummary {
		t.Errorf("expected summary '%s', got '%s'", expectedSummary, searchResp.Summary)
	}

	t.Log("End-to-end test passed successfully!")
}
