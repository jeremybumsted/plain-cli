package plain

import "fmt"

// HelpCenterArticleUpsertInput contains the fields supported by Plain's
// upsertHelpCenterArticle mutation.
type HelpCenterArticleUpsertInput struct {
	HelpCenterID             string
	HelpCenterArticleID      string
	HelpCenterArticleGroupID string
	Title                    string
	Description              string
	ContentHTML              string
	Status                   string
}

// GetHelpCenterArticleBySlug fetches an article by its slug within a help center
func (c *Client) GetHelpCenterArticleBySlug(helpCenterID, slug string) (*HelpCenterArticle, error) {
	query := `
		query HelpCenterArticleBySlug($helpCenterId: ID!, $slug: String!) {
			helpCenterArticleBySlug(helpCenterId: $helpCenterId, slug: $slug) {
				id
				title
				description
				contentHtml
				slug
				status
				createdAt { iso8601 }
				updatedAt { iso8601 }
				articleGroup { id name slug }
			}
		}
	`

	variables := map[string]interface{}{
		"helpCenterId": helpCenterID,
		"slug":         slug,
	}

	var result struct {
		Data struct {
			Article *HelpCenterArticle `json:"helpCenterArticleBySlug"`
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

	if result.Data.Article == nil {
		return nil, &Error{
			StatusCode: 404,
			Message:    fmt.Sprintf("Article not found with slug: %s", slug),
		}
	}

	return result.Data.Article, nil
}

// UpsertHelpCenterArticle creates or updates a help center article.
func (c *Client) UpsertHelpCenterArticle(input HelpCenterArticleUpsertInput) (*HelpCenterArticle, error) {
	mutation := `
		mutation UpsertHelpCenterArticle($input: UpsertHelpCenterArticleInput!) {
			upsertHelpCenterArticle(input: $input) {
				helpCenterArticle {
					id
					title
					description
					contentHtml
					slug
					status
					createdAt { iso8601 }
					updatedAt { iso8601 }
					articleGroup { id name slug }
				}
				error {
					message
					code
					fields {
						field
						message
						type
					}
				}
			}
		}
	`

	inputMap := map[string]interface{}{
		"helpCenterId": input.HelpCenterID,
		"title":        input.Title,
		"contentHtml":  input.ContentHTML,
	}
	if input.HelpCenterArticleID != "" {
		inputMap["helpCenterArticleId"] = input.HelpCenterArticleID
	}
	if input.HelpCenterArticleGroupID != "" {
		inputMap["helpCenterArticleGroupId"] = input.HelpCenterArticleGroupID
	}
	if input.Description != "" {
		inputMap["description"] = input.Description
	}
	if input.Status != "" {
		inputMap["status"] = input.Status
	}

	var result struct {
		Data struct {
			UpsertHelpCenterArticle struct {
				HelpCenterArticle *HelpCenterArticle `json:"helpCenterArticle"`
				Error             *struct {
					Message string `json:"message"`
					Code    string `json:"code"`
				} `json:"error"`
			} `json:"upsertHelpCenterArticle"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := c.graphql(mutation, map[string]interface{}{"input": inputMap}, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	if result.Data.UpsertHelpCenterArticle.Error != nil {
		return nil, fmt.Errorf("%s", result.Data.UpsertHelpCenterArticle.Error.Message)
	}

	if result.Data.UpsertHelpCenterArticle.HelpCenterArticle == nil {
		return nil, fmt.Errorf("article upsert returned no article")
	}

	return result.Data.UpsertHelpCenterArticle.HelpCenterArticle, nil
}
