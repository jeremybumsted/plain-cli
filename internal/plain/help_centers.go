package plain

import "fmt"

// HelpCenterDetail represents detailed help center metadata.
type HelpCenterDetail struct {
	ID           string   `json:"id"`
	Type         string   `json:"type,omitempty"`
	PublicName   string   `json:"publicName,omitempty"`
	InternalName string   `json:"internalName,omitempty"`
	Description  string   `json:"description,omitempty"`
	IsDeleted    bool     `json:"isDeleted"`
	CreatedAt    DateTime `json:"createdAt"`
	UpdatedAt    DateTime `json:"updatedAt"`
}

// GetHelpCenter fetches a single help center by ID.
func (c *Client) GetHelpCenter(helpCenterID string) (*HelpCenterDetail, error) {
	query := `
		query GetHelpCenter($id: ID!) {
			helpCenter(id: $id) {
				id
				type
				publicName
				internalName
				description
				isDeleted
				createdAt { iso8601 }
				updatedAt { iso8601 }
			}
		}
	`

	var result struct {
		Data struct {
			HelpCenter *HelpCenterDetail `json:"helpCenter"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, map[string]interface{}{"id": helpCenterID}, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.HelpCenter == nil {
		return nil, &Error{StatusCode: 404, Message: fmt.Sprintf("Help center not found: %s", helpCenterID)}
	}

	return result.Data.HelpCenter, nil
}
