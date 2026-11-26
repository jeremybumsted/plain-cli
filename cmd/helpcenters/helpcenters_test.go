package helpcenters

import (
	"path/filepath"
	"testing"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func TestHelpCentersCmdStructure(t *testing.T) {
	cmd := &HelpCentersCmd{}
	_ = cmd.List
	_ = cmd.Get
	_ = cmd.Group
}

func TestListCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	err := (&ListCmd{ConfigPath: configPath, Format: "table"}).Run()
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestGroupCreateCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	err := (&GroupCreateCmd{ConfigPath: configPath, Format: "table", HelpCenterID: "hc_123", Name: "Docs"}).Run()
	if err == nil {
		t.Fatal("expected auth error")
	}
}
