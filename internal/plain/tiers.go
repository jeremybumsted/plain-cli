package plain

import "fmt"

// Tier represents a Plain support tier.
type Tier struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	ExternalID            string   `json:"externalId,omitempty"`
	Color                 string   `json:"color,omitempty"`
	IsDefault             bool     `json:"isDefault"`
	DefaultThreadPriority int      `json:"defaultThreadPriority"`
	IsMachineTier         bool     `json:"isMachineTier"`
	CreatedAt             DateTime `json:"createdAt"`
	UpdatedAt             DateTime `json:"updatedAt"`
}

// ListTiers fetches all tiers in the workspace.
func (c *Client) ListTiers() ([]*Tier, error) {
	query := `
		query ListTiers($first: Int, $after: String) {
			tiers(first: $first, after: $after) {
				edges {
					node {
						id
						name
						externalId
						color
						isDefault
						defaultThreadPriority
						isMachineTier
						createdAt { iso8601 }
						updatedAt { iso8601 }
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	tiers := []*Tier{}
	cursor := ""
	for {
		variables := map[string]interface{}{"first": 100}
		if cursor != "" {
			variables["after"] = cursor
		}

		var result struct {
			Data struct {
				Tiers struct {
					Edges []struct {
						Node Tier `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
				} `json:"tiers"`
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

		for _, edge := range result.Data.Tiers.Edges {
			tier := edge.Node
			tiers = append(tiers, &tier)
		}

		if !result.Data.Tiers.PageInfo.HasNextPage || result.Data.Tiers.PageInfo.EndCursor == "" {
			break
		}
		cursor = result.Data.Tiers.PageInfo.EndCursor
	}

	return tiers, nil
}

// GetTier fetches a single tier by ID.
func (c *Client) GetTier(tierID string) (*Tier, error) {
	query := `
		query GetTier($tierId: ID!) {
			tier(tierId: $tierId) {
				id
				name
				externalId
				color
				isDefault
				defaultThreadPriority
				isMachineTier
				createdAt { iso8601 }
				updatedAt { iso8601 }
			}
		}
	`

	var result struct {
		Data struct {
			Tier *Tier `json:"tier"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, map[string]interface{}{"tierId": tierID}, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.Tier == nil {
		return nil, &Error{StatusCode: 404, Message: fmt.Sprintf("Tier not found: %s", tierID)}
	}

	return result.Data.Tier, nil
}
