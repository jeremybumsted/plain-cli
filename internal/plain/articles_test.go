package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHelpCenterArticleBySlug(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenterArticleBySlug": map[string]interface{}{
					"id":          "hca_001",
					"title":       "Getting Started",
					"contentHtml": "<p>Welcome</p>",
					"slug":        "getting-started",
					"status":      "PUBLISHED",
					"updatedAt":   map[string]string{"iso8601": "2025-06-01T00:00:00Z"},
					"articleGroup": map[string]interface{}{
						"id":   "grp_001",
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

	article, err := client.GetHelpCenterArticleBySlug("hc_001", "getting-started")
	if err != nil {
		t.Fatalf("GetHelpCenterArticleBySlug failed: %v", err)
	}

	if article.ID != "hca_001" {
		t.Errorf("Expected article ID 'hca_001', got '%s'", article.ID)
	}
	if article.Title != "Getting Started" {
		t.Errorf("Expected title 'Getting Started', got '%s'", article.Title)
	}
	if article.Slug != "getting-started" {
		t.Errorf("Expected slug 'getting-started', got '%s'", article.Slug)
	}
	if article.Status != "PUBLISHED" {
		t.Errorf("Expected status 'PUBLISHED', got '%s'", article.Status)
	}
	if article.Group == nil {
		t.Fatal("Expected article group to not be nil")
	}
	if article.Group.Name != "Basics" {
		t.Errorf("Expected group name 'Basics', got '%s'", article.Group.Name)
	}
}

func TestGetHelpCenterArticleBySlugNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenterArticleBySlug": nil,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	_, err := client.GetHelpCenterArticleBySlug("hc_001", "nonexistent")
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

func TestUpsertHelpCenterArticle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		variables := body["variables"].(map[string]interface{})
		input := variables["input"].(map[string]interface{})
		if input["helpCenterId"] != "hc_001" {
			t.Fatalf("expected helpCenterId hc_001, got %v", input["helpCenterId"])
		}
		if input["helpCenterArticleId"] != "hca_001" {
			t.Fatalf("expected helpCenterArticleId hca_001, got %v", input["helpCenterArticleId"])
		}
		if input["helpCenterArticleGroupId"] != "hcag_001" {
			t.Fatalf("expected helpCenterArticleGroupId hcag_001, got %v", input["helpCenterArticleGroupId"])
		}
		if input["status"] != "PUBLISHED" {
			t.Fatalf("expected status PUBLISHED, got %v", input["status"])
		}

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"upsertHelpCenterArticle": map[string]interface{}{
					"helpCenterArticle": map[string]interface{}{
						"id":          "hca_001",
						"title":       "Getting Started",
						"description": "Start here",
						"contentHtml": "<p>Welcome</p>",
						"slug":        "getting-started",
						"status":      "PUBLISHED",
						"createdAt":   map[string]string{"iso8601": "2025-05-01T00:00:00Z"},
						"updatedAt":   map[string]string{"iso8601": "2025-06-01T00:00:00Z"},
						"articleGroup": map[string]interface{}{
							"id":   "hcag_001",
							"name": "Basics",
							"slug": "basics",
						},
					},
					"error": nil,
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)
	article, err := client.UpsertHelpCenterArticle(HelpCenterArticleUpsertInput{
		HelpCenterID:             "hc_001",
		HelpCenterArticleID:      "hca_001",
		HelpCenterArticleGroupID: "hcag_001",
		Title:                    "Getting Started",
		Description:              "Start here",
		ContentHTML:              "<p>Welcome</p>",
		Status:                   "PUBLISHED",
	})
	if err != nil {
		t.Fatalf("UpsertHelpCenterArticle failed: %v", err)
	}

	if article.ID != "hca_001" || article.Slug != "getting-started" {
		t.Fatalf("unexpected article response: %+v", article)
	}
	if article.Group == nil || article.Group.ID != "hcag_001" {
		t.Fatalf("expected article group hcag_001, got %+v", article.Group)
	}
}

func TestUpsertHelpCenterArticleMutationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"upsertHelpCenterArticle": map[string]interface{}{
					"helpCenterArticle": nil,
					"error": map[string]interface{}{
						"message": "title is invalid",
						"code":    "input_validation",
					},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)
	_, err := client.UpsertHelpCenterArticle(HelpCenterArticleUpsertInput{
		HelpCenterID: "hc_001",
		Title:        "",
		ContentHTML:  "<p>Welcome</p>",
	})
	if err == nil {
		t.Fatal("expected mutation error, got nil")
	}
	if err.Error() != "title is invalid" {
		t.Fatalf("unexpected error: %v", err)
	}
}
