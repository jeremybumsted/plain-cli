package articles

import (
	"fmt"
	"sort"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// ListCmd represents the "articles list" command
type ListCmd struct {
	ConfigPath   string `help:"Path to config file" default:""`
	Format       string `help:"Output format (table, json, quiet)" default:"table"`
	HelpCenterID string `help:"Help center ID (overrides config)" optional:""`
	Preview      bool   `help:"Include content preview" short:"p"`
}

// Run executes the list command
func (cmd *ListCmd) Run() error {
	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get help center ID
	helpCenterID := cmd.HelpCenterID
	if helpCenterID == "" {
		helpCenterID, err = cfg.GetHelpCenterID()
		if err != nil {
			return err
		}
	}

	// Get client
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	// Fetch articles (include content if preview is requested)
	articles, err := client.ListHelpCenterArticles(helpCenterID, cmd.Preview)
	if err != nil {
		return fmt.Errorf("failed to fetch articles: %w", err)
	}

	// Sort articles alphabetically by title
	sort.Slice(articles, func(i, j int) bool {
		return strings.ToLower(articles[i].Title) < strings.ToLower(articles[j].Title)
	})

	formatter := getFormatter(cmd.Format)

	// Handle different output formats
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(articles)

	case "quiet":
		// Just output article IDs
		for _, article := range articles {
			fmt.Println(article.ID)
		}
		return nil

	default:
		// Table format
		if cmd.Preview {
			return outputTableWithPreview(articles)
		}
		return outputTable(articles)
	}
}

func outputTable(articles []*plain.HelpCenterArticle) error {
	// Print header
	fmt.Printf("%-40s %-15s %-20s %-12s %s\n", "TITLE", "ID", "GROUP", "STATUS", "UPDATED")
	fmt.Println(strings.Repeat("-", 120))

	// Print articles
	for _, article := range articles {
		groupName := "-"
		if article.Group != nil {
			groupName = truncate(article.Group.Name, 20)
		}

		title := truncate(article.Title, 40)
		id := truncate(article.ID, 15)
		status := truncate(article.Status, 12)
		updated := article.UpdatedAt.ISO8601
		if len(updated) > 10 {
			updated = updated[:10]
		}

		fmt.Printf("%-40s %-15s %-20s %-12s %s\n", title, id, groupName, status, updated)
	}

	fmt.Printf("\nTotal: %d articles\n", len(articles))
	return nil
}

func outputTableWithPreview(articles []*plain.HelpCenterArticle) error {
	converter := md.NewConverter("", true, &md.Options{
		HeadingStyle: "atx",
	})

	for i, article := range articles {
		if i > 0 {
			fmt.Println()
		}

		// Article header
		groupName := "-"
		if article.Group != nil {
			groupName = article.Group.Name
		}

		fmt.Printf("Title:   %s\n", article.Title)
		fmt.Printf("ID:      %s\n", article.ID)
		fmt.Printf("Group:   %s\n", groupName)
		fmt.Printf("Status:  %s\n", article.Status)
		fmt.Printf("Updated: %s\n", article.UpdatedAt.ISO8601[:10])

		// Convert HTML to markdown and show preview
		if article.ContentHTML != "" {
			md, err := converter.ConvertString(article.ContentHTML)
			if err == nil {
				// Clean up and truncate to ~150 chars
				md = strings.TrimSpace(md)
				// Remove newlines for preview
				md = strings.ReplaceAll(md, "\n", " ")
				// Remove multiple spaces
				md = strings.Join(strings.Fields(md), " ")

				if len(md) > 150 {
					md = md[:150] + "..."
				}
				fmt.Printf("Preview: %s\n", md)
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Printf("\nTotal: %d articles\n", len(articles))
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
