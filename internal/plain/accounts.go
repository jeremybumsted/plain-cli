package plain

import "fmt"

// Account represents a Plain account.
type Account struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ExternalID string   `json:"externalId,omitempty"`
	URL        string   `json:"url,omitempty"`
	Source     string   `json:"source,omitempty"`
	CreatedAt  DateTime `json:"createdAt"`
	UpdatedAt  DateTime `json:"updatedAt"`
	Tier       *Tier    `json:"tier,omitempty"`
}

// ListAccounts fetches accounts up to the requested limit.
func (c *Client) ListAccounts(limit int) ([]*Account, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		query ListAccounts($first: Int, $after: String) {
			accounts(first: $first, after: $after) {
				edges {
					node {
						id
						name
						externalId
						url
						source
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

	accounts := make([]*Account, 0, limit)
	cursor := ""
	pageSize := limit
	for len(accounts) < limit {
		variables := map[string]interface{}{"first": pageSize}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				Accounts struct {
					Edges []struct {
						Node Account `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
				} `json:"accounts"`
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

		for _, edge := range result.Data.Accounts.Edges {
			account := edge.Node
			accounts = append(accounts, &account)
			if len(accounts) == limit {
				break
			}
		}

		if !result.Data.Accounts.PageInfo.HasNextPage || result.Data.Accounts.PageInfo.EndCursor == "" {
			break
		}
		cursor = result.Data.Accounts.PageInfo.EndCursor
	}

	return accounts, nil
}

// GetAccount fetches a single account by ID.
func (c *Client) GetAccount(accountID string) (*Account, error) {
	query := `
		query GetAccount($accountId: ID!) {
			account(accountId: $accountId) {
				id
				name
				externalId
				url
				source
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
			Account *Account `json:"account"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, map[string]interface{}{"accountId": accountID}, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.Account == nil {
		return nil, &Error{StatusCode: 404, Message: fmt.Sprintf("Account not found: %s", accountID)}
	}

	return result.Data.Account, nil
}
