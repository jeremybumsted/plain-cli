package articles

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// UpsertCmd represents the "articles upsert" command.
type UpsertCmd struct {
	ConfigPath   string `help:"Path to config file" default:""`
	Format       string `help:"Output format (table, json, quiet)" default:"table"`
	HelpCenterID string `help:"Help center ID (overrides config)" optional:""`
	ID           string `help:"Article ID to update" optional:""`
	GroupID      string `name:"group" help:"Article group ID" optional:""`
	Title        string `help:"Article title" required:""`
	Description  string `help:"Article description" required:""`
	Content      string `help:"Inline HTML content" optional:""`
	ContentFile  string `name:"content-file" help:"Path to a file containing HTML content" optional:"" type:"path"`
	Status       string `help:"Article status (draft or published)" default:"draft"`
}

// Run executes the articles upsert command.
func (cmd *UpsertCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	helpCenterID := cmd.HelpCenterID
	if helpCenterID == "" {
		helpCenterID, err = cfg.GetHelpCenterID()
		if err != nil {
			return err
		}
	}

	contentHTML, err := cmd.resolveContentHTML()
	if err != nil {
		return err
	}

	status, err := normalizeArticleStatus(cmd.Status)
	if err != nil {
		return err
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	article, err := client.UpsertHelpCenterArticle(plain.HelpCenterArticleUpsertInput{
		HelpCenterID:             helpCenterID,
		HelpCenterArticleID:      cmd.ID,
		HelpCenterArticleGroupID: cmd.GroupID,
		Title:                    cmd.Title,
		Description:              cmd.Description,
		ContentHTML:              contentHTML,
		Status:                   status,
	})
	if err != nil {
		return fmt.Errorf("failed to save article: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	if cmd.Format == "quiet" {
		return formatter.Print(article.ID)
	}
	if cmd.Format == "json" {
		return formatter.PrintJSON(article)
	}

	if err := formatter.Success("Article saved"); err != nil {
		return err
	}
	if err := formatter.Printf("ID:          %s\n", article.ID); err != nil {
		return err
	}
	if err := formatter.Printf("Title:       %s\n", article.Title); err != nil {
		return err
	}
	if article.Description != "" {
		if err := formatter.Printf("Description: %s\n", article.Description); err != nil {
			return err
		}
	}
	if err := formatter.Printf("Slug:        %s\n", article.Slug); err != nil {
		return err
	}
	if err := formatter.Printf("Status:      %s\n", article.Status); err != nil {
		return err
	}
	if article.Group != nil {
		if err := formatter.Printf("Group:       %s (%s)\n", article.Group.Name, article.Group.ID); err != nil {
			return err
		}
	}
	if link := buildArticleLink(cfg, helpCenterID, article.ID); link != "" {
		if err := formatter.Printf("Open:        %s\n", link); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *UpsertCmd) resolveContentHTML() (string, error) {
	content := cmd.Content
	if cmd.ContentFile != "" {
		data, err := os.ReadFile(cmd.ContentFile)
		if err != nil {
			return "", fmt.Errorf("failed to read content file: %w", err)
		}
		content = string(data)
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return "", fmt.Errorf("article content required: use --content or --content-file")
	}

	return content, nil
}

func normalizeArticleStatus(status string) (string, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "DRAFT", "PUBLISHED":
		return status, nil
	default:
		return "", fmt.Errorf("invalid article status %q, expected draft or published", status)
	}
}

func buildArticleLink(cfg interface{ GetWorkspaceID() (string, error) }, helpCenterID, articleID string) string {
	workspaceID, err := cfg.GetWorkspaceID()
	if err != nil || workspaceID == "" {
		return ""
	}
	return fmt.Sprintf("https://app.plain.com/workspace/%s/help-center/%s/articles/%s/", workspaceID, helpCenterID, articleID)
}
