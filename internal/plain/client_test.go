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

func TestListHelpCenters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenters": map[string]interface{}{
					"edges": []map[string]interface{}{
						{
							"node": map[string]interface{}{
								"id": "hc_123",
							},
						},
						{
							"node": map[string]interface{}{
								"id": "hc_456",
							},
						},
					},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	helpCenters, err := client.ListHelpCenters()
	if err != nil {
		t.Fatalf("ListHelpCenters failed: %v", err)
	}

	if len(helpCenters) != 2 {
		t.Fatalf("Expected 2 help centers, got %d", len(helpCenters))
	}

	if helpCenters[0].ID != "hc_123" {
		t.Errorf("Expected first help center ID 'hc_123', got '%s'", helpCenters[0].ID)
	}

	if helpCenters[1].ID != "hc_456" {
		t.Errorf("Expected second help center ID 'hc_456', got '%s'", helpCenters[1].ID)
	}
}

func TestGetHelpCenterArticle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenterArticle": map[string]interface{}{
					"id":          "hca_123",
					"title":       "Getting Started",
					"contentHtml": "<p>Welcome to our help center</p>",
					"slug":        "getting-started",
					"status":      "PUBLISHED",
					"updatedAt": map[string]interface{}{
						"iso8601": "2024-03-09T10:30:00Z",
					},
					"articleGroup": map[string]interface{}{
						"id":   "hcag_456",
						"name": "Basics",
					},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	article, err := client.GetHelpCenterArticle("hca_123")
	if err != nil {
		t.Fatalf("GetHelpCenterArticle failed: %v", err)
	}

	if article.ID != "hca_123" {
		t.Errorf("Expected article ID 'hca_123', got '%s'", article.ID)
	}

	if article.Title != "Getting Started" {
		t.Errorf("Expected article title 'Getting Started', got '%s'", article.Title)
	}

	if article.ContentHTML != "<p>Welcome to our help center</p>" {
		t.Errorf("Expected article content '<p>Welcome to our help center</p>', got '%s'", article.ContentHTML)
	}

	if article.Slug != "getting-started" {
		t.Errorf("Expected article slug 'getting-started', got '%s'", article.Slug)
	}

	if article.Status != "PUBLISHED" {
		t.Errorf("Expected article status 'PUBLISHED', got '%s'", article.Status)
	}

	if article.Group == nil {
		t.Fatal("Expected article group to not be nil")
	}

	if article.Group.ID != "hcag_456" {
		t.Errorf("Expected article group ID 'hcag_456', got '%s'", article.Group.ID)
	}

	if article.Group.Name != "Basics" {
		t.Errorf("Expected article group name 'Basics', got '%s'", article.Group.Name)
	}
}

func TestGetHelpCenterArticle_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenterArticle": nil,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	_, err := client.GetHelpCenterArticle("hca_nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent article, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestListHelpCenterArticles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenter": map[string]interface{}{
					"articles": map[string]interface{}{
						"edges": []map[string]interface{}{
							{
								"node": map[string]interface{}{
									"id":     "hca_123",
									"title":  "Article 1",
									"slug":   "article-1",
									"status": "PUBLISHED",
									"updatedAt": map[string]interface{}{
										"iso8601": "2024-03-09T10:30:00Z",
									},
									"articleGroup": map[string]interface{}{
										"id":   "hcag_456",
										"name": "Group A",
									},
								},
							},
							{
								"node": map[string]interface{}{
									"id":     "hca_789",
									"title":  "Article 2",
									"slug":   "article-2",
									"status": "DRAFT",
									"updatedAt": map[string]interface{}{
										"iso8601": "2024-03-08T09:00:00Z",
									},
									"articleGroup": nil,
								},
							},
						},
						"pageInfo": map[string]interface{}{
							"hasNextPage": false,
							"endCursor":   "",
						},
					},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	articles, err := client.ListHelpCenterArticles("hc_123", false)
	if err != nil {
		t.Fatalf("ListHelpCenterArticles failed: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(articles))
	}

	if articles[0].ID != "hca_123" {
		t.Errorf("Expected first article ID 'hca_123', got '%s'", articles[0].ID)
	}

	if articles[0].Title != "Article 1" {
		t.Errorf("Expected first article title 'Article 1', got '%s'", articles[0].Title)
	}

	if articles[1].ID != "hca_789" {
		t.Errorf("Expected second article ID 'hca_789', got '%s'", articles[1].ID)
	}

	if articles[1].Group != nil {
		t.Error("Expected second article group to be nil")
	}
}

func TestListHelpCenterArticles_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenter": nil,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	_, err := client.ListHelpCenterArticles("hc_nonexistent", false)
	if err == nil {
		t.Fatal("Expected error for non-existent help center, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestListWorkspaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"myWorkspace": map[string]interface{}{
					"id":   "ws_123",
					"name": "My Workspace",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	workspaces, err := client.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces failed: %v", err)
	}

	if len(workspaces) != 1 {
		t.Fatalf("Expected 1 workspace, got %d", len(workspaces))
	}

	if workspaces[0].ID != "ws_123" {
		t.Errorf("Expected workspace ID 'ws_123', got '%s'", workspaces[0].ID)
	}

	if workspaces[0].Name != "My Workspace" {
		t.Errorf("Expected workspace name 'My Workspace', got '%s'", workspaces[0].Name)
	}
}

func TestListWorkspaces_Error(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
		errMsg   string
	}{
		{
			name: "GraphQL error",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"myWorkspace": nil,
				},
				"errors": []map[string]interface{}{
					{
						"message": "Authentication required",
					},
				},
			},
			errMsg: "GraphQL error: Authentication required",
		},
		{
			name: "Nil myWorkspace",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"myWorkspace": nil,
				},
			},
			errMsg: "unable to fetch workspace: myWorkspace is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(tt.response); err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := NewClientWithURL("test-token", server.URL)

			_, err := client.ListWorkspaces()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err.Error() != tt.errMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}
