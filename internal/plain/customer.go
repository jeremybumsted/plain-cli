package plain

import (
	"fmt"
	"strings"
)

// Customer represents a Plain customer
type Customer struct {
	ID         string        `json:"id"`
	ExternalID string        `json:"externalId,omitempty"`
	FullName   string        `json:"fullName"`
	ShortName  string        `json:"shortName,omitempty"`
	Email      CustomerEmail `json:"email"`
	AvatarURL  string        `json:"avatarUrl,omitempty"`
	CreatedAt  DateTime      `json:"createdAt"`
	UpdatedAt  DateTime      `json:"updatedAt"`
}

// CustomerEmail represents a customer's email information
type CustomerEmail struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}

// CustomersResponse represents a paginated list of customers
type CustomersResponse struct {
	Customers []Customer `json:"customers"`
	Total     int        `json:"total"`
}

// customerFields is the shared GraphQL fragment for customer queries
const customerFields = `
	id
	externalId
	fullName
	shortName
	email {
		email
		isVerified
	}
	avatarUrl
	createdAt { iso8601 }
	updatedAt { iso8601 }
`

// ListCustomers fetches a paginated list of customers
func (c *Client) ListCustomers(first int, after string) (*CustomersResponse, error) {
	if first <= 0 {
		first = 50
	}

	query := `
		query ListCustomers($first: Int, $after: String) {
			customers(first: $first, after: $after) {
				edges {
					node {
						` + customerFields + `
					}
				}
				totalCount
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var result struct {
		Data struct {
			Customers struct {
				Edges []struct {
					Node Customer `json:"node"`
				} `json:"edges"`
				TotalCount int `json:"totalCount"`
			} `json:"customers"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	customers := make([]Customer, 0, len(result.Data.Customers.Edges))
	for _, edge := range result.Data.Customers.Edges {
		customers = append(customers, edge.Node)
	}

	return &CustomersResponse{
		Customers: customers,
		Total:     result.Data.Customers.TotalCount,
	}, nil
}

// GetCustomer fetches a single customer by ID
func (c *Client) GetCustomer(customerID string) (*Customer, error) {
	query := `
		query GetCustomer($customerId: ID!) {
			customer(customerId: $customerId) {
				` + customerFields + `
			}
		}
	`

	variables := map[string]interface{}{
		"customerId": customerID,
	}

	var result struct {
		Data struct {
			Customer *Customer `json:"customer"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.Customer == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Customer not found: %s", customerID),
		}
	}

	return result.Data.Customer, nil
}

// GetCustomerByEmail fetches a customer by email address
func (c *Client) GetCustomerByEmail(email string) (*Customer, error) {
	query := `
		query GetCustomerByEmail($email: String!) {
			customerByEmail(email: $email) {
				` + customerFields + `
			}
		}
	`

	variables := map[string]interface{}{
		"email": email,
	}

	var result struct {
		Data struct {
			CustomerByEmail *Customer `json:"customerByEmail"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.CustomerByEmail == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Customer not found with email: %s", email),
		}
	}

	return result.Data.CustomerByEmail, nil
}

// GetCustomerByExternalID fetches a customer by external ID
func (c *Client) GetCustomerByExternalID(externalID string) (*Customer, error) {
	query := `
		query GetCustomerByExternalId($externalId: ID!) {
			customerByExternalId(externalId: $externalId) {
				` + customerFields + `
			}
		}
	`

	variables := map[string]interface{}{
		"externalId": externalID,
	}

	var result struct {
		Data struct {
			CustomerByExternalId *Customer `json:"customerByExternalId"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, variables, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.CustomerByExternalId == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Customer not found with external ID: %s", externalID),
		}
	}

	return result.Data.CustomerByExternalId, nil
}

// SearchCustomers searches for customers matching a query
func (c *Client) SearchCustomers(searchQuery string, first int) (*CustomersResponse, error) {
	if first <= 0 {
		first = 50
	}

	query := `
		query SearchCustomers($first: Int, $after: String) {
			customers(first: $first, after: $after) {
				edges {
					node {
						` + customerFields + `
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
				totalCount
			}
		}
	`

	query = strings.TrimSpace(query)
	searchQuery = strings.ToLower(strings.TrimSpace(searchQuery))

	customers := make([]Customer, 0, first)
	cursor := ""
	pageSize := first
	if pageSize < 100 {
		pageSize = 100
	}

	for len(customers) < first {
		variables := map[string]interface{}{
			"first": pageSize,
		}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				Customers struct {
					Edges []struct {
						Node Customer `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
					TotalCount int `json:"totalCount"`
				} `json:"customers"`
			} `json:"data"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := c.graphql(query, variables, &result); err != nil {
			return nil, err
		}

		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
		}

		for _, edge := range result.Data.Customers.Edges {
			if customerMatchesQuery(edge.Node, searchQuery) {
				customers = append(customers, edge.Node)
				if len(customers) == first {
					break
				}
			}
		}

		if !result.Data.Customers.PageInfo.HasNextPage || result.Data.Customers.PageInfo.EndCursor == "" {
			break
		}
		cursor = result.Data.Customers.PageInfo.EndCursor
	}

	return &CustomersResponse{
		Customers: customers,
		Total:     len(customers),
	}, nil
}

func customerMatchesQuery(customer Customer, query string) bool {
	if query == "" {
		return true
	}
	for _, candidate := range []string{
		customer.ID,
		customer.ExternalID,
		customer.FullName,
		customer.ShortName,
		customer.Email.Email,
	} {
		if strings.Contains(strings.ToLower(candidate), query) {
			return true
		}
	}
	return false
}
