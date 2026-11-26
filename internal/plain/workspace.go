package plain

import "fmt"

// WorkspaceDetail represents detailed workspace information.
type WorkspaceDetail struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	PublicName      string   `json:"publicName"`
	IsDemoWorkspace bool     `json:"isDemoWorkspace"`
	CreatedAt       DateTime `json:"createdAt"`
	UpdatedAt       DateTime `json:"updatedAt"`
}

// GetWorkspace fetches detailed information for the authenticated workspace.
func (c *Client) GetWorkspace() (*WorkspaceDetail, error) {
	query := `
		query GetWorkspace {
			myWorkspace {
				id
				name
				publicName
				isDemoWorkspace
				createdAt { iso8601 }
				updatedAt { iso8601 }
			}
		}
	`

	var result struct {
		Data struct {
			MyWorkspace *WorkspaceDetail `json:"myWorkspace"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, nil, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.MyWorkspace == nil {
		return nil, &Error{StatusCode: 404, Message: "Workspace not found"}
	}

	return result.Data.MyWorkspace, nil
}
