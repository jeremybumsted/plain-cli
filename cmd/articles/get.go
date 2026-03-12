package articles

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// GetCmd represents the "articles get" command
type GetCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (markdown, json, quiet)" default:"markdown"`
	ID         string `arg:"" name:"id" help:"Article ID (hca_*)"`
}

// Run executes the get command
func (cmd *GetCmd) Run() error {
	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get client
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	// Fetch article
	article, err := client.GetHelpCenterArticle(cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch article: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	// Handle different output formats
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(article)

	case "quiet":
		// Just output article ID
		fmt.Println(article.ID)
		return nil

	default:
		// Markdown format (default for articles)
		return outputMarkdown(article)
	}
}

func outputMarkdown(article *plain.HelpCenterArticle) error {
	// Convert HTML to markdown
	converter := md.NewConverter("", true, &md.Options{
		HeadingStyle: "atx",
	})

	markdown, err := converter.ConvertString(article.ContentHTML)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	// Output article with metadata header
	fmt.Printf("# %s\n\n", article.Title)
	fmt.Printf("Slug: %s\n", article.Slug)
	fmt.Printf("Status: %s\n", article.Status)
	fmt.Printf("Updated: %s\n", article.UpdatedAt.ISO8601)

	if article.Group != nil {
		fmt.Printf("Group: %s\n", article.Group.Name)
	}

	fmt.Println("\n---")
	fmt.Println(strings.TrimSpace(markdown))

	return nil
}
