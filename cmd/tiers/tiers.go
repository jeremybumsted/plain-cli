package tiers

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// TiersCmd represents the tiers command group
type TiersCmd struct {
	List ListCmd `cmd:"" help:"List all tiers"`
	Get  GetCmd  `cmd:"" help:"Get tier details"`
}

func getConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
}

func getClient(cfg *config.Config) (*plain.Client, error) {
	token, err := cfg.GetToken()
	if err != nil {
		return nil, err
	}
	return plain.NewClient(token), nil
}

func getFormatter(format string) *output.Formatter {
	var fmt output.Format
	switch format {
	case "json":
		fmt = output.FormatJSON
	case "quiet":
		fmt = output.FormatQuiet
	default:
		fmt = output.FormatTable
	}
	return output.New(fmt)
}
