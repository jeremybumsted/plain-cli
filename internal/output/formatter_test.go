package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	formatter := New(FormatTable)
	if formatter == nil {
		t.Fatal("New() returned nil")
	}

	if formatter.format != FormatTable {
		t.Errorf("Expected format %v, got %v", FormatTable, formatter.format)
	}

	if formatter.quiet {
		t.Error("FormatTable should not be quiet")
	}
}

func TestQuietMode(t *testing.T) {
	formatter := New(FormatQuiet)
	if !formatter.IsQuiet() {
		t.Error("FormatQuiet should be quiet")
	}

	if !formatter.quiet {
		t.Error("quiet flag should be true")
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatJSON, &buf)

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := formatter.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "key1") {
		t.Error("Output should contain key1")
	}

	if !strings.Contains(output, "value1") {
		t.Error("Output should contain value1")
	}
}

func TestPrintTable(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	headers := []string{"ID", "Name", "Status"}
	rows := [][]string{
		{"1", "Thread A", "Open"},
		{"2", "Thread B", "Closed"},
	}

	err := formatter.PrintTable(headers, rows)
	if err != nil {
		t.Fatalf("PrintTable failed: %v", err)
	}

	output := buf.String()

	// Check headers are present
	if !strings.Contains(output, "ID") || !strings.Contains(output, "Name") || !strings.Contains(output, "Status") {
		t.Error("Output should contain all headers")
	}

	// Check data is present
	if !strings.Contains(output, "Thread A") || !strings.Contains(output, "Thread B") {
		t.Error("Output should contain all row data")
	}
}

func TestPrintTableEmpty(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	headers := []string{"ID", "Name"}
	rows := [][]string{}

	err := formatter.PrintTable(headers, rows)
	if err != nil {
		t.Fatalf("PrintTable failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No results found") {
		t.Error("Should display 'No results found' for empty table")
	}
}

func TestPrintKeyValue(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	pairs := map[string]string{
		"ID":     "123",
		"Title":  "Test Thread",
		"Status": "Open",
	}

	err := formatter.PrintKeyValue(pairs)
	if err != nil {
		t.Fatalf("PrintKeyValue failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID") || !strings.Contains(output, "123") {
		t.Error("Output should contain key-value pairs")
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Success("Operation completed")
	if err != nil {
		t.Fatalf("Success failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Operation completed") {
		t.Error("Output should contain success message")
	}

	if !strings.Contains(output, "✓") {
		t.Error("Output should contain success symbol")
	}
}

func TestSuccessQuiet(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatQuiet, &buf)

	err := formatter.Success("Operation completed")
	if err != nil {
		t.Fatalf("Success failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Error("Quiet mode should not output success messages")
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Error("Something went wrong")
	if err != nil {
		t.Fatalf("Error failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Something went wrong") {
		t.Error("Output should contain error message")
	}

	if !strings.Contains(output, "✗") {
		t.Error("Output should contain error symbol")
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Info("Informational message")
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Informational message") {
		t.Error("Output should contain info message")
	}
}

func TestInfoQuiet(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatQuiet, &buf)

	err := formatter.Info("Informational message")
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Error("Quiet mode should not output info messages")
	}
}

func TestWarning(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Warning("Warning message")
	if err != nil {
		t.Fatalf("Warning failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Warning message") {
		t.Error("Output should contain warning message")
	}

	if !strings.Contains(output, "⚠") {
		t.Error("Output should contain warning symbol")
	}
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Print("Raw output")
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Raw output") {
		t.Error("Output should contain raw text")
	}
}

func TestPrintf(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewWithWriter(FormatTable, &buf)

	err := formatter.Printf("Formatted: %s = %d\n", "value", 42)
	if err != nil {
		t.Fatalf("Printf failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Formatted: value = 42") {
		t.Error("Output should contain formatted text")
	}
}
