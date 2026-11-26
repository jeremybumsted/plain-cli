package plain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCustomers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customers": map[string]interface{}{
					"edges": []map[string]interface{}{
						{
							"node": map[string]interface{}{
								"id":         "c_123",
								"externalId": "ext_1",
								"fullName":   "Test Customer",
								"shortName":  "TC",
								"email": map[string]interface{}{
									"email":      "test@example.com",
									"isVerified": true,
								},
								"avatarUrl": "",
								"createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"},
								"updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"},
							},
						},
						{
							"node": map[string]interface{}{
								"id":         "c_456",
								"externalId": "ext_2",
								"fullName":   "Another Customer",
								"shortName":  "AC",
								"email": map[string]interface{}{
									"email":      "another@example.com",
									"isVerified": false,
								},
								"avatarUrl": "",
								"createdAt": map[string]interface{}{"iso8601": "2026-01-03T00:00:00Z"},
								"updatedAt": map[string]interface{}{"iso8601": "2026-01-04T00:00:00Z"},
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   "",
					},
					"totalCount": 2,
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

	resp, err := client.ListCustomers(50, "")
	if err != nil {
		t.Fatalf("ListCustomers failed: %v", err)
	}

	if len(resp.Customers) != 2 {
		t.Fatalf("Expected 2 customers, got %d", len(resp.Customers))
	}

	if resp.Total != 2 {
		t.Errorf("Expected total 2, got %d", resp.Total)
	}

	if resp.Customers[0].ID != "c_123" {
		t.Errorf("Expected first customer ID 'c_123', got '%s'", resp.Customers[0].ID)
	}

	if resp.Customers[0].FullName != "Test Customer" {
		t.Errorf("Expected first customer fullName 'Test Customer', got '%s'", resp.Customers[0].FullName)
	}

	if resp.Customers[0].Email.Email != "test@example.com" {
		t.Errorf("Expected first customer email 'test@example.com', got '%s'", resp.Customers[0].Email.Email)
	}

	if resp.Customers[1].ID != "c_456" {
		t.Errorf("Expected second customer ID 'c_456', got '%s'", resp.Customers[1].ID)
	}

	if resp.Customers[1].FullName != "Another Customer" {
		t.Errorf("Expected second customer fullName 'Another Customer', got '%s'", resp.Customers[1].FullName)
	}
}

func TestGetCustomer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customer": map[string]interface{}{
					"id":         "c_123",
					"externalId": "ext_1",
					"fullName":   "Test Customer",
					"shortName":  "TC",
					"email": map[string]interface{}{
						"email":      "test@example.com",
						"isVerified": true,
					},
					"avatarUrl": "",
					"createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"},
					"updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"},
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

	customer, err := client.GetCustomer("c_123")
	if err != nil {
		t.Fatalf("GetCustomer failed: %v", err)
	}

	if customer.ID != "c_123" {
		t.Errorf("Expected customer ID 'c_123', got '%s'", customer.ID)
	}

	if customer.ExternalID != "ext_1" {
		t.Errorf("Expected customer externalId 'ext_1', got '%s'", customer.ExternalID)
	}

	if customer.FullName != "Test Customer" {
		t.Errorf("Expected customer fullName 'Test Customer', got '%s'", customer.FullName)
	}

	if customer.ShortName != "TC" {
		t.Errorf("Expected customer shortName 'TC', got '%s'", customer.ShortName)
	}

	if customer.Email.Email != "test@example.com" {
		t.Errorf("Expected customer email 'test@example.com', got '%s'", customer.Email.Email)
	}

	if !customer.Email.IsVerified {
		t.Errorf("Expected customer email to be verified")
	}
}

func TestGetCustomerByEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customerByEmail": map[string]interface{}{
					"id":         "c_123",
					"externalId": "ext_1",
					"fullName":   "Test Customer",
					"shortName":  "TC",
					"email": map[string]interface{}{
						"email":      "test@example.com",
						"isVerified": true,
					},
					"avatarUrl": "",
					"createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"},
					"updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"},
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

	customer, err := client.GetCustomerByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetCustomerByEmail failed: %v", err)
	}

	if customer.ID != "c_123" {
		t.Errorf("Expected customer ID 'c_123', got '%s'", customer.ID)
	}

	if customer.FullName != "Test Customer" {
		t.Errorf("Expected customer fullName 'Test Customer', got '%s'", customer.FullName)
	}

	if customer.Email.Email != "test@example.com" {
		t.Errorf("Expected customer email 'test@example.com', got '%s'", customer.Email.Email)
	}
}

func TestGetCustomerByEmailNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customerByEmail": nil,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	_, err := client.GetCustomerByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatal("Expected error for non-existent customer, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestGetCustomerByExternalID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customerByExternalId": map[string]interface{}{
					"id":         "c_123",
					"externalId": "ext_1",
					"fullName":   "Test Customer",
					"shortName":  "TC",
					"email": map[string]interface{}{
						"email":      "test@example.com",
						"isVerified": true,
					},
					"avatarUrl": "",
					"createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"},
					"updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"},
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

	customer, err := client.GetCustomerByExternalID("ext_1")
	if err != nil {
		t.Fatalf("GetCustomerByExternalID failed: %v", err)
	}

	if customer.ID != "c_123" {
		t.Errorf("Expected customer ID 'c_123', got '%s'", customer.ID)
	}

	if customer.ExternalID != "ext_1" {
		t.Errorf("Expected customer externalId 'ext_1', got '%s'", customer.ExternalID)
	}

	if customer.FullName != "Test Customer" {
		t.Errorf("Expected customer fullName 'Test Customer', got '%s'", customer.FullName)
	}

	if customer.Email.Email != "test@example.com" {
		t.Errorf("Expected customer email 'test@example.com', got '%s'", customer.Email.Email)
	}
}

func TestGetCustomerByExternalIDNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customerByExternalId": nil,
			},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClientWithURL("test-token", server.URL)

	_, err := client.GetCustomerByExternalID("nonexistent_ext")
	if err == nil {
		t.Fatal("Expected error for non-existent customer, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error type, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestSearchCustomers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"customers": map[string]interface{}{
					"edges": []map[string]interface{}{
						{
							"node": map[string]interface{}{
								"id":         "c_123",
								"externalId": "ext_1",
								"fullName":   "Test Customer",
								"shortName":  "TC",
								"email": map[string]interface{}{
									"email":      "test@example.com",
									"isVerified": true,
								},
								"avatarUrl": "",
								"createdAt": map[string]interface{}{"iso8601": "2026-01-01T00:00:00Z"},
								"updatedAt": map[string]interface{}{"iso8601": "2026-01-02T00:00:00Z"},
							},
						},
					},
					"totalCount": 1,
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

	resp, err := client.SearchCustomers("Test", 10)
	if err != nil {
		t.Fatalf("SearchCustomers failed: %v", err)
	}

	if len(resp.Customers) != 1 {
		t.Fatalf("Expected 1 customer, got %d", len(resp.Customers))
	}

	if resp.Total != 1 {
		t.Errorf("Expected total 1, got %d", resp.Total)
	}

	if resp.Customers[0].ID != "c_123" {
		t.Errorf("Expected customer ID 'c_123', got '%s'", resp.Customers[0].ID)
	}

	if resp.Customers[0].FullName != "Test Customer" {
		t.Errorf("Expected customer fullName 'Test Customer', got '%s'", resp.Customers[0].FullName)
	}

	if resp.Customers[0].Email.Email != "test@example.com" {
		t.Errorf("Expected customer email 'test@example.com', got '%s'", resp.Customers[0].Email.Email)
	}
}
