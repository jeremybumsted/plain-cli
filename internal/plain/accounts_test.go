package plain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestListAccounts(t *testing.T) {
	field := regexp.MustCompile(`\{\s*accounts\s*\(\s*first\s*:\s*\$first\b[\s,]+after\s*:\s*\$after\s*,?\s*\)`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if !field.MatchString(req.Query) || req.Variables["first"] != float64(50) {
			t.Errorf("expected accounts query with first=50, got query %q and variables %v", req.Query, req.Variables)
		}
		switch after, present := req.Variables["after"]; {
		case !present:
			_, _ = fmt.Fprint(w, `{"data":{"accounts":{
				"edges":[{"node":{"id":"ten_123","name":"Account A","externalId":"ext-a","url":"https://a.example.com","source":"API","createdAt":{"iso8601":"2026-01-01T00:00:00Z"},"updatedAt":{"iso8601":"2026-01-02T00:00:00Z"}}}],
				"pageInfo":{"hasNextPage":true,"endCursor":"next-page"}
			}}}`)
		case after == "next-page":
			_, _ = fmt.Fprint(w, `{"data":{"accounts":{
				"edges":[{"node":{"id":"ten_456","name":"Account B"}}],
				"pageInfo":{"hasNextPage":false,"endCursor":""}
			}}}`)
		default:
			t.Errorf("unexpected after variable: %v", after)
			http.Error(w, "invalid cursor", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	accounts, err := NewClientWithURL("test-token", server.URL).ListAccounts(50)
	if err != nil {
		t.Fatalf("ListAccounts failed: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}
	if first := accounts[0]; first.ID != "ten_123" || first.Name != "Account A" || first.ExternalID != "ext-a" || first.URL != "https://a.example.com" || first.Source != "API" || first.CreatedAt.ISO8601 != "2026-01-01T00:00:00Z" || first.UpdatedAt.ISO8601 != "2026-01-02T00:00:00Z" {
		t.Errorf("unexpected first account: %+v", first)
	}
	if accounts[1].ID != "ten_456" {
		t.Errorf("unexpected second account: %+v", accounts[1])
	}
}

func TestGetAccount(t *testing.T) {
	field := regexp.MustCompile(`\{\s*account\s*\(\s*accountId\s*:\s*\$accountId\s*,?\s*\)`)
	for _, tt := range []struct {
		name, id, response string
	}{
		{"found", "ten_123", `{"data":{"account":{"id":"ten_123","name":"Account A","tier":{"id":"tier_123","name":"Gold"}}}}`},
		{"not found", "ten_missing", `{"data":{"account":null}}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					Query     string         `json:"query"`
					Variables map[string]any `json:"variables"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("decode request: %v", err)
					http.Error(w, "invalid request", http.StatusBadRequest)
					return
				}
				if !field.MatchString(req.Query) || len(req.Variables) != 1 || req.Variables["accountId"] != tt.id {
					t.Errorf("expected account query with accountId=%q, got query %q and variables %v", tt.id, req.Query, req.Variables)
				}
				_, _ = fmt.Fprint(w, tt.response)
			}))
			defer server.Close()

			account, err := NewClientWithURL("test-token", server.URL).GetAccount(tt.id)
			if tt.name == "not found" {
				apiErr, ok := err.(*Error)
				if !ok || apiErr.StatusCode != 404 || apiErr.Message != "Account not found: ten_missing" || account != nil {
					t.Fatalf("expected nil account and 404 API error, got %+v and %T %v", account, err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetAccount failed: %v", err)
			}
			if account == nil || account.ID != tt.id || account.Name != "Account A" || account.Tier == nil || account.Tier.ID != "tier_123" || account.Tier.Name != "Gold" {
				t.Fatalf("unexpected account response: %+v", account)
			}
		})
	}
}
