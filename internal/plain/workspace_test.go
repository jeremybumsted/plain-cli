package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"myWorkspace": map[string]interface{}{"id": "ws_123", "name": "Workspace", "publicName": "Public Workspace", "isDemoWorkspace": false, "createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"}, "updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"}},
			},
		})
	}))
	defer server.Close()

	workspace, err := NewClientWithURL("test-token", server.URL).GetWorkspace()
	if err != nil {
		t.Fatalf("GetWorkspace failed: %v", err)
	}
	if workspace.ID != "ws_123" || workspace.PublicName != "Public Workspace" {
		t.Fatalf("unexpected workspace response: %+v", workspace)
	}
}
