package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTiers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"tiers": map[string]interface{}{
					"edges": []map[string]interface{}{
						{"node": map[string]interface{}{"id": "tier_123", "name": "Enterprise", "defaultThreadPriority": 3, "isDefault": false, "isMachineTier": false, "createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"}, "updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"}}},
					},
					"pageInfo": map[string]interface{}{"hasNextPage": false, "endCursor": ""},
				},
			},
		})
	}))
	defer server.Close()

	tiers, err := NewClientWithURL("test-token", server.URL).ListTiers()
	if err != nil {
		t.Fatalf("ListTiers failed: %v", err)
	}
	if len(tiers) != 1 || tiers[0].ID != "tier_123" {
		t.Fatalf("unexpected tiers response: %+v", tiers)
	}
}

func TestGetTierNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"tier": nil}})
	}))
	defer server.Close()

	_, err := NewClientWithURL("test-token", server.URL).GetTier("tier_missing")
	apiErr, ok := err.(*Error)
	if !ok || apiErr.StatusCode != 404 {
		t.Fatalf("expected 404 api error, got %T %v", err, err)
	}
}
