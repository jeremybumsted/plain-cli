package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/jeremybumsted/plain-cli/cmd/articles"
	"github.com/jeremybumsted/plain-cli/cmd/config"
	"github.com/jeremybumsted/plain-cli/cmd/threads"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// CLI is the root command structure
type CLI struct {
	Threads  threads.ThreadsCmd `cmd:"" help:"Manage support threads"`
	Articles articles.ArticlesCmd `cmd:"" help:"Manage help center articles"`
	Config   config.ConfigCmd   `cmd:"" help:"Configure Plain CLI"`

	// Global flags
	ConfigPath string `name:"config" help:"Path to config file" type:"path" env:"PLAIN_CONFIG"`
	JSON       bool   `help:"Output as JSON" short:"j"`
	Quiet      bool   `help:"Minimal output" short:"q"`

	// Version command
	Version VersionCmd `cmd:"" help:"Show version information"`
}

// VersionCmd shows version information
type VersionCmd struct{}

// Run executes the version command
func (cmd *VersionCmd) Run() error {
	fmt.Printf("plain-cli version %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("built at: %s\n", date)
	fmt.Println("https://github.com/jeremybumsted/plain-cli")
	return nil
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli,
		kong.Name("plain"),
		kong.Description("Plain support in the CLI - manage threads, customers, and more"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
