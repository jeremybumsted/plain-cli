package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents the output format type
type Format string

const (
	// FormatTable outputs data in a human-readable table format
	FormatTable Format = "table"
	// FormatJSON outputs data as JSON
	FormatJSON Format = "json"
	// FormatQuiet outputs minimal information
	FormatQuiet Format = "quiet"
)

// Formatter handles output formatting for CLI commands
type Formatter struct {
	format Format
	writer io.Writer
	quiet  bool
}

// New creates a new Formatter with the specified format
func New(format Format) *Formatter {
	return &Formatter{
		format: format,
		writer: os.Stdout,
		quiet:  format == FormatQuiet,
	}
}

// NewWithWriter creates a new Formatter with a custom writer (useful for testing)
func NewWithWriter(format Format, writer io.Writer) *Formatter {
	return &Formatter{
		format: format,
		writer: writer,
		quiet:  format == FormatQuiet,
	}
}

// PrintJSON outputs data as JSON
func (f *Formatter) PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintTable outputs data as a formatted table
// headers: column headers
// rows: slice of row data (each row is a slice of strings)
func (f *Formatter) PrintTable(headers []string, rows [][]string) error {
	if len(rows) == 0 {
		return f.Info("No results found")
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	separators := make([]string, len(headers))
	for i := range separators {
		separators[i] = strings.Repeat("-", len(headers[i]))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return w.Flush()
}

// PrintKeyValue outputs key-value pairs in a formatted way
func (f *Formatter) PrintKeyValue(pairs map[string]string) error {
	if f.format == FormatJSON {
		return f.PrintJSON(pairs)
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	for key, value := range pairs {
		fmt.Fprintf(w, "%s:\t%s\n", key, value)
	}
	return w.Flush()
}

// Success prints a success message (unless in quiet mode)
func (f *Formatter) Success(message string) error {
	if f.quiet {
		return nil
	}
	_, err := fmt.Fprintln(f.writer, "✓", message)
	return err
}

// Error prints an error message
func (f *Formatter) Error(message string) error {
	_, err := fmt.Fprintln(f.writer, "✗", message)
	return err
}

// Info prints an informational message (unless in quiet mode)
func (f *Formatter) Info(message string) error {
	if f.quiet {
		return nil
	}
	_, err := fmt.Fprintln(f.writer, message)
	return err
}

// Warning prints a warning message
func (f *Formatter) Warning(message string) error {
	if f.quiet {
		return nil
	}
	_, err := fmt.Fprintln(f.writer, "⚠", message)
	return err
}

// Print outputs raw text
func (f *Formatter) Print(message string) error {
	_, err := fmt.Fprintln(f.writer, message)
	return err
}

// Printf outputs formatted text
func (f *Formatter) Printf(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(f.writer, format, args...)
	return err
}

// FormatData is a generic formatter that handles different output types
func (f *Formatter) FormatData(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.PrintJSON(data)
	case FormatQuiet:
		// In quiet mode, try to extract minimal info
		return f.formatQuiet(data)
	default:
		// Default to some string representation
		return f.Info(fmt.Sprintf("%v", data))
	}
}

// formatQuiet attempts to extract minimal information in quiet mode
func (f *Formatter) formatQuiet(data interface{}) error {
	// This can be customized based on data type
	// For now, just print the raw value
	_, err := fmt.Fprintln(f.writer, data)
	return err
}

// IsQuiet returns true if the formatter is in quiet mode
func (f *Formatter) IsQuiet() bool {
	return f.quiet
}

// GetFormat returns the current format
func (f *Formatter) GetFormat() Format {
	return f.format
}
