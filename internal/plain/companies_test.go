package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCompanies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"companies": map[string]interface{}{
					"edges": []map[string]interface{}{
						{"node": map[string]interface{}{"id": "co_123", "name": "Acme", "domainName": "acme.com", "isDeleted": false, "createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"}, "updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"}}},
					},
					"pageInfo": map[string]interface{}{"hasNextPage": false, "endCursor": ""},
				},
			},
		})
	}))
	defer server.Close()

	companies, err := NewClientWithURL("test-token", server.URL).ListCompanies(50)
	if err != nil {
		t.Fatalf("ListCompanies failed: %v", err)
	}
	if len(companies) != 1 || companies[0].ID != "co_123" {
		t.Fatalf("unexpected companies response: %+v", companies)
	}
}

func TestGetCompanyNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"company": nil}})
	}))
	defer server.Close()

	_, err := NewClientWithURL("test-token", server.URL).GetCompany("co_missing")
	apiErr, ok := err.(*Error)
	if !ok || apiErr.StatusCode != 404 {
		t.Fatalf("expected 404 api error, got %T %v", err, err)
	}
}
