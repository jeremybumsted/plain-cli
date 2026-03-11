package threads

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// formatTime converts a time.Time to a human-readable relative time string
func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// truncateString truncates a string to a maximum length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// confirmAction prompts the user for confirmation before performing an action
// Returns true if the user confirms (y/yes), false otherwise
// If skipPrompt is true, immediately returns true without prompting
func confirmAction(message string, skipPrompt bool) (bool, error) {
	if skipPrompt {
		return true, nil
	}

	fmt.Printf("%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// parseRelativeTime converts relative time strings like "2h", "1d" to absolute time
// Also handles absolute ISO8601 timestamps
// Supported relative formats: 30m, 2h, 1d, 3d, 1w
func parseRelativeTime(input string) (time.Time, error) {
	input = strings.TrimSpace(input)

	// Try parsing as ISO8601 first
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	// Parse relative time format (e.g., "2h", "1d", "3w")
	re := regexp.MustCompile(`^(\d+)([mhdw])$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s (expected format: 30m, 2h, 1d, 3w or ISO8601)", input)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time value: %s", matches[1])
	}

	unit := matches[2]
	var duration time.Duration

	switch unit {
	case "m":
		duration = time.Duration(value) * time.Minute
	case "h":
		duration = time.Duration(value) * time.Hour
	case "d":
		duration = time.Duration(value) * 24 * time.Hour
	case "w":
		duration = time.Duration(value) * 7 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("invalid time unit: %s", unit)
	}

	return time.Now().Add(duration), nil
}

// openEditor opens the user's preferred editor ($EDITOR) for composing text
// If $EDITOR is not set, falls back to prompting for text input
// Returns the edited/entered text
func openEditor(initialText string) (string, error) {
	editor := os.Getenv("EDITOR")

	// If no EDITOR is set, fall back to simple text input
	if editor == "" {
		fmt.Println("$EDITOR not set. Enter text (Ctrl+D when done):")
		reader := bufio.NewReader(os.Stdin)
		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, ""), nil
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "plain-note-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	// Write initial text if provided
	if initialText != "" {
		if _, err := tmpFile.WriteString(initialText); err != nil {
			_ = tmpFile.Close()
			return "", fmt.Errorf("failed to write initial text: %w", err)
		}
	}
	_ = tmpFile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	// Read the edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read edited content: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// extractThreadID extracts the thread ID from a URL or returns the input as-is
// Handles URLs like: https://app.plain.com/threads/thread_xxx
// Returns: thread_xxx or the original input if no thread ID pattern is found
func extractThreadID(input string) string {
	// Pattern: thread_xxx where xxx is alphanumeric
	re := regexp.MustCompile(`thread_[a-zA-Z0-9]+`)
	if match := re.FindString(input); match != "" {
		return match
	}
	return input
}
