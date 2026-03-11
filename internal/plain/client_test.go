package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", client.token)
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestNewClientWithURL(t *testing.T) {
	customURL := "https://custom.example.com/api"
	client := NewClientWithURL("test-token", customURL)

	if client.baseURL != customURL {
		t.Errorf("Expected base URL '%s', got '%s'", customURL, client.baseURL)
	}
}

func TestGetToken(t *testing.T) {
	client := NewClient("my-secret-token")
	if client.GetToken() != "my-secret-token" {
		t.Errorf("Expected token 'my-secret-token', got '%s'", client.GetToken())
	}
}

func TestSetTimeout(t *testing.T) {
	client := NewClient("test-token")
	newTimeout := 60 * time.Second

	client.SetTimeout(newTimeout)

	if client.httpClient.Timeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, client.httpClient.Timeout)
	}
}

func TestRequest_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", authHeader)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type 'application/json'")
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
		}); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClientWithURL("test-token", server.URL)

	// Make request
	var result map[string]string
	err := client.request("GET", "", nil, &result)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", result["status"])
	}
}

func TestRequest_Error(t *testing.T) {
	// Create test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Unauthorized",
		}); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	var result map[string]string
	err := client.request("GET", "", nil, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}

	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
	}

	if apiErr.Message != "Unauthorized" {
		t.Errorf("Expected message 'Unauthorized', got '%s'", apiErr.Message)
	}
}

func TestRequest_WithBody(t *testing.T) {
	// Create test server
	receivedBody := make(map[string]interface{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		}); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	requestBody := map[string]string{
		"key": "value",
	}

	var result map[string]string
	err := client.request("POST", "", requestBody, &result)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if receivedBody["key"] != "value" {
		t.Errorf("Expected request body key='value', got '%v'", receivedBody["key"])
	}
}

func TestStubMethods(t *testing.T) {
	client := NewClient("test-token")

	// Test that deferred methods return "not implemented" errors
	// Phase 1-2: Read operations implemented (ListThreads, GetThread, GetMyThreads, SearchThreads)
	// Phase 3: Write operations implemented (ChangeThreadStatus, SnoozeThread, AssignThread, ChangeThreadPriority, CreateNote)
	// Phase 4: Label and field operations implemented (ListLabelTypes, AddLabels, RemoveLabels, ListThreadFieldSchemas)
	// Still deferred: ReplyToThread (customer-facing operation)
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "ReplyToThread",
			fn: func() error {
				return client.ReplyToThread("123", "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Errorf("%s should return error (not implemented)", tt.name)
			}
			if !strings.Contains(err.Error(), "not implemented") {
				t.Errorf("%s should return 'not implemented' error, got: %v", tt.name, err)
			}
		})
	}
}
