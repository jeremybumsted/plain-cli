package plain

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePriority converts a human-friendly or numeric priority into Plain's
// canonical 0-3 representation.
func ParsePriority(priority string) (int, error) {
	priority = strings.ToLower(strings.TrimSpace(priority))

	if value, err := strconv.Atoi(priority); err == nil {
		if value >= 0 && value <= 3 {
			return value, nil
		}
		return 0, fmt.Errorf("priority must be between 0 and 3, got %d", value)
	}

	switch priority {
	case "low":
		return 0, nil
	case "normal":
		return 1, nil
	case "high":
		return 2, nil
	case "urgent":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid priority %q, expected low, normal, high, urgent, or 0-3", priority)
	}
}
