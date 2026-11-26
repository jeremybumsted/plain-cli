package companies

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// CompaniesCmd represents the companies command group.
type CompaniesCmd struct {
	List ListCmd `cmd:"" help:"List companies"`
	Get  GetCmd  `cmd:"" help:"Get company by ID"`
}

func getConfig(configPath string) (*config.Config, error) { return config.Load(configPath) }

func getClient(cfg *config.Config) (*plain.Client, error) {
	token, err := cfg.GetToken()
	if err != nil {
		return nil, err
	}
	return plain.NewClient(token), nil
}

func getFormatter(format string) *output.Formatter {
	var out output.Format
	switch format {
	case "json":
		out = output.FormatJSON
	case "quiet":
		out = output.FormatQuiet
	default:
		out = output.FormatTable
	}
	return output.New(out)
}
