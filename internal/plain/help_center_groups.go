package plain

import "fmt"

// HelpCenterArticleGroup represents a help center article group and its direct
// children.
type HelpCenterArticleGroup struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Slug               string              `json:"slug"`
	CreatedAt          DateTime            `json:"createdAt"`
	UpdatedAt          DateTime            `json:"updatedAt"`
	ParentArticleGroup *ArticleGroup       `json:"parentArticleGroup,omitempty"`
	ChildArticleGroups []ArticleGroup      `json:"childArticleGroups,omitempty"`
	Articles           []HelpCenterArticle `json:"articles,omitempty"`
}

// GetHelpCenterArticleGroup fetches a help center article group by ID.
func (c *Client) GetHelpCenterArticleGroup(groupID string) (*HelpCenterArticleGroup, error) {
	query := `
		query GetHelpCenterArticleGroup($id: ID!) {
			helpCenterArticleGroup(id: $id) {
				id
				name
				slug
				createdAt { iso8601 }
				updatedAt { iso8601 }
				parentArticleGroup {
					id
					name
					slug
				}
				childArticleGroups(first: 50) {
					edges {
						node {
							id
							name
							slug
						}
					}
				}
				articles(first: 50) {
					edges {
						node {
							id
							title
							description
							slug
							status
							createdAt { iso8601 }
							updatedAt { iso8601 }
							articleGroup {
								id
								name
								slug
							}
						}
					}
				}
			}
		}
	`

	var result struct {
		Data struct {
			HelpCenterArticleGroup *struct {
				ID                 string        `json:"id"`
				Name               string        `json:"name"`
				Slug               string        `json:"slug"`
				CreatedAt          DateTime      `json:"createdAt"`
				UpdatedAt          DateTime      `json:"updatedAt"`
				ParentArticleGroup *ArticleGroup `json:"parentArticleGroup"`
				ChildArticleGroups struct {
					Edges []struct {
						Node ArticleGroup `json:"node"`
					} `json:"edges"`
				} `json:"childArticleGroups"`
				Articles struct {
					Edges []struct {
						Node HelpCenterArticle `json:"node"`
					} `json:"edges"`
				} `json:"articles"`
			} `json:"helpCenterArticleGroup"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(query, map[string]interface{}{"id": groupID}, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.HelpCenterArticleGroup == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Help center article group not found: %s", groupID),
		}
	}

	group := &HelpCenterArticleGroup{
		ID:                 result.Data.HelpCenterArticleGroup.ID,
		Name:               result.Data.HelpCenterArticleGroup.Name,
		Slug:               result.Data.HelpCenterArticleGroup.Slug,
		CreatedAt:          result.Data.HelpCenterArticleGroup.CreatedAt,
		UpdatedAt:          result.Data.HelpCenterArticleGroup.UpdatedAt,
		ParentArticleGroup: result.Data.HelpCenterArticleGroup.ParentArticleGroup,
		ChildArticleGroups: make([]ArticleGroup, 0, len(result.Data.HelpCenterArticleGroup.ChildArticleGroups.Edges)),
		Articles:           make([]HelpCenterArticle, 0, len(result.Data.HelpCenterArticleGroup.Articles.Edges)),
	}

	for _, edge := range result.Data.HelpCenterArticleGroup.ChildArticleGroups.Edges {
		group.ChildArticleGroups = append(group.ChildArticleGroups, edge.Node)
	}
	for _, edge := range result.Data.HelpCenterArticleGroup.Articles.Edges {
		group.Articles = append(group.Articles, edge.Node)
	}

	return group, nil
}

// CreateHelpCenterArticleGroup creates a help center article group.
func (c *Client) CreateHelpCenterArticleGroup(helpCenterID, name, parentGroupID string) (*HelpCenterArticleGroup, error) {
	mutation := `
		mutation CreateHelpCenterArticleGroup($input: CreateHelpCenterArticleGroupInput!) {
			createHelpCenterArticleGroup(input: $input) {
				helpCenterArticleGroup {
					id
					name
					slug
					createdAt { iso8601 }
					updatedAt { iso8601 }
					parentArticleGroup {
						id
						name
						slug
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	input := map[string]interface{}{
		"helpCenterId": helpCenterID,
		"name":         name,
	}
	if parentGroupID != "" {
		input["parentHelpCenterArticleGroupId"] = parentGroupID
	}

	var result struct {
		Data struct {
			CreateHelpCenterArticleGroup struct {
				HelpCenterArticleGroup *HelpCenterArticleGroup `json:"helpCenterArticleGroup"`
				Error                  *struct {
					Message string `json:"message"`
				} `json:"error"`
			} `json:"createHelpCenterArticleGroup"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, map[string]interface{}{"input": input}, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.CreateHelpCenterArticleGroup.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.CreateHelpCenterArticleGroup.Error.Message)
	}
	if result.Data.CreateHelpCenterArticleGroup.HelpCenterArticleGroup == nil {
		return nil, fmt.Errorf("create article group returned no group")
	}

	return result.Data.CreateHelpCenterArticleGroup.HelpCenterArticleGroup, nil
}

// UpdateHelpCenterArticleGroup updates a help center article group's name.
func (c *Client) UpdateHelpCenterArticleGroup(groupID, name string) (*HelpCenterArticleGroup, error) {
	mutation := `
		mutation UpdateHelpCenterArticleGroup($input: UpdateHelpCenterArticleGroupInput!) {
			updateHelpCenterArticleGroup(input: $input) {
				helpCenterArticleGroup {
					id
					name
					slug
					createdAt { iso8601 }
					updatedAt { iso8601 }
					parentArticleGroup {
						id
						name
						slug
					}
				}
				error {
					message
					code
				}
			}
		}
	`

	var result struct {
		Data struct {
			UpdateHelpCenterArticleGroup struct {
				HelpCenterArticleGroup *HelpCenterArticleGroup `json:"helpCenterArticleGroup"`
				Error                  *struct {
					Message string `json:"message"`
				} `json:"error"`
			} `json:"updateHelpCenterArticleGroup"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"helpCenterArticleGroupId": groupID,
			"name":                     name,
		},
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.UpdateHelpCenterArticleGroup.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.UpdateHelpCenterArticleGroup.Error.Message)
	}
	if result.Data.UpdateHelpCenterArticleGroup.HelpCenterArticleGroup == nil {
		return nil, fmt.Errorf("update article group returned no group")
	}

	return result.Data.UpdateHelpCenterArticleGroup.HelpCenterArticleGroup, nil
}

// DeleteHelpCenterArticleGroup deletes a help center article group by ID.
func (c *Client) DeleteHelpCenterArticleGroup(groupID string) error {
	mutation := `
		mutation DeleteHelpCenterArticleGroup($input: DeleteHelpCenterArticleGroupInput!) {
			deleteHelpCenterArticleGroup(input: $input) {
				error {
					message
					code
				}
			}
		}
	`

	var result struct {
		Data struct {
			DeleteHelpCenterArticleGroup struct {
				Error *struct {
					Message string `json:"message"`
				} `json:"error"`
			} `json:"deleteHelpCenterArticleGroup"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"helpCenterArticleGroupId": groupID,
		},
	}

	if err := c.graphql(mutation, variables, &result); err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}
	if result.Data.DeleteHelpCenterArticleGroup.Error != nil {
		return fmt.Errorf("%s", result.Data.DeleteHelpCenterArticleGroup.Error.Message)
	}

	return nil
}
