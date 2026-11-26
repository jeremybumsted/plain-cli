package main

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestRootFlagsPropagateToThreadsList(t *testing.T) {
	cli := &CLI{}
	parser, err := kong.New(cli)
	if err != nil {
		t.Fatalf("failed to build parser: %v", err)
	}

	_, err = parser.Parse([]string{"--config", "/tmp/plain.json", "--json", "threads", "list"})
	if err != nil {
		t.Fatalf("failed to parse CLI: %v", err)
	}

	if cli.Threads.List.ConfigPath != "/tmp/plain.json" {
		t.Fatalf("expected config path to propagate, got %q", cli.Threads.List.ConfigPath)
	}
	if cli.Threads.List.Format != "json" {
		t.Fatalf("expected format to be json, got %q", cli.Threads.List.Format)
	}
}

func TestRootQuietWinsOverJSON(t *testing.T) {
	cli := &CLI{}
	parser, err := kong.New(cli)
	if err != nil {
		t.Fatalf("failed to build parser: %v", err)
	}

	_, err = parser.Parse([]string{"--json", "--quiet", "articles", "get", "hca_123"})
	if err != nil {
		t.Fatalf("failed to parse CLI: %v", err)
	}

	if cli.Articles.Get.Format != "quiet" {
		t.Fatalf("expected quiet format, got %q", cli.Articles.Get.Format)
	}
}

func TestRootConfigPropagatesToConfigShow(t *testing.T) {
	cli := &CLI{}
	parser, err := kong.New(cli)
	if err != nil {
		t.Fatalf("failed to build parser: %v", err)
	}

	_, err = parser.Parse([]string{"--config", "/tmp/plain.json", "config", "show"})
	if err != nil {
		t.Fatalf("failed to parse CLI: %v", err)
	}

	if cli.Config.Show.ConfigPath != "/tmp/plain.json" {
		t.Fatalf("expected config path to propagate, got %q", cli.Config.Show.ConfigPath)
	}
}

func TestRootFlagsPropagateToHelpCenterGroupCreate(t *testing.T) {
	cli := &CLI{}
	parser, err := kong.New(cli)
	if err != nil {
		t.Fatalf("failed to build parser: %v", err)
	}

	_, err = parser.Parse([]string{"--config", "/tmp/plain.json", "--quiet", "help-centers", "group", "create", "--help-center-id", "hc_123", "--name", "Docs"})
	if err != nil {
		t.Fatalf("failed to parse CLI: %v", err)
	}

	if cli.HelpCenters.Group.Create.ConfigPath != "/tmp/plain.json" {
		t.Fatalf("expected config path to propagate, got %q", cli.HelpCenters.Group.Create.ConfigPath)
	}
	if cli.HelpCenters.Group.Create.Format != "quiet" {
		t.Fatalf("expected format to be quiet, got %q", cli.HelpCenters.Group.Create.Format)
	}
}
