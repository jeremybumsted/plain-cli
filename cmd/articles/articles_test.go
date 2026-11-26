package articles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func TestArticlesCmdStructure(t *testing.T) {
	cmd := &ArticlesCmd{}
	_ = cmd.List
	_ = cmd.Get
	_ = cmd.Upsert
}

func TestNormalizeArticleStatus(t *testing.T) {
	status, err := normalizeArticleStatus("published")
	if err != nil {
		t.Fatalf("normalizeArticleStatus failed: %v", err)
	}
	if status != "PUBLISHED" {
		t.Fatalf("expected PUBLISHED, got %q", status)
	}
}

func TestUpsertCmdResolveContentHTMLPrefersFile(t *testing.T) {
	tmpDir := t.TempDir()
	contentPath := filepath.Join(tmpDir, "article.html")
	if err := os.WriteFile(contentPath, []byte("<p>from file</p>"), 0600); err != nil {
		t.Fatalf("failed to write temp content file: %v", err)
	}

	content, err := (&UpsertCmd{Content: "<p>inline</p>", ContentFile: contentPath}).resolveContentHTML()
	if err != nil {
		t.Fatalf("resolveContentHTML failed: %v", err)
	}
	if content != "<p>from file</p>" {
		t.Fatalf("expected file content, got %q", content)
	}
}

func TestUpsertCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	err := (&UpsertCmd{
		ConfigPath:   configPath,
		Format:       "table",
		HelpCenterID: "hc_123",
		Title:        "Test",
		Description:  "Description",
		Content:      "<p>Hello</p>",
	}).Run()
	if err == nil {
		t.Fatal("expected auth error")
	}
}
