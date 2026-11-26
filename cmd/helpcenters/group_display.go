package helpcenters

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

func displayGroup(formatter *output.Formatter, group *plain.HelpCenterArticleGroup) error {
	if err := formatter.Printf("ID:          %s\n", group.ID); err != nil {
		return err
	}
	if err := formatter.Printf("Name:        %s\n", group.Name); err != nil {
		return err
	}
	if err := formatter.Printf("Slug:        %s\n", group.Slug); err != nil {
		return err
	}
	if group.ParentArticleGroup != nil {
		if err := formatter.Printf("Parent:      %s (%s)\n", group.ParentArticleGroup.Name, group.ParentArticleGroup.ID); err != nil {
			return err
		}
	}
	createdAt, err := group.CreatedAt.Time()
	if err == nil && !createdAt.IsZero() {
		if err := formatter.Printf("Created:     %s\n", createdAt.Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
	}
	updatedAt, err := group.UpdatedAt.Time()
	if err == nil {
		if err := formatter.Printf("Updated:     %s\n", updatedAt.Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
	}
	if len(group.ChildArticleGroups) > 0 {
		if err := formatter.Print("Child Groups:"); err != nil {
			return err
		}
		for _, child := range group.ChildArticleGroups {
			if err := formatter.Printf("  - %s (%s)\n", child.Name, child.ID); err != nil {
				return err
			}
		}
	}
	if len(group.Articles) > 0 {
		if err := formatter.Print("Articles:"); err != nil {
			return err
		}
		for _, article := range group.Articles {
			if err := formatter.Printf("  - %s (%s)\n", article.Title, article.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func formatGroupNotFound(err error, groupID string, formatter *output.Formatter) error {
	if apiErr, ok := err.(*plain.Error); ok && apiErr.StatusCode == 404 {
		return formatter.Error(fmt.Sprintf("Article group not found: %s", groupID))
	}
	return err
}
