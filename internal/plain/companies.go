package plain

import "fmt"

// Company represents a Plain company.
type Company struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	DomainName string   `json:"domainName"`
	CreatedAt  DateTime `json:"createdAt"`
	UpdatedAt  DateTime `json:"updatedAt"`
	IsDeleted  bool     `json:"isDeleted"`
	Tier       *Tier    `json:"tier,omitempty"`
}

// ListCompanies fetches companies up to the requested limit.
func (c *Client) ListCompanies(limit int) ([]*Company, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		query ListCompanies($first: Int, $after: String) {
			companies(first: $first, after: $after) {
				edges {
					node {
						id
						name
						domainName
						isDeleted
						createdAt { iso8601 }
						updatedAt { iso8601 }
						tier {
							id
							name
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	companies := make([]*Company, 0, limit)
	cursor := ""
	pageSize := limit
	for len(companies) < limit {
		variables := map[string]interface{}{"first": pageSize}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				Companies struct {
					Edges []struct {
						Node Company `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
				} `json:"companies"`
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

		for _, edge := range result.Data.Companies.Edges {
			company := edge.Node
			companies = append(companies, &company)
			if len(companies) == limit {
				break
			}
		}

		if !result.Data.Companies.PageInfo.HasNextPage || result.Data.Companies.PageInfo.EndCursor == "" {
			break
		}
		cursor = result.Data.Companies.PageInfo.EndCursor
	}

	return companies, nil
}

// GetCompany fetches a single company by ID.
func (c *Client) GetCompany(companyID string) (*Company, error) {
	query := `
		query GetCompany($companyId: ID!) {
			company(companyId: $companyId) {
				id
				name
				domainName
				isDeleted
				createdAt { iso8601 }
				updatedAt { iso8601 }
				tier {
					id
					name
				}
			}
		}
	`

	var result struct {
		Data struct {
			Company *Company `json:"company"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, map[string]interface{}{"companyId": companyID}, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.Company == nil {
		return nil, &Error{StatusCode: 404, Message: fmt.Sprintf("Company not found: %s", companyID)}
	}

	return result.Data.Company, nil
}
