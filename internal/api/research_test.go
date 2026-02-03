package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestResearchResultsStruct tests that ResearchResults struct is properly defined
func TestResearchResultsStruct(t *testing.T) {
	// Test that we can create a ResearchResults struct
	results := ResearchResults{
		TaskID: "task-123",
		Status: "complete",
		Sources: []ResearchSource{
			{
				ID:    "source-1",
				Title: "Test Source",
				URL:   "https://example.com",
				Type:  1, // web
			},
		},
	}

	if results.TaskID != "task-123" {
		t.Errorf("expected TaskID 'task-123', got %s", results.TaskID)
	}

	if results.Status != "complete" {
		t.Errorf("expected Status 'complete', got %s", results.Status)
	}

	if len(results.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(results.Sources))
	}

	if results.Sources[0].ID != "source-1" {
		t.Errorf("expected source ID 'source-1', got %s", results.Sources[0].ID)
	}
}

// TestResearchSourceStruct tests that ResearchSource struct is properly defined
func TestResearchSourceStruct(t *testing.T) {
	source := ResearchSource{
		ID:    "src-abc",
		Title: "Sample Document",
		URL:   "https://example.com/doc",
		Type:  2, // doc type
	}

	if source.ID != "src-abc" {
		t.Errorf("expected ID 'src-abc', got %s", source.ID)
	}

	if source.Title != "Sample Document" {
		t.Errorf("expected Title 'Sample Document', got %s", source.Title)
	}

	if source.URL != "https://example.com/doc" {
		t.Errorf("expected URL 'https://example.com/doc', got %s", source.URL)
	}

	if source.Type != 2 {
		t.Errorf("expected Type 2, got %d", source.Type)
	}
}

// TestStartResearchMethodExists tests that the StartResearch method exists on Client
func TestStartResearchMethodExists(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a mock response with task ID
		response := `)]}'
[[["wrb.fr","Ljjv0c","[\"task-fast-123\"]",null,null,null,"generic"]]]`
		w.Write([]byte(response))
	}))
	defer server.Close()

	// Create client - we just need to verify the method exists
	client := New("test-token", "test-cookies")

	// Call StartResearch with deep=false (fast research)
	_, err := client.StartResearch("project-123", "test query", false)
	// We expect an error since we're not hitting a real server, but the method should exist
	if err == nil {
		t.Log("StartResearch executed without error")
	}

	// Call StartResearch with deep=true (deep research)
	_, err = client.StartResearch("project-123", "test query", true)
	if err == nil {
		t.Log("StartResearch (deep) executed without error")
	}
}

// TestPollResearchResultsMethodExists tests that the PollResearchResults method exists on Client
func TestPollResearchResultsMethodExists(t *testing.T) {
	// Create client
	client := New("test-token", "test-cookies")

	// Call PollResearchResults
	_, err := client.PollResearchResults("project-123")
	// We expect an error since we're not hitting a real server, but the method should exist
	if err == nil {
		t.Log("PollResearchResults executed without error")
	}
}

// TestImportResearchSourcesMethodExists tests that the ImportResearchSources method exists on Client
func TestImportResearchSourcesMethodExists(t *testing.T) {
	// Create client
	client := New("test-token", "test-cookies")

	// Call ImportResearchSources
	sources := []string{"https://example.com/1", "https://example.com/2"}
	err := client.ImportResearchSources("project-123", "task-456", sources)
	// We expect an error since we're not hitting a real server, but the method should exist
	if err == nil {
		t.Log("ImportResearchSources executed without error")
	}
}

// TestParseResearchResults tests parsing of research results from RPC response
func TestParseResearchResults(t *testing.T) {
	// Sample response format from poll research results
	// Format: [task_id, status, [[source_id, title, url, type], ...]]
	responseJSON := `["task-123", "complete", [["src-1", "Source One", "https://example.com/1", 1], ["src-2", "Source Two", "https://example.com/2", 2]]]`

	var data []interface{}
	err := json.Unmarshal([]byte(responseJSON), &data)
	if err != nil {
		t.Fatalf("failed to unmarshal test data: %v", err)
	}

	// Test the parsing logic that will be in the client
	results := parseResearchResultsFromData(data)
	if results == nil {
		t.Fatal("parseResearchResultsFromData returned nil")
	}

	if results.TaskID != "task-123" {
		t.Errorf("expected TaskID 'task-123', got %s", results.TaskID)
	}

	if results.Status != "complete" {
		t.Errorf("expected Status 'complete', got %s", results.Status)
	}

	if len(results.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(results.Sources))
	}

	if results.Sources[0].ID != "src-1" {
		t.Errorf("expected first source ID 'src-1', got %s", results.Sources[0].ID)
	}
}
