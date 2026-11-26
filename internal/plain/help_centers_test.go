package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHelpCenter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"helpCenter": map[string]interface{}{"id": "hc_123", "type": "PUBLIC", "publicName": "Docs", "internalName": "Docs", "description": "Help center", "isDeleted": false, "createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"}, "updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"}},
			},
		})
	}))
	defer server.Close()

	hc, err := NewClientWithURL("test-token", server.URL).GetHelpCenter("hc_123")
	if err != nil {
		t.Fatalf("GetHelpCenter failed: %v", err)
	}
	if hc.ID != "hc_123" || hc.PublicName != "Docs" {
		t.Fatalf("unexpected help center response: %+v", hc)
	}
}
