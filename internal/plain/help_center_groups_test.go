package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHelpCenterArticleGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenterArticleGroup": map[string]interface{}{
					"id":        "hcag_001",
					"name":      "Basics",
					"slug":      "basics",
					"createdAt": map[string]string{"iso8601": "2025-05-01T00:00:00Z"},
					"updatedAt": map[string]string{"iso8601": "2025-06-01T00:00:00Z"},
					"parentArticleGroup": map[string]interface{}{
						"id":   "hcag_parent",
						"name": "Docs",
						"slug": "docs",
					},
					"childArticleGroups": map[string]interface{}{
						"edges": []map[string]interface{}{{
							"node": map[string]interface{}{"id": "hcag_child", "name": "Advanced", "slug": "advanced"},
						}},
					},
					"articles": map[string]interface{}{
						"edges": []map[string]interface{}{{
							"node": map[string]interface{}{
								"id":          "hca_001",
								"title":       "Getting Started",
								"description": "Start here",
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
						}},
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

	group, err := NewClientWithURL("test-token", server.URL).GetHelpCenterArticleGroup("hcag_001")
	if err != nil {
		t.Fatalf("GetHelpCenterArticleGroup failed: %v", err)
	}
	if group.ID != "hcag_001" || group.Name != "Basics" {
		t.Fatalf("unexpected group: %+v", group)
	}
	if len(group.ChildArticleGroups) != 1 || len(group.Articles) != 1 {
		t.Fatalf("unexpected children/articles: %+v", group)
	}
}

func TestCreateHelpCenterArticleGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		input := body["variables"].(map[string]interface{})["input"].(map[string]interface{})
		if input["helpCenterId"] != "hc_001" || input["parentHelpCenterArticleGroupId"] != "hcag_parent" {
			t.Fatalf("unexpected input: %+v", input)
		}

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"createHelpCenterArticleGroup": map[string]interface{}{
					"helpCenterArticleGroup": map[string]interface{}{
						"id":        "hcag_001",
						"name":      "Basics",
						"slug":      "basics",
						"createdAt": map[string]string{"iso8601": "2025-05-01T00:00:00Z"},
						"updatedAt": map[string]string{"iso8601": "2025-06-01T00:00:00Z"},
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

	group, err := NewClientWithURL("test-token", server.URL).CreateHelpCenterArticleGroup("hc_001", "Basics", "hcag_parent")
	if err != nil {
		t.Fatalf("CreateHelpCenterArticleGroup failed: %v", err)
	}
	if group.ID != "hcag_001" || group.Slug != "basics" {
		t.Fatalf("unexpected group: %+v", group)
	}
}

func TestUpdateHelpCenterArticleGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"updateHelpCenterArticleGroup": map[string]interface{}{
					"helpCenterArticleGroup": map[string]interface{}{
						"id":        "hcag_001",
						"name":      "New Name",
						"slug":      "new-name",
						"createdAt": map[string]string{"iso8601": "2025-05-01T00:00:00Z"},
						"updatedAt": map[string]string{"iso8601": "2025-06-01T00:00:00Z"},
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

	group, err := NewClientWithURL("test-token", server.URL).UpdateHelpCenterArticleGroup("hcag_001", "New Name")
	if err != nil {
		t.Fatalf("UpdateHelpCenterArticleGroup failed: %v", err)
	}
	if group.Name != "New Name" {
		t.Fatalf("expected updated name, got %+v", group)
	}
}

func TestDeleteHelpCenterArticleGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"deleteHelpCenterArticleGroup": map[string]interface{}{
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

	if err := NewClientWithURL("test-token", server.URL).DeleteHelpCenterArticleGroup("hcag_001"); err != nil {
		t.Fatalf("DeleteHelpCenterArticleGroup failed: %v", err)
	}
}
