package util

import (
	"strings"
	"testing"
	"time"
)

func TestParseDateToISO8601_ISO8601Format(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, result string)
	}{
		{
			name:    "valid RFC3339 format",
			input:   "2026-03-15T10:30:00Z",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				if result != "2026-03-15T10:30:00Z" {
					t.Errorf("expected 2026-03-15T10:30:00Z, got %s", result)
				}
			},
		},
		{
			name:    "RFC3339 with timezone conversion to UTC",
			input:   "2026-03-15T10:30:00+05:00",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				if result != "2026-03-15T05:30:00Z" {
					t.Errorf("expected UTC conversion to 2026-03-15T05:30:00Z, got %s", result)
				}
			},
		},
		{
			name:    "date only format (YYYY-MM-DD)",
			input:   "2026-03-15",
			wantErr: false,
			validate: func(t *testing.T, result string) {
				if result != "2026-03-15T00:00:00Z" {
					t.Errorf("expected 2026-03-15T00:00:00Z, got %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToISO8601(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateToISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestParseDateToISO8601_RelativeDates(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected time.Time
	}{
		{
			name:     "7 days ago",
			input:    "7d",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -7)),
		},
		{
			name:     "1 day ago",
			input:    "1d",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -1)),
		},
		{
			name:     "30 days ago",
			input:    "30d",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -30)),
		},
		{
			name:     "1 week ago",
			input:    "1w",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -7)),
		},
		{
			name:     "2 weeks ago",
			input:    "2w",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -14)),
		},
		{
			name:     "1 month ago",
			input:    "1M",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, -1, 0)),
		},
		{
			name:     "3 months ago",
			input:    "3M",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, -3, 0)),
		},
		{
			name:     "1 year ago",
			input:    "1y",
			wantErr:  false,
			expected: startOfDay(now.AddDate(-1, 0, 0)),
		},
		{
			name:     "2 years ago",
			input:    "2y",
			wantErr:  false,
			expected: startOfDay(now.AddDate(-2, 0, 0)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToISO8601(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateToISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parsedResult, err := time.Parse(time.RFC3339, result)
				if err != nil {
					t.Errorf("failed to parse result: %v", err)
					return
				}
				if !parsedResult.Equal(tt.expected) {
					t.Errorf("ParseDateToISO8601() = %v, want %v", parsedResult, tt.expected)
				}
			}
		})
	}
}

func TestParseDateToISO8601_HumanReadable(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected time.Time
	}{
		{
			name:     "today",
			input:    "today",
			wantErr:  false,
			expected: startOfDay(now),
		},
		{
			name:     "yesterday",
			input:    "yesterday",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -1)),
		},
		{
			name:     "last-week",
			input:    "last-week",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -7)),
		},
		{
			name:     "last-month",
			input:    "last-month",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, -1, 0)),
		},
		{
			name:     "last-year",
			input:    "last-year",
			wantErr:  false,
			expected: startOfDay(now.AddDate(-1, 0, 0)),
		},
		{
			name:     "today with uppercase",
			input:    "TODAY",
			wantErr:  false,
			expected: startOfDay(now),
		},
		{
			name:     "yesterday with mixed case",
			input:    "YeSteRdaY",
			wantErr:  false,
			expected: startOfDay(now.AddDate(0, 0, -1)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToISO8601(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateToISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parsedResult, err := time.Parse(time.RFC3339, result)
				if err != nil {
					t.Errorf("failed to parse result: %v", err)
					return
				}
				if !parsedResult.Equal(tt.expected) {
					t.Errorf("ParseDateToISO8601() = %v, want %v", parsedResult, tt.expected)
				}
			}
		})
	}
}

func TestParseDateToISO8601_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
			errMsg:  "date string cannot be empty",
		},
		{
			name:    "invalid format",
			input:   "not-a-date",
			wantErr: true,
			errMsg:  "unsupported date format",
		},
		{
			name:    "invalid relative date - missing unit",
			input:   "7",
			wantErr: true,
			errMsg:  "unsupported date format",
		},
		{
			name:    "invalid relative date - missing number",
			input:   "d",
			wantErr: true,
			errMsg:  "unsupported date format",
		},
		{
			name:    "invalid relative date - wrong unit",
			input:   "7x",
			wantErr: true,
			errMsg:  "unsupported date format",
		},
		{
			name:    "garbage input",
			input:   "abc123xyz",
			wantErr: true,
			errMsg:  "unsupported date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateToISO8601(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateToISO8601() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			}
			if !tt.wantErr && result == "" {
				t.Errorf("expected non-empty result, got empty string")
			}
		})
	}
}

func TestValidateDateRange_Valid(t *testing.T) {
	tests := []struct {
		name   string
		after  string
		before string
	}{
		{
			name:   "both empty",
			after:  "",
			before: "",
		},
		{
			name:   "only after provided",
			after:  "7d",
			before: "",
		},
		{
			name:   "only before provided",
			after:  "",
			before: "today",
		},
		{
			name:   "valid range - ISO8601",
			after:  "2026-03-01T00:00:00Z",
			before: "2026-03-15T00:00:00Z",
		},
		{
			name:   "valid range - relative dates",
			after:  "7d",
			before: "1d",
		},
		{
			name:   "valid range - human readable",
			after:  "last-week",
			before: "yesterday",
		},
		{
			name:   "valid range - mixed formats",
			after:  "2026-03-01",
			before: "today",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateRange(tt.after, tt.before)
			if err != nil {
				t.Errorf("ValidateDateRange() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateDateRange_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		after   string
		before  string
		errMsg  string
	}{
		{
			name:   "after is after before",
			after:  "2026-03-15T00:00:00Z",
			before: "2026-03-01T00:00:00Z",
			errMsg: "must be before",
		},
		{
			name:   "after equals before",
			after:  "2026-03-15T00:00:00Z",
			before: "2026-03-15T00:00:00Z",
			errMsg: "must be before",
		},
		{
			name:   "invalid after date",
			after:  "invalid-date",
			before: "2026-03-15T00:00:00Z",
			errMsg: "invalid 'after' date",
		},
		{
			name:   "invalid before date",
			after:  "2026-03-01T00:00:00Z",
			before: "invalid-date",
			errMsg: "invalid 'before' date",
		},
		{
			name:   "both dates invalid",
			after:  "invalid1",
			before: "invalid2",
			errMsg: "invalid 'after' date",
		},
		{
			name:   "after is today, before is yesterday",
			after:  "today",
			before: "yesterday",
			errMsg: "must be before",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateRange(tt.after, tt.before)
			if err == nil {
				t.Errorf("ValidateDateRange() expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateDateRange() error = %v, should contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestStartOfDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "midday",
			input:    time.Date(2026, 3, 15, 12, 30, 45, 123456789, time.UTC),
			expected: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "midnight",
			input:    time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "end of day",
			input:    time.Date(2026, 3, 15, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := startOfDay(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("startOfDay() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseDateToISO8601_ResultFormat(t *testing.T) {
	// Test that all results are valid RFC3339 format
	inputs := []string{
		"today",
		"yesterday",
		"7d",
		"1w",
		"2M",
		"1y",
		"last-week",
		"2026-03-15",
		"2026-03-15T10:30:00Z",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			result, err := ParseDateToISO8601(input)
			if err != nil {
				t.Errorf("ParseDateToISO8601(%q) unexpected error: %v", input, err)
				return
			}

			// Verify result is valid RFC3339
			_, err = time.Parse(time.RFC3339, result)
			if err != nil {
				t.Errorf("ParseDateToISO8601(%q) produced invalid RFC3339: %v", input, err)
			}

			// Verify result ends with Z (UTC)
			if !strings.HasSuffix(result, "Z") {
				t.Errorf("ParseDateToISO8601(%q) result %q should be in UTC (end with Z)", input, result)
			}
		})
	}
}
